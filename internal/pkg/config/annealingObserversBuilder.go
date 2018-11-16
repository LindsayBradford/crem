// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"errors"
	"time"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	. "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type annealingObserversBuilder struct {
	errors            *CompositeError
	config            []AnnealingObserverConfig
	handlers          []logging.Logger
	maximumIterations uint64
}

func (builder *annealingObserversBuilder) initialise() *annealingObserversBuilder {
	builder.errors = new(CompositeError)
	return builder
}

func (builder *annealingObserversBuilder) WithConfig(cremConfig *CREMConfig) *annealingObserversBuilder {
	builder.initialise()
	builder.config = cremConfig.AnnealingObservers
	builder.maximumIterations = uint64(cremConfig.Annealer.Parameters[annealers.MaximumIterations].(int64))
	return builder
}

func (builder *annealingObserversBuilder) WithLogHandlers(handlers []logging.Logger) *annealingObserversBuilder {
	builder.initialise()
	builder.handlers = handlers
	return builder
}

func (builder *annealingObserversBuilder) Build() ([]annealing.Observer, error) {
	var observers []annealing.Observer
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

func (builder *annealingObserversBuilder) buildDefaultObservers() []annealing.Observer {
	defaultObserver := new(observer.AnnealingMessageObserver).
		WithLogHandler(builder.defaultLogger()).
		WithFilter(builder.defaultFilter())
	return []annealing.Observer{defaultObserver}
}

func (builder *annealingObserversBuilder) defaultFilter() *filters.PercentileOfIterationsPerAnnealingFilter {
	filter := new(filters.PercentileOfIterationsPerAnnealingFilter)
	return filter
}

func (builder *annealingObserversBuilder) defaultLogger() logging.Logger {
	return builder.handlers[defaultLoggerIndex]
}

func (builder *annealingObserversBuilder) buildObservers() []annealing.Observer {
	observerList := make([]annealing.Observer, len(builder.config))

	for index, currConfig := range builder.config {
		filter := builder.buildFilter(currConfig)
		logger := builder.findLoggerNamedOrDefault(currConfig)
		observerList[index] = buildObserver(currConfig.Type, logger, filter)
	}

	return observerList
}

func buildObserver(observerType AnnealingObserverType, logger logging.Logger, filter filters.Filter) annealing.Observer {
	var newObserver annealing.Observer
	switch observerType {
	case AttributeObserver:
		newObserver = new(observer.AnnealingAttributeObserver).
			WithLogHandler(logger).
			WithFilter(filter)
	case MessageObserver, UnspecifiedAnnealingObserverType:
		newObserver = new(observer.AnnealingMessageObserver).
			WithLogHandler(logger).
			WithFilter(filter)
	default:
		panic("Should not get here")
	}
	return newObserver
}

func (builder *annealingObserversBuilder) findLoggerNamedOrDefault(currConfig AnnealingObserverConfig) logging.Logger {
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
			currConfig.Logger + "\"] for an Observer"),
	)

	return nil
}

func (builder *annealingObserversBuilder) buildFilter(currConfig AnnealingObserverConfig) filters.Filter {
	var filter filters.Filter
	switch currConfig.IterationFilter {
	case UnspecifiedIterationFilter:
		filter = builder.defaultFilter()
	case EveryNumberOfIterations:
		modulo := currConfig.NumberOfIterations
		filter = new(filters.IterationCountFilter).WithModulo(modulo)
	case EveryElapsedSeconds:
		waitAsDuration := (time.Duration)(currConfig.SecondsBetweenEvents) * time.Second
		filter = new(filters.IterationElapsedTimeFilter).WithWait(waitAsDuration)
	case EveryPercentileOfFinishedIterations:
		percentileOfIterations := currConfig.PercentileOfIterations
		filter = new(filters.PercentileOfIterationsPerAnnealingFilter).
			WithPercentileOfIterations(percentileOfIterations).
			WithMaximumIterations(builder.maximumIterations)
	default:
		panic("Should not reach here")
	}
	return filter
}
