package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Router returns a configured router for the API
func (h *Handler) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", h.Hello).Methods("GET")

	return r
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, map[string]string{"message": "Hello World!"})
}
