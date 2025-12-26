package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/egreco12/transit-backend-go/internal/service"
	"github.com/go-chi/chi/v5"
)

type ArrivalsHandler struct {
	service *service.ArrivalService
}

func NewArrivalsHandler(svc *service.ArrivalService) *ArrivalsHandler {
	return &ArrivalsHandler{service: svc}
}

func (h *ArrivalsHandler) GetArrivals(w http.ResponseWriter, r *http.Request) {
	stopID := chi.URLParam(r, "stopID")
	if stopID == "" {
		http.Error(w, "missing stopID", http.StatusBadRequest)
		return
	}

	arrivals, err := h.service.GetArrivalsForStop(r.Context(), stopID)
	if err != nil {
		http.Error(w, "failed to fetch arrivals: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(arrivals); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
