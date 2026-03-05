package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jhump/protoreflect/dynamic"
	grpcdynamic "github.com/jhump/protoreflect/dynamic/grpcdynamic"
	grpcreflect "github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type grpcServicesRequest struct {
	Target string `json:"target"`
}

type grpcMethodInfo struct {
	Name       string `json:"name"`
	FullMethod string `json:"fullMethod"`
	InputType  string `json:"inputType"`
	OutputType string `json:"outputType"`
}

type grpcServiceInfo struct {
	Name    string           `json:"name"`
	Methods []grpcMethodInfo `json:"methods"`
}

type grpcServicesResponse struct {
	Services []grpcServiceInfo `json:"services"`
}

type grpcInvokeRequest struct {
	Target     string          `json:"target"`
	FullMethod string          `json:"fullMethod"`
	JsonBody   json.RawMessage `json:"jsonBody"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func handleGRPCServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "only POST is allowed"})
		return
	}

	var req grpcServicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body: " + err.Error()})
		return
	}
	if strings.TrimSpace(req.Target) == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "target is required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	conn, client, err := newGRPCReflectionClient(ctx, req.Target)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to connect via reflection: " + err.Error()})
		return
	}
	defer conn.Close()
	defer client.Reset()

	serviceNames, err := client.ListServices()
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to list services (is reflection enabled?): " + err.Error()})
		return
	}

	var resp grpcServicesResponse
	for _, svcName := range serviceNames {
		// фильтруем служебный reflection-сервис
		if strings.HasPrefix(svcName, "grpc.reflection.") {
			continue
		}

		sd, err := client.ResolveService(svcName)
		if err != nil {
			continue
		}
		svc := grpcServiceInfo{Name: svcName}
		for _, m := range sd.GetMethods() {
			fullMethod := "/" + svcName + "/" + m.GetName()
			svc.Methods = append(svc.Methods, grpcMethodInfo{
				Name:       m.GetName(),
				FullMethod: fullMethod,
				InputType:  m.GetInputType().GetFullyQualifiedName(),
				OutputType: m.GetOutputType().GetFullyQualifiedName(),
			})
		}
		resp.Services = append(resp.Services, svc)
	}

	writeJSON(w, http.StatusOK, resp)
}

func handleGRPCInvoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "only POST is allowed"})
		return
	}

	var req grpcInvokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body: " + err.Error()})
		return
	}
	req.Target = strings.TrimSpace(req.Target)
	req.FullMethod = strings.TrimSpace(req.FullMethod)

	if req.Target == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "target is required"})
		return
	}
	if req.FullMethod == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "fullMethod is required"})
		return
	}
	if len(req.JsonBody) == 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "jsonBody is required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	conn, client, err := newGRPCReflectionClient(ctx, req.Target)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to connect via reflection: " + err.Error()})
		return
	}
	defer conn.Close()
	defer client.Reset()

	svcName, methodName, err := splitFullMethod(req.FullMethod)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	sd, err := client.ResolveService(svcName)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to resolve service: " + err.Error()})
		return
	}

	md := sd.FindMethodByName(methodName)
	if md == nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "method not found in service"})
		return
	}

	// JSON -> dynamic message: здесь происходит валидация по схеме
	inMsg := dynamic.NewMessage(md.GetInputType())
	if err := inMsg.UnmarshalJSON(req.JsonBody); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON for request message: " + err.Error()})
		return
	}

	stub := grpcdynamic.NewStub(conn)
	respMsg, err := stub.InvokeRpc(ctx, md, inMsg)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "gRPC call failed: " + err.Error()})
		return
	}

	// Ответ тоже отдаем как JSON
	type jsonMarshaler interface {
		MarshalJSON() ([]byte, error)
	}
	jm, ok := respMsg.(jsonMarshaler)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "response message does not support JSON marshaling"})
		return
	}
	data, err := jm.MarshalJSON()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to marshal response to JSON: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func newGRPCReflectionClient(ctx context.Context, target string) (*grpc.ClientConn, *grpcreflect.Client, error) {
	if target == "" {
		return nil, nil, errors.New("target is empty")
	}

	conn, err := grpc.DialContext(ctx, target,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16*1024*1024)),
	)
	if err != nil {
		return nil, nil, err
	}

	refClient := grpcreflect.NewClient(ctx, reflectionpb.NewServerReflectionClient(conn))
	return conn, refClient, nil
}

func splitFullMethod(full string) (string, string, error) {
	full = strings.TrimSpace(full)
	full = strings.TrimPrefix(full, "/")
	parts := strings.Split(full, "/")
	if len(parts) != 2 {
		return "", "", errors.New("fullMethod must be in form /package.Service/Method")
	}
	return parts[0], parts[1], nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

