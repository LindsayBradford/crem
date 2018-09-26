// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	"github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/config"
)

func BuildScenarioRunner(scenarioConfig *config.CRMConfig) annealing.CallableScenarioRunner {
	newAnnealer, logHandler, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
			Build()

	if buildError != nil {
		logHandler.Error(buildError)
		logHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	var runner annealing.CallableScenarioRunner

	runner = new(annealing.ScenarioRunner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	if scenarioConfig.CpuProfilePath != "" {
		runner = new(annealing.ProfilableScenarioRunner).
			ThatProfiles(runner).
			ToFile(scenarioConfig.CpuProfilePath)
	}

	return runner
}
