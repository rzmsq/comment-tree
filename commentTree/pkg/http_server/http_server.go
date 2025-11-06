package httpserver

import (
	"CommentTree/commentTree/pkg/handler"
	"CommentTree/commentTree/pkg/http_server/config"
	"context"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(handler *handler.Handler, cfg config.ServerConfig) *Server {
	return &Server{server: &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.MaxReadTimeout,
		WriteTimeout: cfg.MaxWriteTimeout,
	}}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
