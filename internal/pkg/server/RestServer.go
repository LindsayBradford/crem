// Copyright (c) 2018 Australian Rivers Institute.

// See: https://astaxie.gitbooks.io/build-web-application-with-golang/en/03.2.html

package server

import (
	"fmt"
	engineApi "github.com/LindsayBradford/crem/cmd/cremengine/engine/api"
	"github.com/LindsayBradford/crem/internal/pkg/server/admin"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type RestServer struct {
	adminMux *admin.Mux
	apiMux   rest.Mux

	cacheMaximumAgeInSeconds uint64
	apiPort                  uint64
	adminPort                uint64

	Logger logging.Logger
}

func (s *RestServer) Initialise() *RestServer {
	s.adminMux = new(admin.Mux).Initialise()
	return s
}

func (s *RestServer) WithApiMux(apiMux rest.Mux) *RestServer {
	s.adminMux = new(admin.Mux).Initialise()
	s.apiMux = apiMux
	s.apiMux.AddHandler("^/$", s.adminMux.StatusHandler)
	return s
}

func (s *RestServer) WithApiPort(apiPort uint64) *RestServer {
	s.apiPort = apiPort
	return s
}

func (s *RestServer) WithAdminPort(adminPort uint64) *RestServer {
	s.adminPort = adminPort
	return s
}

func (s *RestServer) WithCacheMaximumAge(cacheMaximumAge uint64) *RestServer {
	s.cacheMaximumAgeInSeconds = cacheMaximumAge
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
		s.apiMux.SetCacheMaxAge(s.cacheMaximumAgeInSeconds)
		startMuxOnPort(s.apiMux, s.apiPort)
	}()

	go func() {
		s.adminMux.SetCacheMaxAge(s.cacheMaximumAgeInSeconds)
		startMuxOnPort(s.adminMux, s.adminPort)
	}()

	s.adminMux.SetStatus("RUNNING")

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

func (s *RestServer) SetScenario(scenarioFilePath string) {
	if engineApiMux, isEngineApiMux := s.apiMux.(*engineApi.Mux); isEngineApiMux {
		engineApiMux.SetScenario(scenarioFilePath)
	}
}

func (s *RestServer) SetSolution(solutionFilePath string) {
	if engineApiMux, isEngineApiMux := s.apiMux.(*engineApi.Mux); isEngineApiMux {
		engineApiMux.SetSolution(solutionFilePath)
	}
}
