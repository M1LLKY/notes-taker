package server

import (
	"context"
	"net/http"
	"notes-taker/internal/dependencies"
	"notes-taker/internal/handler"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

func New(deps *dependencies.Dependencies) (*HTTPServer, error) {
	h := handler.New(deps)
	r := h.GetRouter()
	server := &http.Server{
		Addr:              deps.Config.HTTPListenAddr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler:           r,
	}
	return &HTTPServer{
		server: server,
	}, nil
}

func (s *HTTPServer) Run() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
