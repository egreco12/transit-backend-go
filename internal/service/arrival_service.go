package service

import (
	"context"
	"sort"

	"github.com/egreco12/transit-backend-go/internal/domain"
	"github.com/egreco12/transit-backend-go/internal/oba"
)

type OBAClient interface {
	ArrivalsForStop(ctx context.Context, stopID string) (*oba.ArrivalsResponse, error)
}

type ArrivalService struct {
	oba OBAClient
}

func NewArrivalService(obaClient OBAClient) *ArrivalService {
	return &ArrivalService{
		oba: obaClient,
	}
}

func (s *ArrivalService) GetArrivalsForStop(ctx context.Context, stopID string) ([]domain.Arrival, error) {
	resp, err := s.oba.ArrivalsForStop(ctx, stopID)
	if err != nil {
		return nil, err
	}

	now := resp.CurrentTime

	routeById := make(map[string]oba.Route)
	for _, r := range resp.Data.References.Routes {
		routeById[r.ID] = r
	}

	tripByID := make(map[string]oba.Trip)
	for _, t := range resp.Data.References.Trips {
		tripByID[t.ID] = t
	}

	var result []domain.Arrival
	for _, ad := range resp.Data.Entry.ArrivalsAndDepartures {
		arrivalTime := ad.ScheduledArrivalTime
		if ad.Predicted && ad.PredictedArrivalTime > 0 {
			arrivalTime = ad.PredictedArrivalTime
		}

		if arrivalTime <= 0 {
			continue
		}

		etaSeconds := int((arrivalTime - now) / 1000)

		route := routeById[ad.RouteID]
		trip := tripByID[ad.TripID]

		headsign := ""
		switch {
		case trip.TripHeadsign != "":
			headsign = trip.TripHeadsign
		case ad.Headsign != "":
			headsign = ad.Headsign
		case route.LongName != "":
			headsign = route.LongName
		default:
			headsign = "Unknown"
		}

		result = append(result, domain.Arrival{
			RouteID:        ad.RouteID,
			RouteShortName: route.ShortName,
			Headsign:       headsign,
			ETASeconds:     etaSeconds,
			ArrivalTimeMs:  arrivalTime,
			Predicted:      ad.Predicted,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ArrivalTimeMs < result[j].ArrivalTimeMs
	})

	return result, nil
}
