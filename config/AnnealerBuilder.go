// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"fmt"

	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/annealing/annealers"
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/errors"
	"github.com/LindsayBradford/crem/logging"
)

type AnnealerBuilder struct {
	config *CREMConfig
	errors *errors.CompositeError

	baseBuilder      annealers.Builder
	loggersBuilder   LogHandlersBuilder
	observersBuilder annealingObserversBuilder
	explorersBuilder solutionExplorerBuilder

	defaultLogHandler logging.Logger
	logHandlers       []logging.Logger
	observers         []annealing.Observer
	annealingExplorer explorer.Explorer
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

func (builder *AnnealerBuilder) WithConfig(suppliedConfig *CREMConfig) *AnnealerBuilder {
	builder.config = suppliedConfig
	builder.initialise()
	return builder
}

func (builder *AnnealerBuilder) RegisteringExplorer(registration ExplorerRegistration) *AnnealerBuilder {
	builder.explorersBuilder.RegisteringExplorer(registration.ExplorerType, registration.ConfigFunction)
	return builder
}

func (builder *AnnealerBuilder) Build() (annealing.Annealer, logging.Logger, error) {
	builder.buildLogHandlers()
	builder.defaultLogHandler.Debug("About to call Builder.Build() ")

	builder.buildObservers()
	builder.buildSolutionExplorer()

	annealerConfig := builder.config.Annealer

	newAnnealer, baseBuildError :=
		builder.buildAnnealerOfType(annealerConfig.Type).
			WithId(builder.config.ScenarioName).
			WithStartingTemperature(annealerConfig.StartingTemperature).
			WithCoolingFactor(annealerConfig.CoolingFactor).
			WithMaxIterations(annealerConfig.MaximumIterations).
			WithLogHandler(builder.defaultLogHandler).
			WithSolutionExplorer(builder.annealingExplorer).
			WithEventNotifier(builder.buildEventNotifier()).
			WithObservers(builder.observers...).
			Build()

	builder.defaultLogHandler.Debug("Call to Builder.Build() finished")

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
		newError := fmt.Errorf("failed to establish log loggers from config: %s", logHandlerErrors.Error())
		builder.errors.Add(newError)
	}

	builder.logHandlers = logHandlers
	builder.setDefaultLogHandler()
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
		newError := fmt.Errorf("failed to establish explorer explorer from config: %s", buildErrors.Error())
		builder.errors.Add(newError)
	}

	builder.annealingExplorer = newExplorer
}

func (builder *AnnealerBuilder) buildAnnealerOfType(annealerType AnnealerType) *annealers.Builder {
	switch annealerType {
	case ElapsedTimeTracking, UnspecifiedAnnealerType:
		return builder.baseBuilder.ElapsedTimeTrackingAnnealer()
	case Simple:
		return builder.baseBuilder.SimpleAnnealer()
	default:
		panic("Should not reach here")
	}
}

func (builder *AnnealerBuilder) buildEventNotifier() annealing.EventNotifier {
	eventNotifierType := builder.config.Annealer.EventNotifier
	switch eventNotifierType {
	case Sequential, UnspecifiedEventNotifierType:
		return new(annealing.SynchronousAnnealingEventNotifier)
	case Concurrent:
		return new(annealing.ConcurrentAnnealingEventNotifier)
	default:
		panic("Should not reach here")
	}
}
