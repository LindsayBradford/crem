// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components/api"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components/scenario"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/LindsayBradford/crm/server"
	"github.com/pkg/errors"
)

const defaultLoggerIndex = 0

var (
	ServerLogger handlers.LogHandler = handlers.DefaultNullLogHandler

	crmServerStatus = server.ServiceStatus{
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
	return buildCrmServerFrom(serverConfig)
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

func buildCrmServerFrom(serverConfig *config.HttpServerConfig) *server.RestServer {
	return new(CrmServer).
		Initialise().
		WithConfig(serverConfig).
		WithApiMux(buildCrmApuMux(serverConfig)).
		WithLogger(ServerLogger).
		WithStatus(crmServerStatus)
}

func start(crmServer *server.RestServer) {
	ServerLogger.Info(server.NameAndVersionString() + " -- Starting")
	crmServer.Start()
}

type CrmServer struct {
	server.RestServer
}

func (cs *CrmServer) Initialise() *CrmServer {
	cs.RestServer.Initialise()
	return cs
}

func (cs *CrmServer) WithConfig(configuration *config.HttpServerConfig) *CrmServer {
	cs.RestServer.WithConfig(configuration)
	return cs
}

func (cs *CrmServer) WithLogger(logger handlers.LogHandler) *CrmServer {
	cs.RestServer.WithLogger(logger)
	return cs
}

func (cs *CrmServer) WithStatus(status server.ServiceStatus) *CrmServer {
	cs.RestServer.WithStatus(status)
	return cs
}

func buildCrmApuMux(serverConfig *config.HttpServerConfig) *api.CrmApiMux {
	return new(api.CrmApiMux).
		Initialise().
		WithJobQueueLength(serverConfig.JobQueueLength)
}
