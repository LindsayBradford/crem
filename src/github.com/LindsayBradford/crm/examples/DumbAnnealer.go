// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

. "github.com/LindsayBradford/crm/annealing"
. "github.com/LindsayBradford/crm/annealing/logging"
. "github.com/LindsayBradford/crm/annealing/shared"
"github.com/LindsayBradford/crm/commandline"
. "github.com/LindsayBradford/crm/logging/formatters"
. "github.com/LindsayBradford/crm/logging/handlers"
. "github.com/LindsayBradford/crm/logging/modulators"
. "github.com/LindsayBradford/crm/logging/shared"
"github.com/LindsayBradford/crm/profiling"
)

func buildDumbAnnealerLogger() LogHandler {
	logBuilder := new(LogHandlerBuilder)

	newLogger, err := logBuilder.
		ForNativeLibraryLogHandler().
		WithFormatter(new(RawMessageFormatter)).
		WithLogLevelDestination(DEBUG, STDOUT).
		WithLogLevelDestination(ANNEALER, STDOUT).
		Build()

	if err != nil {
		newLogger.ErrorWithError(err)
		os.Exit(1)
	}
	return newLogger
}

func buildDumbAnnealer(logHandler LogHandler) Annealer {
	builder := new(AnnealerBuilder)
	humanAudienceObserver := new(AnnealingMessageObserver).
		WithLogHandler(logHandler).
		WithModulator(
			new(IterationModuloLoggingModulator).WithModulo(1))  // No STARTED_ITERATION events, all FINISHED_ITERATION events

	logHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		ElapsedTimeTrackingAnnealer().
		WithDumbObjectiveManager(100).
		WithLogHandler(logHandler).
		WithStartingTemperature(10).
		WithCoolingFactor(0.99).
		WithMaxIterations(500).
		WithObservers(humanAudienceObserver).
		Build()

	logHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if err != nil {
		logHandler.ErrorWithError(err)
		logHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}

func main() {
	logger := buildDumbAnnealerLogger()

	runAnnealer := func() error {
		annealer := buildDumbAnnealer(logger)
		logger.Debug("About to call annealer.Anneal()")
		annealer.Anneal()
		logger.Debug("Call to annealer.Anneal() finished. Exiting Program")
		return nil
	}

	args := CommandLine.ParseArguments()
	if args.CpuProfile != "" {
		logger.Debug("About to generate cpu profile to file [" + args.CpuProfile + "]")
	}

	profiling.CpuProfileOfFunctionToFile(runAnnealer, args.CpuProfile)
	os.Stdout.Sync(); os.Stderr.Sync()  // flush STDOUT & STDERROR streams
}

