// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/dumbannealer/components"
	"github.com/LindsayBradford/crm/profiling"
)

var (
	args               = commandline.ParseArguments()
	annealingFunctions = new(profiling.OptionalProfilingFunctionPair)
)

func main() {
	buildAnnealingRunners()

	var runError error
	if profilingRequested() {
		runError = annealingFunctions.ProfiledFunction()
	} else {
		runError = annealingFunctions.UnProfiledFunction()
	}

	if runError != nil {
		commandline.ExitWithError(runError)
	}

	defer flushStreams()
}

func buildAnnealingRunners() {
	configuration, retrieveError := config.Retrieve(args.ConfigFile)

	if retrieveError != nil {
		commandline.ExitWithError(retrieveError)
	}

	logger := components.BuildLogHandler()

	annealingFunctions.UnProfiledFunction = func() error {
		logger.Info("Configuring with [" + configuration.FilePath + "]")
		annealer := components.BuildDumbAnnealer(configuration)
		logger.Debug("About to call annealer.Anneal()")
		annealer.Anneal()
		logger.Debug("Call to annealer.Anneal() finished.")
		return nil
	}

	annealingFunctions.ProfiledFunction = func() error {
		logger.Debug("About to generate cpu profile to file [" + args.CpuProfile + "]")
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
