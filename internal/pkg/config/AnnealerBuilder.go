// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"fmt"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	errors2 "github.com/pkg/errors"
)

type AnnealerBuilder struct {
	config *CREMConfig
	errors *errors.CompositeError

	baseBuilder      annealers.Builder
	loggersBuilder   LogHandlersBuilder
	observersBuilder annealingObserversBuilder
	explorersBuilder solutionExplorerBuilder
	modelsBuilder    modelBuilder

	defaultLogHandler logging.Logger
	logHandlers       []logging.Logger
	observers         []annealing.Observer
	annealingExplorer explorer.Explorer
	annealingModel    model.Model
}

type ExplorerRegistration struct {
	ExplorerType   string
	ConfigFunction ExplorerConfigFunction
}

type ModelRegistration struct {
	ModelType      string
	ConfigFunction ModelConfigFunction
}

func (builder *AnnealerBuilder) initialise() {
	builder.errors = errors.New("AnnealerBuilder initialisation")
	builder.loggersBuilder.WithConfig(builder.config.Loggers)
	builder.observersBuilder.WithConfig(builder.config)
	builder.explorersBuilder.WithConfig(builder.config)
	builder.modelsBuilder.WithConfig(builder.config)
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

func (builder *AnnealerBuilder) RegisteringModel(registration ModelRegistration) *AnnealerBuilder {
	builder.modelsBuilder.RegisteringModel(registration.ModelType, registration.ConfigFunction)
	return builder
}

func (builder *AnnealerBuilder) Build() (annealing.Annealer, logging.Logger, error) {
	builder.buildLogHandlers()
	builder.defaultLogHandler.Debug("About to call Builder.Build() ")

	builder.buildObservers()
	builder.buildModel()
	builder.buildSolutionExplorer()

	annealerConfig := builder.config.Annealer

	newAnnealer, baseBuildError :=
		builder.buildAnnealerOfType(annealerConfig.Type).
			WithId(builder.config.ScenarioName).
			WithParameters(annealerConfig.Parameters).
			WithLogHandler(builder.defaultLogHandler).
			WithSolutionExplorer(builder.annealingExplorer).
			WithEventNotifier(builder.buildEventNotifier()).
			WithObservers(builder.observers...).
			Build()

	builder.defaultLogHandler.Debug("Call to Builder.Build() finished")

	if baseBuildError != nil {
		newError := errors2.Wrap(baseBuildError, "failed to establish annealer from config")
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

func (builder *AnnealerBuilder) buildModel() {
	myModelName := builder.config.Annealer.Model
	newModel, buildErrors := builder.modelsBuilder.Build(myModelName)

	if loggingModel, isLogger := newModel.(logging.Container); isLogger {
		loggingModel.SetLogHandler(builder.defaultLogHandler)
	}

	if buildErrors != nil {
		newError := errors2.Wrap(buildErrors, "building model from config")
		builder.errors.Add(newError)
	}

	builder.annealingModel = newModel
}

func (builder *AnnealerBuilder) buildSolutionExplorer() {
	myExplorerName := builder.config.Annealer.SolutionExplorer
	newExplorer, buildErrors := builder.explorersBuilder.Build(myExplorerName)

	if buildErrors != nil {
		newError := errors2.Wrap(buildErrors, "building explorer from config")
		builder.errors.Add(newError)
	} else {
		newExplorer.SetModel(builder.annealingModel)
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
