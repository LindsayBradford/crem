// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/handlers"
. "github.com/LindsayBradford/crm/logging/shared"

)

const ERROR_STATUS = 1

var (
	log LogHandler
	annealer Annealer
)

func init() {
	buildLogger()
	buildAnnealer()
}

func buildLogger() {
	logBuilder := new(LogHandlerBuilder)
	logFormatter := new(JsonFormatter)
	newLogger, err := logBuilder.
		// ForNativeLibraryLogHandler().
		ForBareBonesLogHandler().
		WithFormatter(logFormatter).
		WithLogLevelDestination(DEBUG, STDOUT).
		// WithLogLevelDestination(ERROR, STDOUT).
		Build()

	if (err != nil) {
		log.ErrorWithError(err)
		os.Exit(ERROR_STATUS)
	}
	log = newLogger
}

func buildAnnealer() {
	builder := new(AnnealerBuilder)
	// annealerLogger := new(JsonMessageAnnealingLogger).WithLogHandler(log)
	annealerLogger := new(FreeformAnnealingLogger).WithLogHandler(log)

	log.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(1000).
		WithCoolingFactor(0.995).
		WithMaxIterations(5).
		WithObservers(annealerLogger).
		Build()

	log.Debug("Call to AnnealerBuilder.Build() finished")

	if (err != nil) {
		log.ErrorWithError(err)
		log.Error("Exiting program due to failed Annealer build")
		os.Exit(ERROR_STATUS)
	}

	annealer = newAnnealer
}

func main() {
	log.Debug("About to call annealer.Anneal()")
	annealer.Anneal()
	log.Debug("Call to annealer.Anneal() finished. Exiting Program")
}
