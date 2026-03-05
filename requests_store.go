package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const requestsFile = "saved_requests.json"

var requestsMu sync.Mutex

type SavedHTTPConfig struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body"`
}

type SavedGRPCConfig struct {
	Target     string `json:"target"`
	FullMethod string `json:"fullMethod"`
	Body       string `json:"body"`
}

type SavedRequest struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Type      string           `json:"type"` // "http" | "grpc"
	HTTP      *SavedHTTPConfig `json:"http,omitempty"`
	GRPC      *SavedGRPCConfig `json:"grpc,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

type savedRequestsFile struct {
	Items []SavedRequest `json:"items"`
}

type saveRequestInput struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"` // "http" | "grpc"
	HTTP *struct {
		Method string `json:"method"`
		URL    string `json:"url"`
		Body   string `json:"body"`
	} `json:"http,omitempty"`
	GRPC *struct {
		Target     string `json:"target"`
		FullMethod string `json:"fullMethod"`
		Body       string `json:"body"`
	} `json:"grpc,omitempty"`
}

type deleteRequestInput struct {
	ID string `json:"id"`
}

type savedRequestsListResponse struct {
	Items []SavedRequest `json:"items"`
}

func handleSavedRequestsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "only GET is allowed"})
		return
	}

	requestsMu.Lock()
	defer requestsMu.Unlock()

	items, err := loadSavedRequestsLocked()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to load saved requests: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, savedRequestsListResponse{Items: items})
}

func handleSavedRequestSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "only POST is allowed"})
		return
	}

	var input saveRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body: " + err.Error()})
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Type = strings.ToLower(strings.TrimSpace(input.Type))
	if input.Name == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "name is required"})
		return
	}
	if input.Type != "http" && input.Type != "grpc" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "type must be \"http\" or \"grpc\""})
		return
	}

	requestsMu.Lock()
	defer requestsMu.Unlock()

	items, err := loadSavedRequestsLocked()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to load saved requests: " + err.Error()})
		return
	}

	now := time.Now().UTC()
	var saved SavedRequest

	if input.ID != "" {
		// update existing
		found := false
		for i := range items {
			if items[i].ID == input.ID {
				items[i].Name = input.Name
				items[i].Type = input.Type
				if input.Type == "http" && input.HTTP != nil {
					items[i].HTTP = &SavedHTTPConfig{
						Method: input.HTTP.Method,
						URL:    input.HTTP.URL,
						Body:   input.HTTP.Body,
					}
					items[i].GRPC = nil
				} else if input.Type == "grpc" && input.GRPC != nil {
					items[i].GRPC = &SavedGRPCConfig{
						Target:     input.GRPC.Target,
						FullMethod: input.GRPC.FullMethod,
						Body:       input.GRPC.Body,
					}
					items[i].HTTP = nil
				}
				items[i].UpdatedAt = now
				saved = items[i]
				found = true
				break
			}
		}
		if !found {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "request with given id not found"})
			return
		}
	} else {
		// create new
		newID := strconv.FormatInt(now.UnixNano(), 10)
		saved = SavedRequest{
			ID:        newID,
			Name:      input.Name,
			Type:      input.Type,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if input.Type == "http" && input.HTTP != nil {
			saved.HTTP = &SavedHTTPConfig{
				Method: input.HTTP.Method,
				URL:    input.HTTP.URL,
				Body:   input.HTTP.Body,
			}
		} else if input.Type == "grpc" && input.GRPC != nil {
			saved.GRPC = &SavedGRPCConfig{
				Target:     input.GRPC.Target,
				FullMethod: input.GRPC.FullMethod,
				Body:       input.GRPC.Body,
			}
		}
		items = append(items, saved)
	}

	if err := saveSavedRequestsLocked(items); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to save requests: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, saved)
}

func handleSavedRequestDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "only POST or DELETE is allowed"})
		return
	}

	var input deleteRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body: " + err.Error()})
		return
	}
	if strings.TrimSpace(input.ID) == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "id is required"})
		return
	}

	requestsMu.Lock()
	defer requestsMu.Unlock()

	items, err := loadSavedRequestsLocked()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to load saved requests: " + err.Error()})
		return
	}

	newItems := make([]SavedRequest, 0, len(items))
	for _, it := range items {
		if it.ID != input.ID {
			newItems = append(newItems, it)
		}
	}

	if len(newItems) == len(items) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "request with given id not found"})
		return
	}

	if err := saveSavedRequestsLocked(newItems); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to save requests: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func loadSavedRequestsLocked() ([]SavedRequest, error) {
	data, err := os.ReadFile(requestsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []SavedRequest{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []SavedRequest{}, nil
	}

	var file savedRequestsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	if file.Items == nil {
		file.Items = []SavedRequest{}
	}
	return file.Items, nil
}

func saveSavedRequestsLocked(items []SavedRequest) error {
	file := savedRequestsFile{Items: items}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(requestsFile, data, 0o644)
}

