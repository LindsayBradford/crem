// +build windows

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"fmt"
	"os"

	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/SimpleExcelAnnealer/components"
	"github.com/LindsayBradford/crm/profiling"
)

var (
	args               = commandline.ParseArguments()
	annealingFunctions = new(profiling.ProfiledAndUnProfiledFunctionPair)
)

func main() {
	buildAnnealingRunners()

	if profilingRequested() {
		annealingFunctions.ProfiledFunction()
	} else {
		annealingFunctions.UnProfiledFunction()
	}

	defer flushStreams()
}

func buildAnnealingRunners() {
	configuration, retrieveError := config.Retrieve(args.ConfigFile)

	if retrieveError != nil {
		panic(retrieveError)
	}

	logHandlers, logHandlerErrors := components.BuildLogHandlers(configuration.Loggers)

	if logHandlerErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish log handlers from config: %s", logHandlerErrors.Error())
		panic(panicMsg)
	}

	defaultLogHandler := logHandlers[0]
	defaultLogHandler.Info("Configuring with [" + configuration.FilePath + "]")

	observers := components.BuildObservers(configuration, logHandlers)
	annealer := components.BuildAnnealer(configuration, defaultLogHandler, observers...)

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

func profilingRequested() bool {
	return args.CpuProfile != ""
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
