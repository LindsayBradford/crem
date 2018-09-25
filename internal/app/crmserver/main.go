// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/crmserver/server"
	"github.com/LindsayBradford/crm/internal/app/dumbannealer/components"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/pkg/errors"
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

var (
	logger = components.BuildLogHandler() //TODO: picked up wrong logger -- fix

	crmServer = new(CrmServer).
			Initialise().
			WithLogger(logger).
			WithStatus(server.Status{Name: config.ShortApplicationName, Version: config.Version, Message: "DEAD"})
)

func main() {
	args := commandline.ParseArguments()
	RunFromConfigFile(args.ConfigFile)
}

func RunFromConfigFile(configFile string) {
	logger.Info(nameAndVersionString() + " -- Started")

	configuration := retrieveConfiguration(configFile)
	crmServer.
		WithConfig(configuration).
		establishHttpDelegates()

	crmServer.Start()
}

func retrieveConfiguration(configFile string) *config.HttpServerConfig {
	configuration, retrieveError := config.RetrieveHttpServer(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving server configuration")
		panic(wrappingError)
	}

	logger.Info("Configuring with [" + configuration.FilePath + "]")
	return configuration
}

func (cs *CrmServer) establishHttpDelegates() {
	cs.Logger.Debug("Registering API handlers")
	cs.AddApiMapping("/", rootPathHandler)
}

func rootPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, nameAndVersionString())
}

func nameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", config.LongApplicationName, config.Version)
}
