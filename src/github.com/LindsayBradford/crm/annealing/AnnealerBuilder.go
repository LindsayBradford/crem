// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	"errors"
	crmerrors "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/annealing/objectives"
	. "github.com/LindsayBradford/crm/annealing/shared"
)

type AnnealerBuilder struct {
	annealer    Annealer
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
	if err := annealer.SetLogHandler(logHandler); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithStartingTemperature(temperature float64) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.SetTemperature(temperature); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithCoolingFactor(coolingFactor float64) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.SetCoolingFactor(coolingFactor); err != nil {
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

func (this *AnnealerBuilder) WithDumbObjectiveManager(initialObjectiveValue float64) *AnnealerBuilder {
	annealer := this.annealer
	objectiveManager := new(DumbObjectiveManager)
	objectiveManager.SetObjectiveValue(initialObjectiveValue)
	if err := annealer.SetObjectiveManager(objectiveManager); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithMaxIterations(iterations uint) *AnnealerBuilder {
	annealer := this.annealer
	annealer.SetMaxIterations(iterations)
	return this
}

func (this *AnnealerBuilder) WithObservers(observers ...AnnealingObserver) *AnnealerBuilder {
	if observers == nil {
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