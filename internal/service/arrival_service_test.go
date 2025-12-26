package service

import (
	"context"
	"testing"
	"time"

	"github.com/egreco12/transit-backend-go/internal/oba"
)

type fakeOBAClient struct {
	resp *oba.ArrivalsResponse
	err  error
}

func (f *fakeOBAClient) ArrivalsForStop(ctx context.Context, stopID string) (*oba.ArrivalsResponse, error) {
	return f.resp, f.err
}

func TestGetArrivalsForStop_MapsResponseCorrectly(t *testing.T) {
	now := time.Now().UnixMilli()

	ad := oba.ArrivalAndDeparture{
		RouteID:              "ROUTE_1",
		TripID:               "TRIP_1",
		ScheduledArrivalTime: now + 5*60*1000,
		PredictedArrivalTime: now + 4*60*1000,
		Predicted:            true,
		StopID:               "STOP_1",
		VehicleID:            "VEH_1",
		Headsign:             "Headsign from AD",
	}

	route := oba.Route{
		ID:        "ROUTE_1",
		ShortName: "10",
		LongName:  "Downtown Express",
	}

	trip := oba.Trip{
		ID:           "TRIP_1",
		RouteID:      "ROUTE_1",
		TripHeadsign: "To Downtown",
	}

	resp := &oba.ArrivalsResponse{
		Code:        200,
		CurrentTime: now,
		Data: oba.Data{
			Entry: oba.Entry{
				ArrivalsAndDepartures: []oba.ArrivalAndDeparture{ad},
			},
			References: oba.References{
				Routes: []oba.Route{route},
				Trips:  []oba.Trip{trip},
			},
		},
		Text:    "OK",
		Version: 2,
	}

	svc := NewArrivalService(&fakeOBAClient{resp: resp})

	arrivals, err := svc.GetArrivalsForStop(context.Background(), "STOP_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(arrivals) != 1 {
		t.Fatalf("expected 1 arrival, got %d", len(arrivals))
	}

	a := arrivals[0]

	if a.RouteID != "ROUTE_1" {
		t.Errorf("RouteID = %s, want ROUTE_1", a.RouteID)
	}
	if a.RouteShortName != "10" {
		t.Errorf("RouteShortName = %s, want 10", a.RouteShortName)
	}
	// We prefer trip headsign over AD headsign
	if a.Headsign != "To Downtown" {
		t.Errorf("Headsign = %s, want To Downtown", a.Headsign)
	}
	if !a.Predicted {
		t.Errorf("Predicted = false, want true")
	}

	// ETA should be ~4 minutes, allow a range
	if a.ETASeconds < 3*60 || a.ETASeconds > 5*60 {
		t.Errorf("ETASeconds = %d, want roughly 4min", a.ETASeconds)
	}
}

func TestGetArrivalsForStop_EmptyOnNoData(t *testing.T) {
	resp := &oba.ArrivalsResponse{
		Code:        200,
		CurrentTime: time.Now().UnixMilli(),
		Data: oba.Data{
			Entry: oba.Entry{
				ArrivalsAndDepartures: []oba.ArrivalAndDeparture{},
			},
			References: oba.References{},
		},
	}

	svc := NewArrivalService(&fakeOBAClient{resp: resp})

	arrivals, err := svc.GetArrivalsForStop(context.Background(), "STOP_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(arrivals) != 0 {
		t.Errorf("expected 0 arrivals, got %d", len(arrivals))
	}
}

func TestGetArrivalsForStop_PropagatesError(t *testing.T) {
	svc := NewArrivalService(&fakeOBAClient{resp: nil, err: assertErr{}})

	_, err := svc.GetArrivalsForStop(context.Background(), "STOP_1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

type assertErr struct{}

func (assertErr) Error() string { return "boom" }
