package server

import (
	"encoding/json"
	"errors"
	"flights/internal/tracker"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

func (s *Server) TrackPath(ctx *fasthttp.RequestCtx) {
	contentType := ctx.Request.Header.ContentType()
	if string(contentType) != "application/json" {
		log.Printf("content type error: %v", string(contentType))
		ctx.Error(ErrWrongContentType.Error(), fasthttp.StatusBadRequest)
		return
	}

	var data TrackPathRequest

	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		log.Printf("error parsing request: %v", err)
		ctx.Error(ErrInvalidJson.Error(), fasthttp.StatusBadRequest)
		return
	}

	if len(data.Flights) == 0 {
		log.Printf("empty request")
		ctx.Error(ErrEmptyRequest.Error(), fasthttp.StatusBadRequest)
		return
	}

	graph, err := tracker.BuildGraph(data.Flights)
	if err != nil {
		log.Println(err.Error())
		ctx.Error(ErrInvalidAirportCodes.Error(), fasthttp.StatusBadRequest)
		return
	}

	start, end, err := graph.FindStartAndEndPoint()
	if err != nil {
		log.Println(err.Error())

		if errors.Is(err, tracker.ErrGraphIsUnconnected) {
			ctx.Error(ErrMultipleTrips.Error(), fasthttp.StatusBadRequest)
		} else {
			ctx.Error(ErrRouteTrip.Error(), fasthttp.StatusBadRequest)
		}
		return
	}

	body, err := json.Marshal(&TrackPathResponse{StartEndFlights: []string{start, end}})
	if err != nil {
		log.Println(err.Error())
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
}

func (s *Server) Health(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(fmt.Sprintf("Health: alive at %v", time.Now().Unix()))
}
