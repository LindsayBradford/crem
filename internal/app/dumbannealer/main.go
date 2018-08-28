// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/dumbannealer/components"
)

var (
	args = commandline.ParseArguments()
)

func main() {
	configuration, retrieveError := config.Retrieve(args.ConfigFile)

	if retrieveError != nil {
		commandline.Exit(retrieveError)
	}

	logger := components.BuildLogHandler()

	logger.Info("Configuring with [" + configuration.FilePath + "]")
	scenarioRunner := components.BuildScenarioRunner(configuration)

	runError := scenarioRunner.Run()

	if runError != nil {
		commandline.Exit(runError)
	}

	defer flushStreams()
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
