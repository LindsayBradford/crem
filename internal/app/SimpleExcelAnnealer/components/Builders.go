// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"os"

	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/logging/handlers"
)

func BuildAnnealer(configuration *config.CRMConfig) (shared.Annealer, handlers.LogHandler) {
	newAnnealer, humanLogHandler, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(configuration).
			RegisteringExplorer(buildSimpleExcelExplorerRegistration()).
			Build()

	if buildError != nil {
		humanLogHandler.ErrorWithError(buildError)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer, humanLogHandler
}

func buildSimpleExcelExplorerRegistration() config.ExplorerRegistration {
	return config.ExplorerRegistration{
		ExplorerType: "SimpleExcelSolutionExplorer",
		ConfigFunction: func(config config.SolutionExplorerConfig) solution.SolutionExplorer {
			return new(SimpleExcelSolutionExplorer).
				WithPenalty(config.Penalty).
				WithName(config.Name).
				WithInputFile(config.InputFile)
		},
	}
}
