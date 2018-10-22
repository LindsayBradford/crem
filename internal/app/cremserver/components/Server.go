// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/config"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components/api"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components/scenario"
	"github.com/LindsayBradford/crem/logging/handlers"
	"github.com/LindsayBradford/crem/server"
	"github.com/pkg/errors"
)

const defaultLoggerIndex = 0

var (
	ServerLogger handlers.LogHandler = handlers.DefaultNullLogHandler

	cremServerStatus = server.ServiceStatus{
		ServiceName: config.ShortApplicationName,
		Version:     config.Version,
		Status:      "DEAD"}
)

func RunServerFromConfigFile(configFile string) {
	configuredServer := buildServerFromFrom(configFile)
	start(configuredServer)
}

func buildServerFromFrom(configFile string) *server.RestServer {
	serverConfig := retrieveServerConfiguration(configFile)
	buildLoggerFrom(serverConfig)
	return buildCremServerFrom(serverConfig)
}

func retrieveServerConfiguration(configFile string) *config.HttpServerConfig {
	configuration, retrieveError := config.RetrieveHttpServer(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving server configuration")
		panic(wrappingError)
	}

	return configuration
}

func buildLoggerFrom(configuration *config.HttpServerConfig) {
	establishServerLogger(configuration)
	scenario.LogHandler = ServerLogger
}

func establishServerLogger(configuration *config.HttpServerConfig) {
	loggers, _ := new(config.LogHandlersBuilder).WithConfig(configuration.Loggers).Build()
	ServerLogger = loggers[defaultLoggerIndex]
	ServerLogger.Info("Configuring with [" + configuration.FilePath + "]")
}

func buildCremServerFrom(serverConfig *config.HttpServerConfig) *server.RestServer {
	return new(CremServer).
		Initialise().
		WithConfig(serverConfig).
		WithApiMux(buildCremApuMux(serverConfig)).
		WithLogger(ServerLogger).
		WithStatus(cremServerStatus)
}

func start(cremServer *server.RestServer) {
	ServerLogger.Info(server.NameAndVersionString() + " -- Starting")
	cremServer.Start()
}

type CremServer struct {
	server.RestServer
}

func (cs *CremServer) Initialise() *CremServer {
	cs.RestServer.Initialise()
	return cs
}

func (cs *CremServer) WithConfig(configuration *config.HttpServerConfig) *CremServer {
	cs.RestServer.WithConfig(configuration)
	return cs
}

func (cs *CremServer) WithLogger(logger handlers.LogHandler) *CremServer {
	cs.RestServer.WithLogger(logger)
	return cs
}

func (cs *CremServer) WithStatus(status server.ServiceStatus) *CremServer {
	cs.RestServer.WithStatus(status)
	return cs
}

func buildCremApuMux(serverConfig *config.HttpServerConfig) *api.CremApiMux {
	return new(api.CremApiMux).
		Initialise().
		WithJobQueueLength(serverConfig.JobQueueLength)
}
