// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/config"
	"github.com/LindsayBradford/crem/cmd/cremengine/engine"
	"github.com/LindsayBradford/crem/cmd/cremengine/engine/api"
	data2 "github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/config/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/server/admin"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

var (
	ServerLogger = loggers.DefaultTestingLogger

	engineStatus = admin.ServiceStatus{
		ServiceName: config.ShortApplicationName,
		Version:     config.Version,
		Status:      "DEAD"}
)

type EngineConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	loggingInterpreter *interpreter.LoggingConfigInterpreter

	engine engine.Engine
	logger logging.Logger
}

func NewEngineConfigInterpreter() *EngineConfigInterpreter {
	interpreter := new(EngineConfigInterpreter).initialise()
	return interpreter
}

func (i *EngineConfigInterpreter) initialise() *EngineConfigInterpreter {
	i.errors = compositeErrors.New("Scenario Configuration")
	i.loggingInterpreter = interpreter.NewLoggingConfigInterpreter()
	i.engine = engine.NullEngine
	return i
}

func (i *EngineConfigInterpreter) Interpret(engineConfig data2.HttpServerConfig) *EngineConfigInterpreter {
	i.buildLogger(engineConfig)
	i.buildEngine(engineConfig)
	return i
}

func (i *EngineConfigInterpreter) buildLogger(engineConfig data2.HttpServerConfig) {
	ServerLogger = i.loggingInterpreter.Interpret(&engineConfig.Logger).LogHandler()
}

func (i *EngineConfigInterpreter) buildEngine(engineConfig data2.HttpServerConfig) {
	apiMux := buildApiMux(engineConfig)

	i.engine = engine.NewBaseEngine().
		WithApiPort(engineConfig.ApiPort).
		WithAdminPort(engineConfig.AdminPort).
		WithApiMux(apiMux).
		WithCacheMaximumAge(engineConfig.CacheMaximumAgeInSeconds).
		WithLogHandler(ServerLogger).
		WithStatus(engineStatus)

	// TODO: Job Queue Length?
}

func buildApiMux(serverConfig data2.HttpServerConfig) *api.Mux {
	return new(api.Mux).Initialise()
}

func (i *EngineConfigInterpreter) Engine() engine.Engine {
	return i.engine
}

func (i *EngineConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
