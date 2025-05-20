package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"movieapp.com/rating/internal/controller/rating"
	model "movieapp.com/rating/pkg"
)

type Handler struct {
	ctrl *rating.Controller
}

// New creates a new rating service controller
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

// handling http requests to our service
func (h *Handler) Handle(w http.ResponseWriter, req *http.Request)  {
	recordId := model.RecordId(req.FormValue("id"))

	if recordId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(req.FormValue("type"))

	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(req.Context(), recordId, recordType)

		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("Response encode error %v\n", err)
		}

	case http.MethodPut:
		userId := model.UserID(req.FormValue("userId"))
		v, err := strconv.ParseFloat(req.FormValue("value"), 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h.ctrl.PutRating(req.Context(), recordId, recordType, &model.Rating{UserID: userId, Value: model.RatingValue(v)}); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}