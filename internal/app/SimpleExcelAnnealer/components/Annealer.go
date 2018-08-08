// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	. "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

func BuildLogHandlers(loggingConfig []config.LoggerConfig) ([]LogHandler, error) {
	handlerList, buildError :=
		new(config.LogHandlersBuilder).
			WithConfig(loggingConfig).
			Build()

	compositeError, ok := buildError.(*CompositeError)
	if ok && compositeError.Size() > 0 {
		return nil, compositeError
	}

	return handlerList, nil
}

func BuildObservers(configuration *config.CRMConfig, loggers []LogHandler) []AnnealingObserver {
	observerList :=
		new(config.AnnealingObserversBuilder).
			WithConfig(configuration).
			WithLogHandlers(loggers).
			Build()
	return observerList
}

func BuildAnnealer(configuration *config.CRMConfig, humanLogHandler LogHandler, observers ...AnnealingObserver) Annealer {
	builder := new(AnnealerBuilder)

	humanLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	annealerConfig := configuration.Annealer

	solutionExplorers := BuildSolutionExplorers(configuration)
	mySolutionExplorer := findMyExplorer(solutionExplorers, annealerConfig)

	newAnnealer, err := builder.
		AnnealerOfType(annealerConfig.Type).
		WithStartingTemperature(annealerConfig.StartingTemperature).
		WithCoolingFactor(annealerConfig.CoolingFactor).
		WithMaxIterations(annealerConfig.MaximumIterations).
		WithLogHandler(humanLogHandler).
		WithSolutionExplorer(mySolutionExplorer).
		WithEventNotifier(annealerConfig.EventNotifier).
		WithObservers(observers...).
		Build()

	humanLogHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if err != nil {
		humanLogHandler.ErrorWithError(err)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}

func findMyExplorer(solutionExplorers []solution.SolutionExplorer, annealerConfig config.AnnealingConfig) solution.SolutionExplorer {
	var mySolutionExplorer solution.SolutionExplorer
	for _, explorer := range solutionExplorers {
		if annealerConfig.SolutionExplorer == explorer.Name() {
			mySolutionExplorer = explorer
		}
	}
	return mySolutionExplorer
}

func BuildSolutionExplorers(configuration *config.CRMConfig) []solution.SolutionExplorer {
	explorerConfig := configuration.SolutionExplorers

	explorerList := make([]solution.SolutionExplorer, len(explorerConfig))
	for index, currConfig := range explorerConfig {
		var explorer solution.SolutionExplorer

		// TODO: Works for this example, but generally speaking, this explorer is local, whereas the other two
		// TODO: are "pre-canned"... how to handle gracefully going forward?

		switch currConfig.Type {
		case "NullSolutionExplorer", "":
			explorer = new(solution.NullSolutionExplorer).
				WithName(currConfig.Name)
		case "DumbSolutionExplorer":
			explorer = new(solution.DumbSolutionExplorer).
				WithName(currConfig.Name)
		case "SimpleExcelSolutionExplorer":
			explorer = new(SimpleExcelSolutionExplorer).
				WithPenalty(currConfig.Penalty).
				WithName(currConfig.Name)
		}

		// TODO: I'm throwing away errors... bad.. fix it.
		explorerList[index] = explorer
	}
	return explorerList
}
