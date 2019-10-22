package server

import (
	lr "github.com/Wallruzz9114/bookey/util/logger"
)

// Server ...
type Server struct {
	logger *lr.Logger
}

// New ...
func New(logger *lr.Logger) *Server {
	return &Server{logger: logger}
}

// Logger ...
func (server *Server) Logger() *lr.Logger {
	return server.logger
}
