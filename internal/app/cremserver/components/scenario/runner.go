// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"os"

	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/config"
	"github.com/LindsayBradford/crem/logging/handlers"
	"github.com/pkg/errors"
)

var (
	LogHandler handlers.LogHandler
)

func RunScenarioFromConfig(cremConfig *config.CREMConfig) {
	scenarioRunner := BuildScenarioRunner(cremConfig)
	runScenario(scenarioRunner)
	flushStreams()
}

func BuildScenarioRunner(scenarioConfig *config.CREMConfig) annealing.CallableScenarioRunner {
	newAnnealer, _, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
			Build()

	if buildError != nil {
		LogHandler.Error(buildError)
	}

	var runner annealing.CallableScenarioRunner

	runner = new(annealing.ScenarioRunner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	return runner
}

func runScenario(scenarioRunner annealing.CallableScenarioRunner) {
	if runError := scenarioRunner.Run(); runError != nil {
		wrappingError := errors.Wrap(runError, "running scenario")
		LogHandler.Error(wrappingError)
	}
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
