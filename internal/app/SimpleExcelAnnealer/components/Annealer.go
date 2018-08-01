// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"
	"time"

	. "github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/annealing/logging"
	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	. "github.com/LindsayBradford/crm/logging/handlers"
	"github.com/LindsayBradford/crm/logging/modulators"
)

func BuildObservers(configuration *config.CRMConfig, loggers []LogHandler) []AnnealingObserver {
	observerConfig := configuration.AnnealingObservers
	observerList := make([]AnnealingObserver, len(observerConfig))

	for index, currConfig := range observerConfig {
		var newObserver AnnealingObserver

		var filter modulators.LoggingModulator
		switch currConfig.IterationFilter {
		case "AllowAll":
			filter = new(modulators.NullModulator)
		case "FinishedIterationModulo":
			modulo := (uint64)(currConfig.FilterRate)
			filter = new(modulators.IterationModuloLoggingModulator).
				WithModulo(modulo)
		case "FinishedIterationEveryElapsedSeconds":
			waitAsDuration := (time.Duration)(currConfig.FilterRate) * time.Second
			filter = new(modulators.IterationElapsedTimeLoggingModulator).WithWait(waitAsDuration)
		}

		var observerLogger LogHandler
		for _, logger := range loggers {
			if logger.Name() == currConfig.Logger {
				observerLogger = logger
			}
		}

		switch currConfig.Type {
		case "AttributeObserver":
			newObserver = new(logging.AnnealingAttributeObserver).
				WithLogHandler(observerLogger).
				WithModulator(filter)
		case "MessageObserver":
			newObserver = new(logging.AnnealingMessageObserver).
				WithLogHandler(observerLogger).
				WithModulator(filter)
		}

		observerList[index] = newObserver
	}

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
				WithPenalty(currConfig.Penalty).WithName(currConfig.Name)
		}

		// TODO: I'm throwing away errors... bad.. fix it.
		explorerList[index] = explorer
	}
	return explorerList
}
