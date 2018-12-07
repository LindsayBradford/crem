// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"os"

	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/pkg/errors"
)

var (
	LogHandler logging.Logger
)

func RunScenarioFromConfig(cremConfig *config.CREMConfig) error {
	scenarioRunner, runnerError := BuildScenarioRunner(cremConfig)
	if runnerError != nil {
		return runnerError
	}

	runScenario(scenarioRunner)
	flushStreams()
	return nil
}

func BuildScenarioRunner(scenarioConfig *config.CREMConfig) (scenario.CallableRunner, error) {
	newAnnealer, _, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
			RegisteringModel(buildCatchmentModelRegistration()).
			Build()

	if buildError != nil {
		LogHandler.Error(buildError)
		return nil, buildError
	}

	var runner scenario.CallableRunner

	runner = new(scenario.Runner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	return runner, nil
}

func buildCatchmentModelRegistration() config.ModelRegistration {
	return config.ModelRegistration{
		ModelType: "CatchmentModel",
		ConfigFunction: func(config config.ModelConfig) model.Model {
			return NewCatchmentModel().
				WithName(config.Name).
				WithParameters(config.Parameters)
		},
	}
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
