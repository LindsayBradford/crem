// +build windows

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	"github.com/LindsayBradford/crm/annealing/shared"
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

	var logHandlers []handlers.LogHandler
	logHandlers, defaultLogHandler = components.BuildLogHandlers(config)

	observers := components.BuildObservers(config, logHandlers)
	explorer := components.BuildSolutionExplorer(config)
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

func profilingRequested() bool {
	return args.CpuProfile != ""
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
