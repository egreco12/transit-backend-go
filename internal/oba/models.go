package oba

type ArrivalsResponse struct {
	Code        int64  `json:"code"`
	CurrentTime int64  `json:"currentTime"`
	Data        Data   `json:"data"`
	Text        string `json:"text"`
	Version     int64  `json:"version"`
}

type Data struct {
	Entry      Entry      `json:"entry"`
	References References `json:"references"`
}

type Entry struct {
	ArrivalsAndDepartures []ArrivalAndDeparture `json:"arrivalsAndDepartures"`
}

type ArrivalAndDeparture struct {
	RouteID                string `json:"routeId"`
	TripID                 string `json:"tripId"`
	ScheduledArrivalTime   int64  `json:"scheduledArrivalTime"`
	PredictedArrivalTime   int64  `json:"predictedArrivalTime"`
	ScheduledDepartureTime int64  `json:"scheduledDepartureTime"`
	PredictedDepartureTime int64  `json:"predictedDepartureTime"`
	StopSequence           int    `json:"stopSequence"`
	Predicted              bool   `json:"predicted"`
	StopID                 string `json:"stopId"`
	VehicleID              string `json:"vehicleId"`
	Headsign               string `json:"headsign"`
}

type References struct {
	Routes []Route `json:"routes"`
	Trips  []Trip  `json:"trips"`
}

type Route struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

type Trip struct {
	ID           string `json:"id"`
	RouteID      string `json:"routeId"`
	TripHeadsign string `json:"tripHeadsign"`
}
