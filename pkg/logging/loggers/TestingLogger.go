// Copyright (c) 2018 Australian Rivers Institute.

package loggers

import (
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/formatters"
)

var DefaultTestingLogger = buildTestingLogger()

func buildTestingLogger() logging.Logger {
	builder := new(Builder)
	testLogger, _ := builder.
		ForNativeLibraryLogHandler().
		WithName("DefaultTestingLogger").
		WithFormatter(new(formatters.RawMessageFormatter)).
		WithLogLevelDestination(logging.DEBUG, logging.STDOUT).
		Build()
	return testLogger
}
