// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	"errors"
	crmerrors "github.com/LindsayBradford/crm/errors"
)

type AnnealerBuilder struct {
	annealer Annealer
	buildErrors *crmerrors.CompositeError
}

func (this *AnnealerBuilder) SingleObjectiveAnnealer() *AnnealerBuilder {
	this.annealer = &singleObjectiveAnnealer{}
	this.annealer.Initialise()
	this.buildErrors = crmerrors.New("AnnealerBuilder Error")
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

func (this *AnnealerBuilder) WithMaxIterations(iterations uint) *AnnealerBuilder {
	annealer := this.annealer
	annealer.setMaxIterations(iterations)
	return this
}

func (this *AnnealerBuilder) WithObservers(observers ...AnnealingObserver) *AnnealerBuilder {
	if (observers == nil) {
		this.buildErrors.Add(errors.New("Attempt to assign an observers list to Annealer"))
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


