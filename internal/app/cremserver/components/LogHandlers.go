// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	. "github.com/LindsayBradford/crem/annealing/logging"
	. "github.com/LindsayBradford/crem/logging/formatters"
	. "github.com/LindsayBradford/crem/logging/handlers"
	. "github.com/LindsayBradford/crem/logging/shared"
)

func BuildLogHandler() LogHandler {
	logBuilder := new(LogHandlerBuilder)

	newLogger, err := logBuilder.
		ForNativeLibraryLogHandler().
		WithFormatter(new(RawMessageFormatter)).
		WithLogLevelDestination(DEBUG, STDOUT).
		WithLogLevelDestination(AnnealerLogLevel, STDOUT).
		Build()

	if err != nil {
		newLogger.Error(err)
		os.Exit(1)
	}
	return newLogger
}
