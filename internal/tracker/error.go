package tracker

import "errors"

var (
	ErrInvalidEdge                = errors.New("edge must contain exactly 2 not empty vertices")
	ErrGraphIsUnconnected         = errors.New("graph is unconnected")
	ErrUnableFindStartAndEndPoint = errors.New("unable find start and end point in cyclic graph")
)
