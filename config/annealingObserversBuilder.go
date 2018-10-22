// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"errors"
	"time"

	"github.com/LindsayBradford/crem/annealing/logging"
	observers "github.com/LindsayBradford/crem/annealing/shared"
	. "github.com/LindsayBradford/crem/errors"
	"github.com/LindsayBradford/crem/logging/filters"
	"github.com/LindsayBradford/crem/logging/handlers"
)

type annealingObserversBuilder struct {
	errors            *CompositeError
	config            []AnnealingObserverConfig
	handlers          []handlers.LogHandler
	maximumIterations uint64
}

func (builder *annealingObserversBuilder) initialise() *annealingObserversBuilder {
	builder.errors = new(CompositeError)
	return builder
}

func (builder *annealingObserversBuilder) WithConfig(cremConfig *CREMConfig) *annealingObserversBuilder {
	builder.initialise()
	builder.config = cremConfig.AnnealingObservers
	builder.maximumIterations = cremConfig.Annealer.MaximumIterations
	return builder
}

func (builder *annealingObserversBuilder) WithLogHandlers(handlers []handlers.LogHandler) *annealingObserversBuilder {
	builder.initialise()
	builder.handlers = handlers
	return builder
}

func (builder *annealingObserversBuilder) Build() ([]observers.AnnealingObserver, error) {
	var observers []observers.AnnealingObserver
	if len(builder.config) == 0 {
		observers = builder.buildDefaultObservers()
	} else {
		observers = builder.buildObservers()
	}

	if builder.errors.Size() > 0 {
		return nil, builder.errors
	}
	return observers, nil
}

func (builder *annealingObserversBuilder) buildDefaultObservers() []observers.AnnealingObserver {
	defaultObserver := new(logging.AnnealingMessageObserver).
		WithLogHandler(builder.defaultLogger()).
		WithFilter(builder.defaultFilter())
	return []observers.AnnealingObserver{defaultObserver}
}

func (builder *annealingObserversBuilder) defaultFilter() *filters.PercentileOfIterationsPerAnnealingFilter {
	filter := new(filters.PercentileOfIterationsPerAnnealingFilter)
	return filter
}

func (builder *annealingObserversBuilder) defaultLogger() handlers.LogHandler {
	return builder.handlers[defaultLoggerIndex]
}

func (builder *annealingObserversBuilder) buildObservers() []observers.AnnealingObserver {
	observerList := make([]observers.AnnealingObserver, len(builder.config))

	for index, currConfig := range builder.config {
		filter := builder.buildFilter(currConfig)
		logger := builder.findLoggerNamedOrDefault(currConfig)
		observerList[index] = buildObserver(currConfig.Type, logger, filter)
	}

	return observerList
}

func buildObserver(observerType AnnealingObserverType, logger handlers.LogHandler, filter filters.LoggingFilter) observers.AnnealingObserver {
	var newObserver observers.AnnealingObserver
	switch observerType {
	case AttributeObserver:
		newObserver = new(logging.AnnealingAttributeObserver).
			WithLogHandler(logger).
			WithFilter(filter)
	case MessageObserver, UnspecifiedAnnealingObserverType:
		newObserver = new(logging.AnnealingMessageObserver).
			WithLogHandler(logger).
			WithFilter(filter)
	default:
		panic("Should not get here")
	}
	return newObserver
}

func (builder *annealingObserversBuilder) findLoggerNamedOrDefault(currConfig AnnealingObserverConfig) handlers.LogHandler {
	if currConfig.Logger == "" {
		return builder.handlers[defaultLoggerIndex]
	}

	for _, logger := range builder.handlers {
		if logger.Name() == currConfig.Logger {
			return logger
		}
	}

	builder.errors.Add(
		errors.New("configuration specifies a non-existent logger [\"" +
			currConfig.Logger + "\"] for an AnnealingObserver"),
	)

	return nil
}

func (builder *annealingObserversBuilder) buildFilter(currConfig AnnealingObserverConfig) filters.LoggingFilter {
	var filter filters.LoggingFilter
	switch currConfig.IterationFilter {
	case UnspecifiedIterationFilter:
		filter = builder.defaultFilter()
	case EveryNumberOfIterations:
		modulo := currConfig.NumberOfIterations
		filter = new(filters.IterationCountLoggingFilter).WithModulo(modulo)
	case EveryElapsedSeconds:
		waitAsDuration := (time.Duration)(currConfig.SecondsBetweenEvents) * time.Second
		filter = new(filters.IterationElapsedTimeFilter).WithWait(waitAsDuration)
	case EveryPercentileOfFinishedIterations:
		percentileOfIterations := currConfig.PercentileOfIterations
		filter = new(filters.PercentileOfIterationsPerAnnealingFilter).
			WithPercentileOfIterations(percentileOfIterations).
			WithMaxIterations(builder.maximumIterations)
	default:
		panic("Should not reach here")
	}
	return filter
}
