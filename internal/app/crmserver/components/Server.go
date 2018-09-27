// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/config"
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
	crmServer := new(CrmServer).
		Initialise().
		WithLogger(ServerLogger).
		WithStatus(crmServerStatus)

	ServerLogger.Info(nameAndVersionString() + " -- Started")

	configuration := retrieveServerConfiguration(configFile)
	crmServer.
		WithConfig(configuration).
		establishApiHandlers()

	crmServer.Start()
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

func (cs *CrmServer) establishApiHandlers() {
	cs.Logger.Debug("Registering API handlers")
	cs.AddApiMapping("/", rootPathHandler)
}

func rootPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, nameAndVersionString())
}

func nameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", config.LongApplicationName, config.Version)
}
