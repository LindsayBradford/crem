// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"fmt"
	"os"

	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/SimpleExcelAnnealer/components"
	"github.com/LindsayBradford/crm/logging/handlers"
)

func BuildDumbAnnealer(annealerConfig *config.CRMConfig, logHandler handlers.LogHandler) Annealer {

	logHandlers, defaultLogHandler := components.BuildLogHandlers(annealerConfig)
	observers := BuildObservers(annealerConfig, logHandlers)
	explorer := BuildSolutionExplorer(annealerConfig)

	newAnnealer, buildError :=
		new(config.AnnealerBuilderViaConfig).
			WithConfig(annealerConfig).
			WithLogHandler(defaultLogHandler).
			WithObservers(observers...).
			WithExplorer(explorer).
			Build()

	if buildError != nil {
		logHandler.ErrorWithError(buildError)
		logHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}

func BuildSolutionExplorer(configuration *config.CRMConfig) solution.SolutionExplorer {
	myExplorerName := configuration.Annealer.SolutionExplorer

	explorer, buildErrors :=
		new(config.SolutionExplorerBuilder).
			WithConfig(configuration).
			Build(myExplorerName)

	if buildErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish solution explorer from config: %s", buildErrors.Error())
		panic(panicMsg)
	}
	return explorer
}

func BuildObservers(configuration *config.CRMConfig, logHandlers []handlers.LogHandler) []AnnealingObserver {
	observers, observerErrors :=
		new(config.AnnealingObserversBuilder).
			WithConfig(configuration).
			WithLogHandlers(logHandlers).
			Build()
	if observerErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish annealing observes from config: %s", observerErrors.Error())
		panic(panicMsg)
	}
	return observers
}
