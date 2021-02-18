// Copyright (c) 2020 Australian Rivers Institute.

package engine

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/config"
	"github.com/LindsayBradford/crem/internal/pkg/server"
	"github.com/LindsayBradford/crem/internal/pkg/server/admin"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

type Engine interface {
	LogHandler() logging.Logger
	SetScenario(scenarioFilePath string)
	Run() error
}

func NewBaseEngine() *BaseEngine {
	engine := new(BaseEngine)
	engine.Initialise()
	return engine
}

var _ Engine = new(BaseEngine)

type BaseEngine struct {
	server.RestServer
}

func (s *BaseEngine) WithLogHandler(logger logging.Logger) *BaseEngine {
	s.RestServer.WithLogger(logger)
	return s
}

func (s *BaseEngine) WithAdminPort(adminPort uint64) *BaseEngine {
	s.RestServer.WithAdminPort(adminPort)
	return s
}

func (s *BaseEngine) WithApiMux(apiMux rest.Mux) *BaseEngine {
	s.RestServer.WithApiMux(apiMux)
	return s
}

func (s *BaseEngine) WithApiPort(apiPort uint64) *BaseEngine {
	s.RestServer.WithApiPort(apiPort)
	return s
}

func (s *BaseEngine) WithCacheMaximumAge(cacheMaximumAge uint64) *BaseEngine {
	s.RestServer.WithCacheMaximumAge(cacheMaximumAge)
	return s
}

func (s *BaseEngine) WithStatus(status admin.ServiceStatus) *BaseEngine {
	s.RestServer.WithStatus(status)
	return s
}

func (s *BaseEngine) LogHandler() logging.Logger {
	return s.RestServer.Logger
}

func (s *BaseEngine) Run() error {
	s.LogHandler().Info(config.NameAndVersionString() + " -- Starting")
	s.Start()
	return nil
}

func (s *BaseEngine) SetScenario(scenarioFilePath string) {
	s.RestServer.SetScenario(scenarioFilePath)
}

var NullEngine Engine = new(nullEngine)

type nullEngine struct{}

func (s *nullEngine) LogHandler() logging.Logger          { return loggers.NewNullLogger() }
func (s *nullEngine) Run() error                          { return nil }
func (s *nullEngine) SetScenario(scenarioFilePath string) {}
