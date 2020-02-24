// Copyright (c) 2018 Australian Rivers Institute.

package loggers

import (
	annealingObserver "github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
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
		WithLogLevelDestination(model.LogLevel, logging.STDOUT).
		WithLogLevelDestination(annealingObserver.AnnealingLogLevel, logging.STDOUT).
		Build()

	return testLogger
}

var DefaultTestingAnnealingObserver = buildDefaultTestingAnnealingObserver()

func buildDefaultTestingAnnealingObserver() *annealingObserver.AnnealingMessageObserver {
	filter := new(filters.IterationCountFilter)
	messageObserver := new(annealingObserver.AnnealingMessageObserver).
		WithLogHandler(DefaultTestingLogger).
		WithFilter(filter)
	return messageObserver
}

var DefaultTestingEventNotifier = buildTestingEventNotifier()

func buildTestingEventNotifier() observer.EventNotifier {
	messageObserver := buildDefaultTestingAnnealingObserver()

	eventNotifier := new(observer.SynchronousAnnealingEventNotifier)
	eventNotifier.AddObserver(messageObserver)

	return eventNotifier
}

var NullTestingEventNotifier = buildNullTestingEventNotifier()

func buildNullTestingEventNotifier() observer.EventNotifier {
	filter := new(filters.IterationCountFilter)
	messageObserver := new(annealingObserver.AnnealingMessageObserver).
		WithLogHandler(new(NullLogger)).
		WithFilter(filter)

	eventNotifier := new(observer.SynchronousAnnealingEventNotifier)
	eventNotifier.AddObserver(messageObserver)

	return eventNotifier
}
