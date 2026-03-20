package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type HandlerResponse struct {
	Value string `json:"value"`
}
type AuthenticatedResponse struct {
	Value string `json:"value"`
	Token string `json:"token"`
}
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type GetHandlerRequest struct {
	TaskId string `json:"taskId"`
}

type FilterRequest struct {
	Filter FilterParametesRequest `json:"filter"`
}

type FilterParametesRequest struct {
	Name       string           `json:"name"`
	Parameters ParametesRequest `json:"parameters"`
}

type ParametesRequest struct {
	Value any `json:"value"`
}

func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	// Возвращаем токен
	return parts[1], nil
}
func ExtractFiltersFromBody(r *http.Request) (*FilterRequest, error) {
	var req FilterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

func CreateAuthRequest(r *http.Request) (*AuthRequest, error) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

func CreateGetHandlerRequest(r *http.Request) (*GetHandlerRequest, error) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		return nil, errors.New("taskID is empty")
	}
	return &GetHandlerRequest{TaskId: taskID}, nil
}

func ProcessErrorAndResponse(w http.ResponseWriter, r any, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r != nil {
		if err := json.NewEncoder(w).Encode(r); err != nil {
			http.Error(w, "InternalError", http.StatusInternalServerError)
		}
	}

}
