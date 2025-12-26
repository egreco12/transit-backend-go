package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/egreco12/transit-backend-go/internal/domain"
	"github.com/go-chi/chi/v5"
)

type fakeArrivalService struct {
	resp []domain.Arrival
	err  error
}

func (f *fakeArrivalService) GetArrivalsForStop(ctx context.Context, stopID string) ([]domain.Arrival, error) {
	return f.resp, f.err
}

func TestGetArrivals_Handler(t *testing.T) {
	arr := domain.Arrival{
		RouteID:        "ROUTE_1",
		RouteShortName: "10",
		Headsign:       "Downtown",
		ETASeconds:     240,
		ArrivalTimeMs:  1234567890,
		Predicted:      true,
	}

	svc := &fakeArrivalService{resp: []domain.Arrival{arr}}

	h := &ArrivalsHandler{service: svc}

	r := chi.NewRouter()
	r.Get("/api/stops/{stopID}/arrivals", h.GetArrivals)

	req := httptest.NewRequest(http.MethodGet, "/api/stops/1_75403/arrivals", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("status = %d, want %d", status, http.StatusOK)
	}

	var got []domain.Arrival
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 arrival, got %d", len(got))
	}

	if got[0].RouteShortName != "10" {
		t.Errorf("RouteShortName = %s, want 10", got[0].RouteShortName)
	}
	if got[0].Headsign != "Downtown" {
		t.Errorf("Headsign = %s, want Downtown", got[0].Headsign)
	}
}
