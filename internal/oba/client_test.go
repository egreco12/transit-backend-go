package oba

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestArrivalsForStop_SendsCorrectRequestAndParsesResponse(t *testing.T) {
	const apiKey = "test-key"
	const stopID = "STOP_1"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}

		expectedPath := "/arrivals-and-departures-for-stop/" + stopID + ".json"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %s, want %s", r.URL.Path, expectedPath)
		}

		if gotKey := r.URL.Query().Get("key"); gotKey != apiKey {
			t.Errorf("query key = %s, want %s", gotKey, apiKey)
		}

		now := time.Now().UnixMilli()

		resp := ArrivalsResponse{
			Code:        200,
			CurrentTime: now,
			Data: Data{
				Entry: Entry{
					ArrivalsAndDepartures: []ArrivalAndDeparture{
						{
							RouteID:              "ROUTE_1",
							TripID:               "TRIP_1",
							ScheduledArrivalTime: now + 5*60*1000,
							PredictedArrivalTime: now + 4*60*1000,
							Predicted:            true,
							StopID:               stopID,
							VehicleID:            "VEH_1",
							Headsign:             "To Downtown",
						},
					},
				},
				References: References{
					Routes: []Route{
						{
							ID:        "ROUTE_1",
							ShortName: "10",
							LongName:  "Downtown Express",
						},
					},
					Trips: []Trip{
						{
							ID:           "TRIP_1",
							RouteID:      "ROUTE_1",
							TripHeadsign: "To Downtown",
						},
					},
				},
			},
			Text:    "OK",
			Version: 2,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode fake response: %v", err)
		}
	}))
	defer ts.Close()

	client := NewClient(ts.URL, apiKey)

	got, err := client.ArrivalsForStop(context.Background(), stopID)
	if err != nil {
		t.Fatalf("ArrivalsForStop returned error: %v", err)
	}

	if got == nil {
		t.Fatal("ArrivalsForStop returned nil response")
	}

	if got.Code != 200 {
		t.Errorf("Code = %d, want 200", got.Code)
	}

	if len(got.Data.Entry.ArrivalsAndDepartures) != 1 {
		t.Fatalf("len(arrivals) = %d, want 1", len(got.Data.Entry.ArrivalsAndDepartures))
	}

	ad := got.Data.Entry.ArrivalsAndDepartures[0]
	if ad.RouteID != "ROUTE_1" {
		t.Errorf("RouteID = %s, want ROUTE_1", ad.RouteID)
	}
	if ad.StopID != stopID {
		t.Errorf("StopID = %s, want %s", ad.StopID, stopID)
	}
}

func TestArrivalsForStop_HandlesErrorStatusCode(t *testing.T) {
	const apiKey = "test-key"
	const stopID = "STOP_1"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, apiKey)

	resp, err := client.ArrivalsForStop(context.Background(), stopID)
	if err == nil {
		t.Fatalf("expected error for 401 response, got nil")
	}
	if resp != nil {
		t.Fatalf("expected nil response on error, got %#v", resp)
	}
}
