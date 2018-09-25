// Copyright (c) 2018 Australian Rivers Institute.

// See: https://astaxie.gitbooks.io/build-web-application-with-golang/en/03.2.html

package server

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/logging/handlers"
)

type Response struct {
	ResponseCode int
	Message      string
	Time         string
}

type RestServer struct {
	adminMux *AdminMux
	apiMux   *ApiMux

	configuration *config.HttpServerConfig
	Logger        handlers.LogHandler
}

type StartableMux interface {
	Start(portAddress string)
}

func (s *RestServer) Initialise() *RestServer {
	s.adminMux = new(AdminMux).Initialise().WithType("Admin")
	s.apiMux = new(ApiMux).Initialise().WithType("API")
	return s
}

func (s *RestServer) WithLogger(logger handlers.LogHandler) *RestServer {
	s.Logger = logger
	s.adminMux.Logger = logger
	s.apiMux.Logger = logger
	return s
}

func (s *RestServer) WithConfig(configuration *config.HttpServerConfig) *RestServer {
	s.configuration = configuration
	return s
}

func (s *RestServer) WithStatus(status Status) *RestServer {
	s.adminMux.Status = status
	return s
}

func (s *RestServer) Start() {
	go func() {
		startMuxOnPort(s.adminMux, s.configuration.AdminPort)
	}()

	go func() {
		startMuxOnPort(s.apiMux, s.configuration.ApiPort)
	}()

	s.adminMux.WaitForShutdownSignal()
	s.shutdown()
}

func startMuxOnPort(mux StartableMux, portNumber uint64) {
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

func (s *RestServer) AddApiMapping(address string, handlerFunction http.HandlerFunc) {
	s.apiMux.handlerMap[address] = handlerFunction
}
