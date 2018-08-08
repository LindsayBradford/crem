// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"errors"
	"time"

	"github.com/LindsayBradford/crm/annealing/logging"
	observers "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/errors"
	"github.com/LindsayBradford/crm/logging/filters"
	"github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealingObserversBuilder struct {
	errors            *CompositeError
	config            []AnnealingObserverConfig
	handlers          []handlers.LogHandler
	maximumIterations uint64
}

func (builder *AnnealingObserversBuilder) initialise() *AnnealingObserversBuilder {
	builder.errors = new(CompositeError)
	return builder
}

func (builder *AnnealingObserversBuilder) WithConfig(crmConfig *CRMConfig) *AnnealingObserversBuilder {
	builder.initialise()
	builder.config = crmConfig.AnnealingObservers
	builder.maximumIterations = crmConfig.Annealer.MaximumIterations
	return builder
}

func (builder *AnnealingObserversBuilder) WithLogHandlers(handlers []handlers.LogHandler) *AnnealingObserversBuilder {
	builder.initialise()
	builder.handlers = handlers
	return builder
}

func (builder *AnnealingObserversBuilder) Build() ([]observers.AnnealingObserver, error) {
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

func (builder *AnnealingObserversBuilder) buildDefaultObservers() []observers.AnnealingObserver {
	defaultObserver := new(logging.AnnealingMessageObserver).
		WithLogHandler(builder.defaultLogger()).
		WithFilter(builder.defaultFilter())
	return []observers.AnnealingObserver{defaultObserver}
}

func (builder *AnnealingObserversBuilder) defaultFilter() *filters.PercentileOfIterationsPerAnnealingFilter {
	filter := new(filters.PercentileOfIterationsPerAnnealingFilter)
	return filter
}

func (builder *AnnealingObserversBuilder) defaultLogger() handlers.LogHandler {
	return builder.handlers[defaultLoggerIndex]
}

func (builder *AnnealingObserversBuilder) buildObservers() []observers.AnnealingObserver {
	observerList := make([]observers.AnnealingObserver, len(builder.config))

	for index, currConfig := range builder.config {
		filter := builder.buildFilter(currConfig)
		logger, loggerError := builder.findLoggerNamedOrDefault(currConfig)

		if loggerError != nil {
			builder.errors.Add(loggerError)
		}

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

func (builder *AnnealingObserversBuilder) findLoggerNamedOrDefault(currConfig AnnealingObserverConfig) (handlers.LogHandler, error) {
	if currConfig.Logger == "" {
		return builder.handlers[defaultLoggerIndex], nil
	}

	for _, logger := range builder.handlers {
		if logger.Name() == currConfig.Logger {
			return logger, nil
		}
	}
	return nil, errors.New("configuration specifies a non-existent logger [\"" + currConfig.Logger + "\"] for an AnnealingObserver")
}

func (builder *AnnealingObserversBuilder) buildFilter(currConfig AnnealingObserverConfig) filters.LoggingFilter {
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
