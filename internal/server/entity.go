package server

type TrackPathRequest struct {
	Flights [][]string `json:"flights"`
}

type TrackPathResponse struct {
	StartEndFlights []string `json:"start-end-flights"`
}
