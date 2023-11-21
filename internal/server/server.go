package server

import (
	fasthttprouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Server struct {
	serv *fasthttp.Server
}

func NewServer() *Server {
	server := &Server{
		serv: &fasthttp.Server{},
	}

	router := fasthttprouter.New()
	router.POST("/calculate", server.TrackPath)
	router.GET("/health", server.Health)

	server.serv.Handler = router.Handler

	return server
}

func (s *Server) Listen(addr string) error {
	return fasthttp.ListenAndServe(addr, s.serv.Handler)
}
