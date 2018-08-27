// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"os"

	"github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/logging/handlers"
)

func BuildScenarioRunner(scenarioConfig *config.CRMConfig) (*annealing.ScenarioRunner, handlers.LogHandler) {
	newAnnealer, humanLogHandler, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
			RegisteringExplorer(buildSimpleExcelExplorerRegistration()).
			Build()

	if buildError != nil {
		humanLogHandler.Error(buildError)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	runner := new(annealing.ScenarioRunner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		Concurrently(scenarioConfig.RunConcurrently)

	return runner, humanLogHandler
}

func buildSimpleExcelExplorerRegistration() config.ExplorerRegistration {
	return config.ExplorerRegistration{
		ExplorerType: "SimpleExcelSolutionExplorer",
		ConfigFunction: func(config config.SolutionExplorerConfig) solution.Explorer {
			return new(SimpleExcelSolutionExplorer).
				WithPenalty(config.Penalty).
				WithName(config.Name).
				WithInputFile(config.InputFile)
		},
	}
}
