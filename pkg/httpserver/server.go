package httpserver

import (
	"context"
	"net/http"
)

type Server struct {
	App *http.Server
}

func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		App: &http.Server{
			Handler: handler,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() error {
	return s.App.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.App.Shutdown(ctx)
}
