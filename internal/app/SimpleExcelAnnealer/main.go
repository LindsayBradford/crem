// +build windows

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	"github.com/LindsayBradford/crm/commandline"
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
	humanAudienceLogger := components.BuildHumanLogger()
	machineAudienceLogger := components.BuildMachineLogger()

	annealer := components.BuildAnnealer(humanAudienceLogger, machineAudienceLogger)

	annealingFunctions.UnProfiledFunction = func() error {
		humanAudienceLogger.Debug("About to call annealer.Anneal()")
		annealer.Anneal()
		humanAudienceLogger.Debug("Call to annealer.Anneal() finished.")
		return nil
	}

	annealingFunctions.ProfiledFunction = func() error {
		humanAudienceLogger.Debug("About to generate cpu profile to file [" + args.CpuProfile + "]")
		defer humanAudienceLogger.Debug("Cpu profiling to file [" + args.CpuProfile + "] now generated")
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
