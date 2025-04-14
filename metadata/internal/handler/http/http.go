package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieapp.com/metadata/internal/controller/metadata"
	"movieapp.com/metadata/internal/repository"
)

// identifies movie metadata HTTP handler
type Handler struct {
	ctrl *metadata.Controller
}

// create new metadata HTTP handler
func New(ctrl *metadata.Controller) *Handler{
	return &Handler{ctrl}
}

// GetMetadata handler GET \metadata requests
func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request)  {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	m, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, repository.ErrNotFound){
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response encode error %v\n", err)
	}
}