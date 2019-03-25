// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"os"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

func RunServerScenarioFromConfig(cremConfig *config.CREMConfig) error {
	scenarioRunner, runnerError := BuildScenarioRunner(cremConfig)
	if runnerError != nil {
		return runnerError
	}

	go runScenario(scenarioRunner)
	return nil
}

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
		WithSaver(new(scenario.Saver)).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	return runner, nil
}

func buildCatchmentModelRegistration() config.ModelRegistration {
	return config.ModelRegistration{
		ModelType: "CatchmentModel",
		ConfigFunction: func(config config.ModelConfig) model.Model {
			return NewModel().
				WithName(config.Name).
				WithOleFunctionWrapper(threading.GetMainThreadChannel().Call).
				WithParameters(config.Parameters)
		},
	}
}

func runScenario(scenarioRunner scenario.CallableRunner) {
	if runError := scenarioRunner.Run(); runError != nil {
		wrappingError := errors.Wrap(runError, "running scenario")
		LogHandler.Error(wrappingError)
	}
	flushStreams()
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
