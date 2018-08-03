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
	"github.com/LindsayBradford/crm/logging/filters"
	. "github.com/LindsayBradford/crm/logging/handlers"
	"github.com/LindsayBradford/crm/logging/shared"
)

func BuildObservers(configuration *config.CRMConfig, loggers []LogHandler) []AnnealingObserver {
	if len(configuration.AnnealingObservers) == 0 {
		return buildDefaultObservers()
	} else {
		return buildObservers(configuration, loggers)
	}
}

func buildDefaultObservers() []AnnealingObserver {
	defaultLogger := buildDefaultLogger()
	defaultFilter := buildDefaultFilter()

	defaultObserver := new(logging.AnnealingMessageObserver).
		WithLogHandler(defaultLogger).
		WithFilter(defaultFilter)

	return []AnnealingObserver{defaultObserver}
}

func buildDefaultFilter() *filters.PercentileOfIterationsPerAnnealingFilter {
	filter := new(filters.PercentileOfIterationsPerAnnealingFilter)
	return filter
}

func buildDefaultLogger() LogHandler {
	logBuilder := new(LogHandlerBuilder)
	defaultLogHandler, _ := logBuilder.
		ForDefaultLogHandler().
		WithLogLevelDestination(logging.AnnealerLogLevel, shared.STDOUT).
		Build()
	return defaultLogHandler
}

func buildObservers(configuration *config.CRMConfig, loggers []LogHandler) []AnnealingObserver {
	observerConfig := configuration.AnnealingObservers
	observerList := make([]AnnealingObserver, len(observerConfig))

	for index, currConfig := range observerConfig {
		filter := buildFilter(currConfig, configuration)
		logger := findLoggerNamed(loggers, currConfig.Logger)

		observerList[index] = buildObserver(currConfig.Type, logger, filter)
	}

	return observerList
}

func buildObserver(observerType config.AnnealingObserverType, logger LogHandler, filter filters.LoggingFilter) AnnealingObserver {
	var newObserver AnnealingObserver
	switch observerType {
	case config.AttributeObserver:
		newObserver = new(logging.AnnealingAttributeObserver).
			WithLogHandler(logger).
			WithFilter(filter)
	case config.MessageObserver, config.UnspecifiedAnnealingObserverType:
		newObserver = new(logging.AnnealingMessageObserver).
			WithLogHandler(logger).
			WithFilter(filter)
	default:
		panic("Should not get here")
	}
	return newObserver
}

func findLoggerNamed(loggers []LogHandler, name string) LogHandler {
	var observerLogger LogHandler
	for _, logger := range loggers {
		if logger.Name() == name {
			observerLogger = logger
		}
	}
	return observerLogger
}

func buildFilter(currConfig config.AnnealingObserverConfig, configuration *config.CRMConfig) filters.LoggingFilter {
	var filter filters.LoggingFilter
	switch currConfig.IterationFilter {
	case config.UnspecifiedIterationFilter:
		filter = new(filters.PercentileOfIterationsPerAnnealingFilter).
			WithPercentileOfIterations(100).
			WithMaxIterations(configuration.Annealer.MaximumIterations)
	case config.EveryNumberOfIterations:
		modulo := currConfig.NumberOfIterations
		filter = new(filters.IterationCountLoggingFilter).WithModulo(modulo)
	case config.EveryElapsedSeconds:
		waitAsDuration := (time.Duration)(currConfig.SecondsBetweenEvents) * time.Second
		filter = new(filters.IterationElapsedTimeFilter).WithWait(waitAsDuration)
	case config.EveryPercentileOfFinishedIterations:
		percentileOfIterations := currConfig.PercentileOfIterations
		filter = new(filters.PercentileOfIterationsPerAnnealingFilter).
			WithPercentileOfIterations(percentileOfIterations).
			WithMaxIterations(configuration.Annealer.MaximumIterations)
	default:
		panic("Should not reach here")
	}
	return filter
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
