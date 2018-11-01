// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/annealing"
	. "github.com/LindsayBradford/crem/annealing/explorer"
	cremerrors "github.com/LindsayBradford/crem/errors"
	. "github.com/LindsayBradford/crem/logging/handlers"
)

type Builder struct {
	annealer    annealing.Annealer
	buildErrors *cremerrors.CompositeError
}

func (builder *Builder) ElapsedTimeTrackingAnnealer() *Builder {
	return builder.forAnnealer(&ElapsedTimeTrackingAnnealer{})
}

func (builder *Builder) SimpleAnnealer() *Builder {
	return builder.forAnnealer(&SimpleAnnealer{})
}

func (builder *Builder) forAnnealer(annealer annealing.Annealer) *Builder {
	builder.annealer = annealer
	builder.annealer.Initialise()
	if builder.buildErrors == nil {
		builder.buildErrors = cremerrors.NewComposite("Failed to build valid annealer")
	}
	return builder
}

func (builder *Builder) WithId(title string) *Builder {
	annealerBeingBuilt := builder.annealer
	if title != "" {
		annealerBeingBuilt.SetId(title)
	}
	return builder
}

func (builder *Builder) WithLogHandler(logHandler LogHandler) *Builder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetLogHandler(logHandler); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *Builder) WithStartingTemperature(temperature float64) *Builder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetTemperature(temperature); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *Builder) WithCoolingFactor(coolingFactor float64) *Builder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetCoolingFactor(coolingFactor); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *Builder) WithSolutionExplorer(explorer Explorer) *Builder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetSolutionExplorer(explorer); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *Builder) WithEventNotifier(delegate annealing.EventNotifier) *Builder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetEventNotifier(delegate); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *Builder) WithDumbSolutionExplorer(initialObjectiveValue float64) *Builder {
	annealerBeingBuilt := builder.annealer
	explorer := new(DumbExplorer)
	explorer.SetObjectiveValue(initialObjectiveValue)
	explorer.SetScenarioId(annealerBeingBuilt.Id())
	annealerBeingBuilt.SetSolutionExplorer(explorer)
	return builder
}

func (builder *Builder) WithMaxIterations(iterations uint64) *Builder {
	annealerBeingBuilt := builder.annealer
	annealerBeingBuilt.SetMaxIterations(iterations)
	return builder
}

func (builder *Builder) WithObservers(observers ...annealing.Observer) *Builder {
	annealerBeingBuilt := builder.annealer

	for _, currObserver := range observers {
		if err := annealerBeingBuilt.AddObserver(currObserver); err != nil {
			builder.buildErrors.Add(err)
		}
	}

	return builder
}

func (builder *Builder) Build() (annealing.Annealer, *cremerrors.CompositeError) {
	annealerBeingBuilt := builder.annealer
	buildErrors := builder.buildErrors
	if buildErrors.Size() == 0 {
		return annealerBeingBuilt, nil
	} else {
		return nil, buildErrors
	}
}
