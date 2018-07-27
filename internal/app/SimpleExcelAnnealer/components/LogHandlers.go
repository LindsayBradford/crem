// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing/logging"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
)

func BuildHumanLogger() LogHandler {
	logBuilder := new(LogHandlerBuilder)

	newLogger, err := logBuilder.
		ForNativeLibraryLogHandler().
		WithFormatter(new(RawMessageFormatter)).
		WithLogLevelDestination(DEBUG, STDOUT).
		WithLogLevelDestination(AnnealerLogLevel, STDOUT).
		Build()

	if err != nil {
		newLogger.ErrorWithError(err)
		os.Exit(1)
	}
	return newLogger
}

func BuildMachineLogger() LogHandler {
	logBuilder := new(LogHandlerBuilder)
	newLogger, err := logBuilder.
		ForBareBonesLogHandler().
		WithFormatter(new(JsonFormatter)).
		WithLogLevelDestination(AnnealerLogLevel, DISCARD).
		Build()

	if err != nil {
		newLogger.ErrorWithError(err)
		os.Exit(1)
	}
	return newLogger
}
