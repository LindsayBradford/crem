// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	crmerrors "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealerBuilder struct {
	annealer    Annealer
	buildErrors *crmerrors.CompositeError
}

func (builder *AnnealerBuilder) AnnealerOfType(annealerType config.AnnealerType) *AnnealerBuilder {

	builder.buildErrors = crmerrors.NewComposite("failed to build configured annealer")

	switch annealerType {
	case config.ElapsedTimeTracking, config.UnspecifiedAnnealerType:
		return builder.ElapsedTimeTrackingAnnealer()
	case config.OSThreadLocked:
		return builder.OSThreadLockedAnnealer()
	case config.Simple:
		return builder.SimpleAnnealer()
	default:
		panic("Should not reach here")
	}
	return builder
}

func (builder *AnnealerBuilder) OSThreadLockedAnnealer() *AnnealerBuilder {
	builder.annealer = &OSThreadLockedAnnealer{}
	builder.annealer.Initialise()
	if builder.buildErrors == nil {
		builder.buildErrors = crmerrors.NewComposite("Failed to build valid OS thread-locked annealer")
	}
	return builder
}

func (builder *AnnealerBuilder) ElapsedTimeTrackingAnnealer() *AnnealerBuilder {
	builder.annealer = &ElapsedTimeTrackingAnnealer{}
	builder.annealer.Initialise()
	if builder.buildErrors == nil {
		builder.buildErrors = crmerrors.NewComposite("Failed to build valid elapsed-timed tracking annealer")
	}
	return builder
}

func (builder *AnnealerBuilder) SimpleAnnealer() *AnnealerBuilder {
	builder.annealer = &SimpleAnnealer{}
	builder.annealer.Initialise()
	if builder.buildErrors == nil {
		builder.buildErrors = crmerrors.NewComposite("Failed to build valid simple annealer")
	}
	return builder
}

func (builder *AnnealerBuilder) WithLogHandler(logHandler LogHandler) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetLogHandler(logHandler); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *AnnealerBuilder) WithStartingTemperature(temperature float64) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetTemperature(temperature); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *AnnealerBuilder) WithCoolingFactor(coolingFactor float64) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetCoolingFactor(coolingFactor); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *AnnealerBuilder) WithSolutionExplorer(explorer SolutionExplorer) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetSolutionExplorer(explorer); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *AnnealerBuilder) WithEventNotifier(eventNotifierType config.EventNotifierType) *AnnealerBuilder {
	switch eventNotifierType {
	case config.Sequential, config.UnspecifiedEventNotifierType:
		return builder.withEventNotifier(new(SynchronousAnnealingEventNotifier))
	case config.Concurrent:
		return builder.withEventNotifier(new(ConcurrentAnnealingEventNotifier))
	default:
		panic("Should not reach here")
	}
	return builder
}

func (builder *AnnealerBuilder) withEventNotifier(delegate AnnealingEventNotifier) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetEventNotifier(delegate); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *AnnealerBuilder) WithDumbSolutionExplorer(initialObjectiveValue float64) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	explorer := new(DumbSolutionExplorer)
	explorer.SetObjectiveValue(initialObjectiveValue)
	annealerBeingBuilt.SetSolutionExplorer(explorer)
	return builder
}

func (builder *AnnealerBuilder) WithMaxIterations(iterations uint64) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	annealerBeingBuilt.SetMaxIterations(iterations)
	return builder
}

func (builder *AnnealerBuilder) WithObservers(observers ...AnnealingObserver) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer

	for _, currObserver := range observers {
		if err := annealerBeingBuilt.AddObserver(currObserver); err != nil {
			builder.buildErrors.Add(err)
		}
	}

	return builder
}

func (builder *AnnealerBuilder) Build() (Annealer, *crmerrors.CompositeError) {
	annealerBeingBuilt := builder.annealer
	buildErrors := builder.buildErrors
	if buildErrors.Size() == 0 {
		return annealerBeingBuilt, nil
	} else {
		return nil, buildErrors
	}
}
