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
	configuration := retrieveServerConfiguration(configFile)
	establishServerLogger(configuration)

	crmServer := new(CrmServer).
		Initialise().
		WithConfig(configuration).
		WithApiMux(buildCrmApuMix()).
		WithLogger(ServerLogger).
		WithStatus(crmServerStatus)

	ServerLogger.Info(nameAndVersionString() + " -- Starting")
	crmServer.Start()
}

func buildCrmApuMix() *CrmApiMux {
	newMux := new(CrmApiMux).Initialise().WithType("API)")
	return newMux
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

type CrmApiMux struct {
	server.ApiMux
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	return cam
}

func (cam *CrmApiMux) WithType(muxType string) *CrmApiMux {
	cam.ApiMux.WithType(muxType)
	return cam
}

func (cam *CrmApiMux) rootPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, nameAndVersionString())
}

func nameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", config.LongApplicationName, config.Version)
}
