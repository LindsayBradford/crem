// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"fmt"
	"os"

	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/logging/handlers"
)

func BuildLogHandlers(configuration *config.CRMConfig) ([]handlers.LogHandler, handlers.LogHandler) {
	logHandlers, logHandlerErrors :=
		new(config.LogHandlersBuilder).
			WithConfig(configuration.Loggers).
			Build()
	if logHandlerErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish log handlers from config: %s", logHandlerErrors.Error())
		panic(panicMsg)
	}

	defaultLogHandler := logHandlers[0]

	defer func() {
		defaultLogHandler.Info("Configuring with [" + configuration.FilePath + "]")
	}()

	return logHandlers, defaultLogHandler
}

func BuildObservers(configuration *config.CRMConfig, logHandlers []handlers.LogHandler) []shared.AnnealingObserver {
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

func BuildSolutionExplorer(configuration *config.CRMConfig) solution.SolutionExplorer {
	myExplorerName := configuration.Annealer.SolutionExplorer

	explorer, buildErrors :=
		new(config.SolutionExplorerBuilder).
			WithConfig(configuration).
			RegisteringExplorer(
				"SimpleExcelSolutionExplorer",
				func(config config.SolutionExplorerConfig) solution.SolutionExplorer {
					return new(SimpleExcelSolutionExplorer).
						WithPenalty(config.Penalty).
						WithName(config.Name)
				},
			).Build(myExplorerName)

	if buildErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish solution explorer from config: %s", buildErrors.Error())
		panic(panicMsg)
	}
	return explorer
}

func BuildAnnealer(configuration *config.CRMConfig, humanLogHandler handlers.LogHandler, explorer solution.SolutionExplorer, observers ...shared.AnnealingObserver) shared.Annealer {
	newAnnealer, buildError :=
		new(config.AnnealerBuilderViaConfig).
			WithConfig(configuration).
			WithLogHandler(humanLogHandler).
			WithExplorer(explorer).
			WithObservers(observers...).
			Build()

	if buildError != nil {
		humanLogHandler.ErrorWithError(buildError)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}
