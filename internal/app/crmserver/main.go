// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components"
)

func main() {
	buildSharedLogger()
	decideBetweenScenarioAndServer()
}

func buildSharedLogger() {
	logger := components.BuildLogHandler()

	components.ScenarioLogger = logger
	components.ServerLogger = logger
}

func decideBetweenScenarioAndServer() {
	args := commandline.ParseArguments()
	if shouldRunScenario(args) {
		components.RunScenarioFromConfigFile(args.ScenarioFile)
	} else {
		components.RunServerFromConfigFile(args.ServerConfigFile)
	}
}

func shouldRunScenario(args *commandline.Arguments) bool {
	return args.ScenarioFile != ""
}
