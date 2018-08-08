// +build windows

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"fmt"
	"os"

	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/SimpleExcelAnnealer/components"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/LindsayBradford/crm/profiling"
)

var (
	args               = commandline.ParseArguments()
	annealingFunctions = new(profiling.ProfiledAndUnProfiledFunctionPair)
	defaultLogHandler  handlers.LogHandler
)

func main() {
	buildAnnealingFunctions()

	if profilingRequested() {
		annealingFunctions.ProfiledFunction()
	} else {
		annealingFunctions.UnProfiledFunction()
	}

	defer flushStreams()
}

func buildAnnealingFunctions() {
	annealer := buildAnnealerOffConfig()

	annealingFunctions.UnProfiledFunction = func() error {
		defaultLogHandler.Debug("About to call annealer.Anneal()")
		annealer.Anneal()
		defaultLogHandler.Debug("Call to annealer.Anneal() finished.")
		return nil
	}

	annealingFunctions.ProfiledFunction = func() error {
		defaultLogHandler.Debug("About to generate cpu profile to file [" + args.CpuProfile + "]")
		defer defaultLogHandler.Debug("Cpu profiling to file [" + args.CpuProfile + "] now generated")
		return profiling.CpuProfileOfFunctionToFile(annealingFunctions.UnProfiledFunction, args.CpuProfile)
	}
}

func buildAnnealerOffConfig() shared.Annealer {
	config := retrieveConfig()
	logHandlers := buildLogHandlers(config)
	observers := buildObservers(config, logHandlers)
	explorer := buildSolutionExplorer(config)
	annealer := components.BuildAnnealer(config, defaultLogHandler, explorer, observers...)
	return annealer
}

func retrieveConfig() *config.CRMConfig {
	configuration, retrieveError := config.Retrieve(args.ConfigFile)
	if retrieveError != nil {
		panic(retrieveError)
	}
	return configuration
}

func buildLogHandlers(configuration *config.CRMConfig) []handlers.LogHandler {
	logHandlers, logHandlerErrors :=
		new(config.LogHandlersBuilder).
			WithConfig(configuration.Loggers).
			Build()
	if logHandlerErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish log handlers from config: %s", logHandlerErrors.Error())
		panic(panicMsg)
	}

	defer func() {
		defaultLogHandler = logHandlers[0]
		defaultLogHandler.Info("Configuring with [" + configuration.FilePath + "]")
	}()

	return logHandlers
}

func buildObservers(configuration *config.CRMConfig, logHandlers []handlers.LogHandler) []shared.AnnealingObserver {
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

func buildSolutionExplorer(configuration *config.CRMConfig) solution.SolutionExplorer {
	myExplorerName := configuration.Annealer.SolutionExplorer

	explorer, buildErrors :=
		new(config.SolutionExplorerBuilder).
			WithConfig(configuration).
			RegisteringExplorer(
				"SimpleExcelSolutionExplorer",
				func(config config.SolutionExplorerConfig) solution.SolutionExplorer {
					return new(components.SimpleExcelSolutionExplorer).
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

func profilingRequested() bool {
	return args.CpuProfile != ""
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
