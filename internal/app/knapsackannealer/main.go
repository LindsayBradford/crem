// +build windows

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/internal/app/knapsackannealer/components"
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
	humanAudienceLogger := components.BuildHumanLogger()
	machineAudienceLogger := components.BuildMachineLogger()

	annealer := components.BuildAnnealer(humanAudienceLogger, machineAudienceLogger)

	annealingRunners.UnProfiledFunction = func() error {
		humanAudienceLogger.Debug("About to call annealer.Anneal()")
		annealer.Anneal()
		humanAudienceLogger.Debug("Call to annealer.Anneal() finished.")
		return nil
	}

	annealingRunners.ProfiledFunction = func() error {
		humanAudienceLogger.Debug("About to generate cpu profile to file [" + args.CpuProfile + "]")
		defer humanAudienceLogger.Debug("Cpu profiling to file [" + args.CpuProfile + "] now generated")
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
