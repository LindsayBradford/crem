// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealerBuilderViaConfig struct {
	config      AnnealingConfig
	baseBuilder AnnealerBuilder

	logHandler handlers.LogHandler
	explorer   solution.SolutionExplorer
	observers  []AnnealingObserver
}

func (builder *AnnealerBuilderViaConfig) WithConfig(suppliedConfig *CRMConfig) *AnnealerBuilderViaConfig {
	builder.config = suppliedConfig.Annealer
	return builder
}

func (builder *AnnealerBuilderViaConfig) WithLogHandler(handler handlers.LogHandler) *AnnealerBuilderViaConfig {
	builder.logHandler = handler
	return builder
}

func (builder *AnnealerBuilderViaConfig) WithExplorer(explorer solution.SolutionExplorer) *AnnealerBuilderViaConfig {
	builder.explorer = explorer
	return builder
}

func (builder *AnnealerBuilderViaConfig) WithObservers(observers ...AnnealingObserver) *AnnealerBuilderViaConfig {
	builder.observers = observers
	return builder
}

func (builder *AnnealerBuilderViaConfig) Build() (Annealer, error) {
	builder.logHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, baseBuildError :=
		builder.buildAnnealerOfType(builder.config.Type).
			WithStartingTemperature(builder.config.StartingTemperature).
			WithCoolingFactor(builder.config.CoolingFactor).
			WithMaxIterations(builder.config.MaximumIterations).
			WithLogHandler(builder.logHandler).
			WithSolutionExplorer(builder.explorer).
			WithEventNotifier(builder.buildEventNotifier(builder.config.EventNotifier)).
			WithObservers(builder.observers...).
			Build()

	builder.logHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if baseBuildError != nil {
		return nil, baseBuildError
	}
	return newAnnealer, nil
}

func (builder *AnnealerBuilderViaConfig) buildEventNotifier(eventNotifierType EventNotifierType) AnnealingEventNotifier {
	switch eventNotifierType {
	case Sequential, UnspecifiedEventNotifierType:
		return new(SynchronousAnnealingEventNotifier)
	case Concurrent:
		return new(ConcurrentAnnealingEventNotifier)
	default:
		panic("Should not reach here")
	}
}

func (builder *AnnealerBuilderViaConfig) buildAnnealerOfType(annealerType AnnealerType) *AnnealerBuilder {
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
