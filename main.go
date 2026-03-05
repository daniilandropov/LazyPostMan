package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Root UI – простая HTML-страница, дальше расширим под полноценный UI
	mux.HandleFunc("/", serveIndex)

	// gRPC reflection: получить список сервисов и методов
	mux.HandleFunc("/api/grpc/services", handleGRPCServices)
	// gRPC invoke: провалидировать JSON и выполнить вызов
	mux.HandleFunc("/api/grpc/invoke", handleGRPCInvoke)

	// Saved requests (как коллекции в Postman)
	mux.HandleFunc("/api/requests", handleSavedRequestsList)
	mux.HandleFunc("/api/requests/save", handleSavedRequestSave)
	mux.HandleFunc("/api/requests/delete", handleSavedRequestDelete)

	// TODO: позже можно добавить HTTP proxy-эндпоинт:
	// - /api/http/request

	addr := ":8089"
	log.Printf("lazyPostman UI is running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
