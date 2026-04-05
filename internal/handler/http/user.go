package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juanpblasi/go-template/internal/service"
	"github.com/juanpblasi/go-template/pkg/errors"
)

type userHandler struct {
	svc service.UserService
}

func RegisterUserRoutes(r chi.Router, svc service.UserService) {
	h := &userHandler{svc: svc}
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUser)
	})
}

// ErrorResponse is a standard format for HTTP errors
type ErrorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func mapErrorToHTTP(w http.ResponseWriter, err error) {
	if errors.IsType(err, errors.ErrNotFound) {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.IsType(err, errors.ErrInvalidRequest) {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		mapErrorToHTTP(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	user, err := h.svc.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		mapErrorToHTTP(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
