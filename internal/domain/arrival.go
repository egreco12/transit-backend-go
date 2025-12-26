package domain

type Arrival struct {
	RouteID        string `json:"routeId"`
	RouteShortName string `json:"routeShortName"`
	Headsign       string `json:"headsign"`
	ETASeconds     int    `json:"etaSeconds"`
	ArrivalTimeMs  int64  `json:"arrivalTimeEpochMs"`
	Predicted      bool   `json:"predicted"`
}
