// Copyright (c) 2018 Australian Rivers Institute.

package bootstrap

import (
	"os"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/pkg/errors"
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
			Build()

	if buildError != nil {
		LogHandler.Error(buildError)
		return nil, buildError
	}

	saver := scenario.NewSaver().
		WithOutputType(configOutputTypeToSolutionOutputType(scenarioConfig.OutputType)).
		WithOutputPath(scenarioConfig.OutputPath)

	var runner scenario.CallableRunner

	runner = new(scenario.Runner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithSaver(saver).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	if scenarioConfig.CpuProfilePath != "" {
		profilableRunner := new(scenario.ProfilingRunner).
			ThatProfiles(runner).
			ToFile(scenarioConfig.CpuProfilePath)

		runner = profilableRunner
	}

	return runner, nil
}

func configOutputTypeToSolutionOutputType(outputType config.ScenarioOutputType) encoding.OutputType {
	return encoding.OutputType(outputType.String())
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
