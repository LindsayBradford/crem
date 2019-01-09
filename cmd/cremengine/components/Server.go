// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/server"
	"github.com/LindsayBradford/crem/internal/pkg/server/admin"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type Server struct {
	server.RestServer
}

func (s *Server) Initialise() *Server {
	s.RestServer.Initialise()
	return s
}

// mainThreadChannel *threading.MainThreadChannel
func (s *Server) WithConfig(configuration *config.HttpServerConfig) *Server {
	s.RestServer.WithConfig(configuration)
	return s
}

func (s *Server) WithLogger(logger logging.Logger) *Server {
	s.RestServer.WithLogger(logger)
	return s
}

func (s *Server) WithStatus(status admin.ServiceStatus) *Server {
	s.RestServer.WithStatus(status)
	return s
}
