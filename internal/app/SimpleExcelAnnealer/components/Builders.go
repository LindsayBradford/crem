// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"os"

	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/annealing/solution"
	"github.com/LindsayBradford/crem/config"
	"github.com/LindsayBradford/crem/logging/handlers"
)

func BuildScenarioRunner(scenarioConfig *config.CREMConfig, wrapper func(f func()), tearDown func()) (annealing.CallableScenarioRunner, handlers.LogHandler) {
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

	var runner annealing.CallableScenarioRunner

	runner = new(annealing.ScenarioRunner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithTearDownFunction(tearDown).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	runner = new(annealing.SpreadsheetSafeScenarioRunner).ThatLocks(runner)

	if scenarioConfig.CpuProfilePath != "" {
		profilableRunner := new(annealing.ProfilableScenarioRunner).
			ThatProfiles(runner).
			ToFile(scenarioConfig.CpuProfilePath)

		runner = profilableRunner
	}

	return runner, humanLogHandler
}

func buildSimpleExcelExplorerRegistration(wrapper func(f func())) config.ExplorerRegistration {
	return config.ExplorerRegistration{
		ExplorerType: "SimpleExcelSolutionExplorer",
		ConfigFunction: func(config config.SolutionExplorerConfig) solution.Explorer {
			return new(SimpleExcelSolutionExplorer).
				WithPenalty(config.Penalty).
				WithName(config.Name).
				WithInputFile(config.InputFiles["ExcelDataSource"]).
				WithOleFunctionWrapper(wrapper)
		},
	}
}
