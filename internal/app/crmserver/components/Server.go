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

var (
	ServerLogger handlers.LogHandler = handlers.DefaultNullLogHandler

	crmServerStatus = server.Status{
		Name:    config.ShortApplicationName,
		Version: config.Version,
		Message: "DEAD"}
)

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

func (cs *CrmServer) WithStatus(status server.Status) *CrmServer {
	cs.RestServer.WithStatus(status)
	return cs
}

func RunServerFromConfigFile(configFile string) {
	configuration := retrieveServerConfiguration(configFile)
	establishServerLogger(configuration)

	crmServer := new(CrmServer).
		Initialise().
		WithConfig(configuration).
		WithApiMux(buildCrmApuMux()).
		WithLogger(ServerLogger).
		WithStatus(crmServerStatus)

	ServerLogger.Info(server.NameAndVersionString() + " -- Starting")
	scenario.LogHandler = ServerLogger
	crmServer.Start()
}

func buildCrmApuMux() *api.CrmApiMux {
	return new(api.CrmApiMux).Initialise()
}

func establishServerLogger(configuration *config.HttpServerConfig) {
	loggers, _ := new(config.LogHandlersBuilder).WithConfig(configuration.Loggers).Build()
	ServerLogger = loggers[0]
}

func retrieveServerConfiguration(configFile string) *config.HttpServerConfig {
	configuration, retrieveError := config.RetrieveHttpServer(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving server configuration")
		panic(wrappingError)
	}

	ServerLogger.Info("Configuring with [" + configuration.FilePath + "]")
	return configuration
}
