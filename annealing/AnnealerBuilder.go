// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	. "github.com/LindsayBradford/crm/annealing/solution"
	. "github.com/LindsayBradford/crm/annealing/shared"
	crmerrors "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealerBuilder struct {
	annealer    Annealer
	buildErrors *crmerrors.CompositeError
}

func (this *AnnealerBuilder) OSThreadLockedAnnealer() *AnnealerBuilder {
	this.annealer = &OSThreadLockedAnnealer{}
	this.annealer.Initialise()
	this.buildErrors = crmerrors.NewComposite("Failed to build valid OS thread-locked annealer")
	return this
}

func (this *AnnealerBuilder) ElapsedTimeTrackingAnnealer() *AnnealerBuilder {
	this.annealer = &ElapsedTimeTrackingAnnealer{}
	this.annealer.Initialise()
	this.buildErrors = crmerrors.NewComposite("Failed to build valid elapsed-timed tracking annealer")
	return this
}

func (this *AnnealerBuilder) SimpleAnnealer() *AnnealerBuilder {
	this.annealer = &SimpleAnnealer{}
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

func (this *AnnealerBuilder) WithStateTourer(tourer SolutionTourer) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.SetSolutionTourer(tourer); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithEventNotifier(delegate AnnealingEventNotifier) *AnnealerBuilder {
	annealer := this.annealer
	if err := annealer.SetEventNotifier(delegate); err != nil {
		this.buildErrors.Add(err)
	}
	return this
}

func (this *AnnealerBuilder) WithDumbObjectiveManager(initialObjectiveValue float64) *AnnealerBuilder {
	annealer := this.annealer
	objectiveManager := new(DumbSolutionTourer)
	objectiveManager.SetObjectiveValue(initialObjectiveValue)
	annealer.SetSolutionTourer(objectiveManager)
	return this
}

func (this *AnnealerBuilder) WithMaxIterations(iterations uint) *AnnealerBuilder {
	annealer := this.annealer
	annealer.SetMaxIterations(iterations)
	return this
}

func (this *AnnealerBuilder) WithObservers(observers ...AnnealingObserver) *AnnealerBuilder {
	annealer := this.annealer

	for _, currObserver := range observers {
		if err := annealer.AddObserver(currObserver); err != nil {
			this.buildErrors.Add(err)
		}
	}

	return this
}

func (this *AnnealerBuilder) Build() (Annealer, *crmerrors.CompositeError) {
	if this.buildErrors.Size() == 0 {
		return this.annealer, nil
	} else {
		return this.annealer, this.buildErrors
	}
}
