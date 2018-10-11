// Copyright (c) 2018 Australian Rivers Institute.

package handlers

import (
	"github.com/LindsayBradford/crm/logging/formatters"
	"github.com/LindsayBradford/crm/logging/shared"
)

var DefaultTestingLogHandler = buildTestingLogHandler()

func buildTestingLogHandler() LogHandler {
	builder := new(LogHandlerBuilder)
	testLogger, _ := builder.
		ForNativeLibraryLogHandler().
		WithName("DefaultTestingLogHandler").
		WithFormatter(new(formatters.RawMessageFormatter)).
		WithLogLevelDestination(shared.DEBUG, shared.STDOUT).
		Build()
	return testLogger
}
