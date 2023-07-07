package api

import (
	"context"
	"home-bar/configs"
	"home-bar/internal"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(config *configs.Config, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:           config.ServerConfig.Host + ":" + config.ServerConfig.Port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

func (s *Server) RunServer() error {
	internal.PrintMessage("Server: %v", s.server.Addr)

	return s.server.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	internal.PrintMessage("Server shutdown")

	return s.server.Shutdown(ctx)
}
