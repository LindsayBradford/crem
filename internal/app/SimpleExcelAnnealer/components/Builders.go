// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"os"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
)

func BuildScenarioRunner(scenarioConfig *config.CREMConfig, wrapper func(f func()), tearDown func()) (scenario.CallableRunner, logging.Logger) {
	newAnnealer, humanLogHandler, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
			RegisteringExplorer(buildSimpleExcelExplorerRegistration(wrapper)).
			Build()

	if buildError != nil {
		humanLogHandler.Error(buildError)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	var runner scenario.CallableRunner

	runner = new(scenario.Runner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithTearDownFunction(tearDown).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	runner = new(scenario.SpreadsheetSafeScenarioRunner).ThatLocks(runner)

	if scenarioConfig.CpuProfilePath != "" {
		profilableRunner := new(scenario.ProfilableRunner).
			ThatProfiles(runner).
			ToFile(scenarioConfig.CpuProfilePath)

		runner = profilableRunner
	}

	return runner, humanLogHandler
}

func buildSimpleExcelExplorerRegistration(wrapper func(f func())) config.ExplorerRegistration {
	return config.ExplorerRegistration{
		ExplorerType: "SimpleExcelSolutionExplorer",
		ConfigFunction: func(config config.SolutionExplorerConfig) explorer.Explorer {
			return new(SimpleExcelSolutionExplorer).
				WithParameters(config.Parameters).
				WithName(config.Name).
				WithOleFunctionWrapper(wrapper)
		},
	}
}
