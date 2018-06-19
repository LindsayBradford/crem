// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging"
)

const ERROR_STATUS = 1

var (
	globalLogHandler LogHandler
)

func init() {
	logBuilder := new(LogHandlerBuilder)
	logFormatter := new(MessageFormatter)
	logHandler, err := logBuilder.
		ForNativeLibraryLogHandler().
		WithFormatter(logFormatter).
		WithLogLevelDestination(DEBUG, STDOUT).
		Build()

	if (err != nil) {
		globalLogHandler.ErrorWithError(err)
		os.Exit(ERROR_STATUS)
	}

	globalLogHandler = logHandler
}

func main() {
	builder := new(AnnealerBuilder)

	annealer, err := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(1000).
		WithCoolingFactor(0.995).
		WithMaxIterations(5).
		WithObservers(
			new(JsonMessageAnnealingLogger).WithLogHandler(globalLogHandler)).
		Build()

	if (err != nil) {
		globalLogHandler.ErrorWithError(err)
		globalLogHandler.Error("Exiting due to failed Annealer build")
		os.Exit(ERROR_STATUS)
	}

	annealer.Anneal()
}
