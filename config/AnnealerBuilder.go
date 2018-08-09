// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"fmt"

	"github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealerBuilder struct {
	config *CRMConfig

	baseBuilder      annealing.AnnealerBuilder
	loggersBuilder   LogHandlersBuilder
	observersBuilder AnnealingObserversBuilder
	explorersBuilder SolutionExplorerBuilder

	defaultLogHandler handlers.LogHandler
	logHandlers       []handlers.LogHandler
	observers         []AnnealingObserver
	annealingExplorer solution.SolutionExplorer
}

type ExplorerRegistration struct {
	ExplorerType   string
	ConfigFunction ExplorerConfigFunction
}

func (builder *AnnealerBuilder) initialise() {
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
		return nil, builder.defaultLogHandler, baseBuildError
	}
	return newAnnealer, builder.defaultLogHandler, nil
}

func (builder *AnnealerBuilder) buildLogHandlers() {
	logHandlers, logHandlerErrors := builder.loggersBuilder.Build()

	if logHandlerErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish log handlers from config: %s", logHandlerErrors.Error())
		panic(panicMsg)
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
		new(AnnealingObserversBuilder).
			WithConfig(builder.config).
			WithLogHandlers(builder.logHandlers).
			Build()

	if observerErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish annealing observes from config: %s", observerErrors.Error())
		panic(panicMsg)
	}

	builder.observers = configuredObservers
}

func (builder *AnnealerBuilder) buildSolutionExplorer() {
	myExplorerName := builder.config.Annealer.SolutionExplorer
	newExplorer, buildErrors := builder.explorersBuilder.Build(myExplorerName)

	if buildErrors != nil {
		panicMsg := fmt.Sprintf("failed to establish solution newExplorer from config: %s", buildErrors.Error())
		panic(panicMsg)
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
