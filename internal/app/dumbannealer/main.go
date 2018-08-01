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
	configuration := config.Retrieve(args.ConfigFile)

	logger := components.BuildLogHandler()

	annealingFunctions.UnProfiledFunction = func() error {
		logger.Info("Configuring with [" + configuration.FilePath + "]")
		annealer := components.BuildDumbAnnealer(configuration, logger)
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
