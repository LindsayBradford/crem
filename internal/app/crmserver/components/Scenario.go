// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"os"

	"github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/pkg/errors"
)

var (
	ScenarioLogger handlers.LogHandler = handlers.DefaultNullLogHandler
)

func RunScenarioFromConfigFile(configFile string) {
	configuration := retrieveScenarioConfiguration(configFile)
	scenarioRunner := BuildScenarioRunner(configuration)
	runScenario(scenarioRunner)
	flushStreams()
}

func retrieveScenarioConfiguration(configFile string) *config.CRMConfig {
	configuration, retrieveError := config.RetrieveCrm(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		panic(wrappingError)
	}

	ScenarioLogger.Info("Configuring with [" + configuration.FilePath + "]")
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
