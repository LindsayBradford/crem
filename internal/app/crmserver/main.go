// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components"
	"github.com/LindsayBradford/crm/internal/app/crmserver/server"
	"github.com/LindsayBradford/crm/logging/handlers"
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
		runServerFromConfigFile(args.ServerConfigFile)
	}
}

func shouldRunScenario(args *commandline.Arguments) bool {
	return args.ScenarioFile != ""
}
