// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	"errors"
	crmerrors "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealerBuilder struct {
	annealer Annealer
	buildErrors *crmerrors.CompositeError
}

func (this *AnnealerBuilder) ElapsedTimeTrackingAnnealer() *AnnealerBuilder {
	this.annealer = &ElapsedTimeTrackingAnnealer{}
	this.annealer.Initialise()
	this.buildErrors = crmerrors.NewComposite("Failed to build valid elapsed-timed tracking annealer")
	return this
}

func (this *AnnealerBuilder) WithLogHandler(logHandler LogHandler) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.setLogHandler(logHandler); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithStartingTemperature(temperature float64) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.setTemperature(temperature); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithCoolingFactor(coolingFactor float64) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.setCoolingFactor(coolingFactor); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithObjectiveManager(manager ObjectiveManager) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.SetObjectiveManager(manager); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithMaxIterations(iterations uint) *AnnealerBuilder {
	annealer := this.annealer
	annealer.setMaxIterations(iterations)
	return this
}

func (this *AnnealerBuilder) WithObservers(observers ...AnnealingObserver) *AnnealerBuilder {
	if (observers == nil) {
		this.buildErrors.Add(errors.New("Invalid attempt to supply a non-existant observers list to annealer"))
	}

	annealer := this.annealer

	for _, currObserver := range observers {
		if err := annealer.AddObserver(currObserver); err != nil {
			this.buildErrors.Add(err)
		}
	}

	return this
}

func (this *AnnealerBuilder) Build() (Annealer, error) {
	if this.buildErrors.Size() == 0 {
		return this.annealer, nil
	} else {
		return this.annealer, this.buildErrors
	}
}


