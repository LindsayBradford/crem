// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
"os"

. "github.com/LindsayBradford/crm/annealing"
.	"github.com/LindsayBradford/crm/annealing/objectives"
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/annealing/logging"
	"github.com/LindsayBradford/crm/commandline"
. "github.com/LindsayBradford/crm/logging/formatters"
. "github.com/LindsayBradford/crm/logging/handlers"
. "github.com/LindsayBradford/crm/logging/modulators"
. "github.com/LindsayBradford/crm/logging/shared"
"github.com/LindsayBradford/crm/profiling"

)

const ERROR_STATUS = 1

var (
	humanLogHandler   LogHandler
	machineLogHandler LogHandler
	annealer          Annealer
)

func init() {
	buildHumanLogger()
	buildMachineLogger()
	buildAnnealer()
}

func buildHumanLogger() {
	logBuilder := new(LogHandlerBuilder)

	newLogger, err := logBuilder.
		ForNativeLibraryLogHandler().
		WithFormatter(new(RawMessageFormatter)).
		// WithLogLevelDestination(DEBUG, STDOUT).
		WithLogLevelDestination(DEBUG, DISCARD).
		// WithLogLevelDestination(INFO, DISCARD).
		WithLogLevelDestination(ANNEALER, STDOUT).
		// WithLogLevelDestination(ANNEALER, DISCARD).
		Build()

	if err != nil {
		humanLogHandler.ErrorWithError(err)
		os.Exit(ERROR_STATUS)
	}
	humanLogHandler = newLogger
}

func buildMachineLogger() {
	logBuilder := new(LogHandlerBuilder)
	newLogger, err := logBuilder.
		ForBareBonesLogHandler().
		WithFormatter(new(JsonFormatter)).
		WithLogLevelDestination(ANNEALER, DISCARD).
		// WithLogLevelDestination(ANNEALER, STDOUT).
		Build()

	if err != nil {
		machineLogHandler.ErrorWithError(err)
		os.Exit(ERROR_STATUS)
	}
	machineLogHandler = newLogger
}

func buildAnnealer() {
	builder := new(AnnealerBuilder)
	machineAudienceObserver := new(AnnealingAttributeObserver).
		WithLogHandler(machineLogHandler).
		WithModulator(new(NullModulator))
	humanAudienceObserver := new(AnnealingMessageObserver).
		WithLogHandler(humanLogHandler).
		// WithModulator(new(NullModulator))
		// WithModulator(new(IterationElapsedTimeLoggingModulator).WithWait(1 * time.Second))
		WithModulator(new(IterationModuloLoggingModulator).WithModulo(200))

	humanLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		ElapsedTimeTrackingAnnealer().
		WithLogHandler(humanLogHandler).
		WithObjectiveManager(new(DumbObjectiveManager)).
		WithStartingTemperature(10).
		WithCoolingFactor(0.997).
		WithMaxIterations(2000).
		WithObservers(machineAudienceObserver, humanAudienceObserver).
		Build()

	humanLogHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if err != nil {
		humanLogHandler.ErrorWithError(err)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(ERROR_STATUS)
	}

	annealer = newAnnealer
}

func main() {
	args := CommandLine.ParseArguments()
	if args.CpuProfile != "" {
		humanLogHandler.Debug("About to generate cpu profile to file [" + args.CpuProfile + "]")
	}
	profiling.CpuProfileOfFunctionToFile(runAnnealer, args.CpuProfile)
	os.Stdout.Sync(); os.Stderr.Sync()  // flush STDOUT & STDERROR streams
}

func runAnnealer() error {
	humanLogHandler.Debug("About to call annealer.Anneal()")
	annealer.Anneal()
	humanLogHandler.Debug("Call to annealer.Anneal() finished. Exiting Program")
	return nil
}