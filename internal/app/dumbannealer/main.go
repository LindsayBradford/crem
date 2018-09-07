// (c) 2018 Australian Rivers Institute.
package main

import (
	"os"

	"github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/dumbannealer/components"
	"github.com/pkg/errors"
)

var (
	logger = components.BuildLogHandler()
)

func main() {
	args := commandline.ParseArguments()
	RunFromConfigFile(args.ConfigFile)
}

func RunFromConfigFile(configFile string) {
	configuration := retrieveConfiguration(configFile)
	scenarioRunner := components.BuildScenarioRunner(configuration)
	runScenario(scenarioRunner)
	flushStreams()
}

func retrieveConfiguration(configFile string) *config.CRMConfig {
	configuration, retrieveError := config.RetrieveCrm(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving dumb annealer configuration")
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
