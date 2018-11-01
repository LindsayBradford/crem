// Copyright (c) 2018 Australian Rivers Institute.

// See: https://astaxie.gitbooks.io/build-web-application-with-golang/en/03.2.html

package server

import (
	"fmt"

	"github.com/LindsayBradford/crem/config"
	"github.com/LindsayBradford/crem/logging"
	"github.com/LindsayBradford/crem/server/admin"
	"github.com/LindsayBradford/crem/server/rest"
)

type RestServer struct {
	adminMux *admin.Mux
	apiMux   rest.Mux

	configuration *config.HttpServerConfig
	Logger        logging.Logger
}

func (s *RestServer) WithConfig(configuration *config.HttpServerConfig) *RestServer {
	s.configuration = configuration
	return s
}

func (s *RestServer) Initialise() *RestServer {
	s.adminMux = new(admin.Mux).Initialise()
	return s
}

func (s *RestServer) WithApiMux(apiMux rest.Mux) *RestServer {
	s.adminMux = new(admin.Mux).Initialise()
	s.apiMux = apiMux
	s.apiMux.AddHandler("/", s.adminMux.StatusHandler)
	return s
}

func (s *RestServer) WithLogger(logger logging.Logger) *RestServer {
	s.Logger = logger
	s.adminMux.SetLogger(logger)
	s.apiMux.SetLogger(logger)
	return s
}

func (s *RestServer) WithStatus(status admin.ServiceStatus) *RestServer {
	s.adminMux.Status = status
	return s
}

func (s *RestServer) Start() {
	go func() {
		s.adminMux.WithCacheMaxAge(s.configuration.CacheMaximumAgeInSeconds)
		startMuxOnPort(s.adminMux, s.configuration.AdminPort)
	}()

	go func() {
		s.adminMux.WithCacheMaxAge(s.configuration.CacheMaximumAgeInSeconds)
		startMuxOnPort(s.apiMux, s.configuration.ApiPort)
	}()

	s.adminMux.WaitForShutdownSignal()
	s.shutdown()
}

func startMuxOnPort(mux rest.Mux, portNumber uint64) {
	portAddress := toPortAddress(portNumber)
	mux.Start(portAddress)
}

func toPortAddress(portNumber uint64) string {
	return fmt.Sprintf(":%d", portNumber)
}

func (s *RestServer) shutdown() {
	s.apiMux.Shutdown()
	s.adminMux.Shutdown()
	s.Logger.Info("Shutdown complete")
}

func (s *RestServer) AddApiMapping(address string, handlerFunction rest.HandlerFunc) {
	s.apiMux.AddHandler(address, handlerFunction)
}
