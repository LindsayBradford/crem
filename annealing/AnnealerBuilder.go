// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	. "github.com/LindsayBradford/crem/annealing/shared"
	. "github.com/LindsayBradford/crem/annealing/solution"
	cremerrors "github.com/LindsayBradford/crem/errors"
	. "github.com/LindsayBradford/crem/logging/handlers"
)

type AnnealerBuilder struct {
	annealer    Annealer
	buildErrors *cremerrors.CompositeError
}

func (builder *AnnealerBuilder) ElapsedTimeTrackingAnnealer() *AnnealerBuilder {
	return builder.forAnnealer(&ElapsedTimeTrackingAnnealer{})
}

func (builder *AnnealerBuilder) SimpleAnnealer() *AnnealerBuilder {
	return builder.forAnnealer(&SimpleAnnealer{})
}

func (builder *AnnealerBuilder) forAnnealer(annealer Annealer) *AnnealerBuilder {
	builder.annealer = annealer
	builder.annealer.Initialise()
	if builder.buildErrors == nil {
		builder.buildErrors = cremerrors.NewComposite("Failed to build valid annealer")
	}
	return builder
}

func (builder *AnnealerBuilder) WithId(title string) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if title != "" {
		annealerBeingBuilt.SetId(title)
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

func (builder *AnnealerBuilder) WithSolutionExplorer(explorer Explorer) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetSolutionExplorer(explorer); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *AnnealerBuilder) WithEventNotifier(delegate AnnealingEventNotifier) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	if err := annealerBeingBuilt.SetEventNotifier(delegate); err != nil {
		builder.buildErrors.Add(err)
	}
	return builder
}

func (builder *AnnealerBuilder) WithDumbSolutionExplorer(initialObjectiveValue float64) *AnnealerBuilder {
	annealerBeingBuilt := builder.annealer
	explorer := new(DumbExplorer)
	explorer.SetObjectiveValue(initialObjectiveValue)
	explorer.SetScenarioId(annealerBeingBuilt.Id())
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

func (builder *AnnealerBuilder) Build() (Annealer, *cremerrors.CompositeError) {
	annealerBeingBuilt := builder.annealer
	buildErrors := builder.buildErrors
	if buildErrors.Size() == 0 {
		return annealerBeingBuilt, nil
	} else {
		return nil, buildErrors
	}
}
