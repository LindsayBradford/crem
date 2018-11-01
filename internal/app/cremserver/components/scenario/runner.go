// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"os"

	"github.com/LindsayBradford/crem/config"
	"github.com/LindsayBradford/crem/logging"
	"github.com/LindsayBradford/crem/scenario"
	"github.com/pkg/errors"
)

var (
	LogHandler logging.Logger
)

func RunScenarioFromConfig(cremConfig *config.CREMConfig) {
	scenarioRunner := BuildScenarioRunner(cremConfig)
	runScenario(scenarioRunner)
	flushStreams()
}

func BuildScenarioRunner(scenarioConfig *config.CREMConfig) scenario.CallableRunner {
	newAnnealer, _, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
			Build()

	if buildError != nil {
		LogHandler.Error(buildError)
	}

	var runner scenario.CallableRunner

	runner = new(scenario.Runner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	return runner
}

func runScenario(scenarioRunner scenario.CallableRunner) {
	if runError := scenarioRunner.Run(); runError != nil {
		wrappingError := errors.Wrap(runError, "running scenario")
		LogHandler.Error(wrappingError)
	}
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
