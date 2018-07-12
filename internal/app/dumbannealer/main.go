// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/internal/app/dumbannealer/components"
	"github.com/LindsayBradford/crm/profiling"
)

var (
	args             = commandline.ParseArguments()
	annealingRunners = new(profiling.ProfiledAndUnProfiledFunctionPair)
)

func main() {
	buildAnnealingRunners()

	if profilingRequested() {
		annealingRunners.ProfiledFunction()
	} else {
		annealingRunners.UnProfiledFunction()
	}

	defer flushStreams()
}

func buildAnnealingRunners() {
	logger := components.BuildLogHandler()

	annealingRunners.UnProfiledFunction = func() error {
		annealer := components.BuildDumbAnnealer(logger)
		logger.Debug("About to call annealer.Anneal()")
		annealer.Anneal()
		logger.Debug("Call to annealer.Anneal() finished.")
		return nil
	}

	annealingRunners.ProfiledFunction = func() error {
		logger.Debug("About to generate cpu profile to file [" + args.CpuProfile + "]")
		return profiling.CpuProfileOfFunctionToFile(annealingRunners.UnProfiledFunction, args.CpuProfile)
	}
}

func profilingRequested() bool {
	return args.CpuProfile != ""
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
