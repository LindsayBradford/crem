// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components"
	"github.com/LindsayBradford/crm/internal/app/crmserver/server"
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
	logger = components.BuildLogHandler()

	crmServer = new(CrmServer).
			Initialise().
			WithLogger(logger).
			WithStatus(server.Status{Name: config.ShortApplicationName, Version: config.Version, Message: "DEAD"})
)

func main() {
	args := commandline.ParseArguments()
	if shouldRunScenario(args) {
		RunScenarioFromConfigFile(args.ScenarioFile)
	} else {
		RunServerFromConfigFile(args.ServerConfigFile)
	}
}

func shouldRunScenario(args *commandline.Arguments) bool {
	return args.ScenarioFile != ""
}

func RunScenarioFromConfigFile(configFile string) {
	configuration := retrieveScenarioConfiguration(configFile)
	scenarioRunner := components.BuildScenarioRunner(configuration)
	runScenario(scenarioRunner)
	flushStreams()
}

func retrieveScenarioConfiguration(configFile string) *config.CRMConfig {
	configuration, retrieveError := config.RetrieveCrm(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		panic(wrappingError)
	}

	logger.Info("Configuring with [" + configuration.FilePath + "]")
	return configuration
}

func runScenario(scenarioRunner annealing.CallableScenarioRunner) {
	if runError := scenarioRunner.Run(); runError != nil {
		wrappingError := errors.Wrap(runError, "running dumb annealer scenario")
		panic(wrappingError)
	}
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}

func RunServerFromConfigFile(configFile string) {
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
