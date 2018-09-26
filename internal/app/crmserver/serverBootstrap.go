// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/config"
	"github.com/pkg/errors"
)

func runServerFromConfigFile(configFile string) {
	logger.Info(nameAndVersionString() + " -- Started")

	configuration := retrieveServerConfiguration(configFile)
	crmServer.
		WithConfig(configuration).
		establishHttpDelegates()

	crmServer.Start()
}

func retrieveServerConfiguration(configFile string) *config.HttpServerConfig {
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
