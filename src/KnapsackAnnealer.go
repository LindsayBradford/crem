// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"
	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
	. "github.com/LindsayBradford/crm/logging"
	. "github.com/LindsayBradford/crm/logging/modulators"
	. "github.com/LindsayBradford/crm/logging/formatters"
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
		WithLogLevelDestination(DEBUG, STDOUT).
		// WithLogLevelDestination(DEBUG, DISCARD).
		// WithLogLevelDestination(INFO, DISCARD).
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
		// WithLogLevelDestination(INFO, DISCARD).
		WithLogLevelDestination(INFO, STDOUT).
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
	humanAudienceLogger := new(FreeformAnnealingLogger).
		WithLogHandler(humanLogHandler).
		WithModulator(new(NullModulator))
		// WithModulator(new(IterationElapsedTimeLoggingModulator).WithWait(1 * time.Second))
		// WithModulator(new(IterationModuloLoggingModulator).WithModulo(5))

	humanLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(1000).
		WithCoolingFactor(0.995).
		WithMaxIterations(5).
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