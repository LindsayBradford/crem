// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"fmt"

	"github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/errors"
	"github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealerBuilder struct {
	config *CRMConfig
	errors *errors.CompositeError

	baseBuilder      annealing.AnnealerBuilder
	loggersBuilder   logHandlersBuilder
	observersBuilder annealingObserversBuilder
	explorersBuilder solutionExplorerBuilder

	defaultLogHandler handlers.LogHandler
	logHandlers       []handlers.LogHandler
	observers         []AnnealingObserver
	annealingExplorer solution.Explorer
}

type ExplorerRegistration struct {
	ExplorerType   string
	ConfigFunction ExplorerConfigFunction
}

func (builder *AnnealerBuilder) initialise() {
	builder.errors = new(errors.CompositeError)
	builder.loggersBuilder.WithConfig(builder.config.Loggers)
	builder.observersBuilder.WithConfig(builder.config)
	builder.explorersBuilder.WithConfig(builder.config)
}

func (builder *AnnealerBuilder) WithConfig(suppliedConfig *CRMConfig) *AnnealerBuilder {
	builder.config = suppliedConfig
	builder.initialise()
	return builder
}

func (builder *AnnealerBuilder) RegisteringExplorer(registration ExplorerRegistration) *AnnealerBuilder {
	builder.explorersBuilder.RegisteringExplorer(registration.ExplorerType, registration.ConfigFunction)
	return builder
}

func (builder *AnnealerBuilder) Build() (Annealer, handlers.LogHandler, error) {
	builder.buildLogHandlers()
	builder.defaultLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	builder.buildObservers()
	builder.buildSolutionExplorer()

	annealerConfig := builder.config.Annealer

	newAnnealer, baseBuildError :=
		builder.buildAnnealerOfType(annealerConfig.Type).
			WithStartingTemperature(annealerConfig.StartingTemperature).
			WithCoolingFactor(annealerConfig.CoolingFactor).
			WithMaxIterations(annealerConfig.MaximumIterations).
			WithLogHandler(builder.defaultLogHandler).
			WithSolutionExplorer(builder.annealingExplorer).
			WithEventNotifier(builder.buildEventNotifier()).
			WithObservers(builder.observers...).
			Build()

	builder.defaultLogHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if baseBuildError != nil {
		newError := fmt.Errorf("failed to establish annealer from config: %s", baseBuildError.Error())
		builder.errors.Add(newError)
	}

	if builder.errors.Size() > 0 {
		return nil, builder.defaultLogHandler, builder.errors
	}

	return newAnnealer, builder.defaultLogHandler, nil
}

func (builder *AnnealerBuilder) buildLogHandlers() {
	logHandlers, logHandlerErrors := builder.loggersBuilder.Build()

	if logHandlerErrors != nil {
		newError := fmt.Errorf("failed to establish log handlers from config: %s", logHandlerErrors.Error())
		builder.errors.Add(newError)
	}

	builder.logHandlers = logHandlers
	builder.setDefaultLogHandler()

	defer func() {
		builder.defaultLogHandler.Info("Configuring with [" + builder.config.FilePath + "]")
	}()
}

func (builder *AnnealerBuilder) setDefaultLogHandler() {
	builder.defaultLogHandler = builder.logHandlers[0]
}

func (builder *AnnealerBuilder) buildObservers() {
	configuredObservers, observerErrors :=
		new(annealingObserversBuilder).
			WithConfig(builder.config).
			WithLogHandlers(builder.logHandlers).
			Build()

	if observerErrors != nil {
		newError := fmt.Errorf("failed to establish annealing observes from config: %s", observerErrors.Error())
		builder.errors.Add(newError)
	}

	builder.observers = configuredObservers
}

func (builder *AnnealerBuilder) buildSolutionExplorer() {
	myExplorerName := builder.config.Annealer.SolutionExplorer
	newExplorer, buildErrors := builder.explorersBuilder.Build(myExplorerName)

	if buildErrors != nil {
		newError := fmt.Errorf("failed to establish solution explorer from config: %s", buildErrors.Error())
		builder.errors.Add(newError)
	}

	builder.annealingExplorer = newExplorer
}

func (builder *AnnealerBuilder) buildAnnealerOfType(annealerType AnnealerType) *annealing.AnnealerBuilder {
	switch annealerType {
	case ElapsedTimeTracking, UnspecifiedAnnealerType:
		return builder.baseBuilder.ElapsedTimeTrackingAnnealer()
	case OSThreadLocked:
		return builder.baseBuilder.OSThreadLockedAnnealer()
	case Simple:
		return builder.baseBuilder.SimpleAnnealer()
	default:
		panic("Should not reach here")
	}
}

func (builder *AnnealerBuilder) buildEventNotifier() AnnealingEventNotifier {
	eventNotifierType := builder.config.Annealer.EventNotifier
	switch eventNotifierType {
	case Sequential, UnspecifiedEventNotifierType:
		return new(SynchronousAnnealingEventNotifier)
	case Concurrent:
		return new(ConcurrentAnnealingEventNotifier)
	default:
		panic("Should not reach here")
	}
}
