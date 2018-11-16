// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components/api"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/server"
	"github.com/LindsayBradford/crem/internal/pkg/server/admin"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/pkg/errors"
)

const defaultLoggerIndex = 0

var (
	ServerLogger logging.Logger = loggers.DefaultNullLogLogger

	cremServerStatus = admin.ServiceStatus{
		ServiceName: config.ShortApplicationName,
		Version:     config.Version,
		Status:      "DEAD"}
)

func RunServerFromConfigFile(configFile string) {
	configuredServer := buildServerFromConfigFile(configFile)
	start(configuredServer)
}

func buildServerFromConfigFile(configFile string) *server.RestServer {
	serverConfig := retrieveServerConfiguration(configFile)
	buildLoggerFrom(serverConfig)
	return buildServerFrom(serverConfig)
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

func buildServerFrom(serverConfig *config.HttpServerConfig) *server.RestServer {
	return new(Server).
		Initialise().
		WithConfig(serverConfig).
		WithApiMux(buildApiMux(serverConfig)).
		WithLogger(ServerLogger).
		WithStatus(cremServerStatus)
}

func buildApiMux(serverConfig *config.HttpServerConfig) *api.Mux {
	return new(api.Mux).
		Initialise().
		WithJobQueueLength(serverConfig.JobQueueLength)
}

func start(server *server.RestServer) {
	ServerLogger.Info(rest.NameAndVersionString() + " -- Starting")
	server.Start()
}
