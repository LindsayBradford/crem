// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
"os"


. "github.com/LindsayBradford/crm/annealing"
. "github.com/LindsayBradford/crm/logging"
. "github.com/LindsayBradford/crm/logging/formatters"
. "github.com/LindsayBradford/crm/logging/handlers"
. "github.com/LindsayBradford/crm/logging/modulators"
. "github.com/LindsayBradford/crm/logging/shared"


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

	if (err != nil) {
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

	if (err != nil) {
		machineLogHandler.ErrorWithError(err)
		os.Exit(ERROR_STATUS)
	}
	machineLogHandler = newLogger
}

func buildAnnealer() {
	builder := new(AnnealerBuilder)
	machineAudienceLogger := new(AnnealingAttributeLogger).
		WithLogHandler(machineLogHandler).
		WithModulator(new(NullModulator))
	humanAudienceLogger := new(AnnealingMessageLogger).
		WithLogHandler(humanLogHandler).
		WithModulator(new(NullModulator))
		// WithModulator(new(IterationElapsedTimeLoggingModulator).WithWait(1 * time.Second))
		// WithModulator(new(IterationModuloLoggingModulator).WithModulo(5))

	humanLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		ElapsedTimeTrackingAnnealer().
		WithLogHandler(humanLogHandler).
		WithStartingTemperature(10).
		WithCoolingFactor(0.997).
		WithMaxIterations(1000).
		WithObservers(machineAudienceLogger, humanAudienceLogger).
		Build()

	humanLogHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if (err != nil) {
		humanLogHandler.ErrorWithError(err)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(ERROR_STATUS)
	}

	annealer = newAnnealer
}

func main() {
	humanLogHandler.Debug("About to call annealer.Anneal()")
	annealer.Anneal()
	humanLogHandler.Debug("Call to annealer.Anneal() finished. Exiting Program")
}