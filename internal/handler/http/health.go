package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHealthRoutes(r chi.Router) {
	r.Get("/healthz", healthzHandler)
	r.Get("/ready", readyHandler)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// Here you would check DB connection or other external dependencies
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
