package server

import "errors"

var (
	ErrEmptyRequest        = errors.New("request can't be empty")
	ErrInvalidAirportCodes = errors.New("invalid airport codes")
	ErrMultipleTrips       = errors.New("multiple paths detected: unable find departure and arrival airport")
	ErrRouteTrip           = errors.New("route trip detected: unable find departure and arrival airport")
	ErrWrongContentType    = errors.New("wrong content type")
	ErrInvalidJson         = errors.New("invalid body: only strings are acceptable")
)
