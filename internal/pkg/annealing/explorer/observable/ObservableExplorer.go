// Copyright (c) 2019 Australian Rivers Institute.

package observable

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

type ContainedObservable struct {
	acceptanceProbability float64
	changeIsDesirable     bool
	changeAccepted        bool
	objectiveValue        float64
	objectiveValueChange  float64
	temperature           float64
}

func (ke *ContainedObservable) Temperature() float64 {
	return ke.temperature
}

func (ke *ContainedObservable) SetTemperature(temperature float64) {
	ke.temperature = temperature
}

func (ke *ContainedObservable) ChangeIsDesirable() bool {
	return ke.changeIsDesirable
}

func (ke *ContainedObservable) SetChangeIsDesirable(changeIsDesirable bool) {
	ke.changeIsDesirable = changeIsDesirable
}

func (ke *ContainedObservable) ChangeInObjectiveValue() float64 {
	return ke.objectiveValueChange
}

func (ke *ContainedObservable) SetChangeInObjectiveValue(change float64) {
	ke.objectiveValueChange = change
}

func (ke *ContainedObservable) ObjectiveValue() float64 {
	return ke.objectiveValue
}

func (ke *ContainedObservable) SetObjectiveValue(objectiveValue float64) {
	ke.objectiveValue = objectiveValue
}

func (ke *ContainedObservable) ChangeAccepted() bool {
	return ke.changeAccepted
}

func (ke *ContainedObservable) SetChangeAccepted(changeAccepted bool) {
	ke.changeAccepted = changeAccepted
}

func (ke *ContainedObservable) AcceptanceProbability() float64 {
	return ke.acceptanceProbability
}

func (ke *ContainedObservable) SetAcceptanceProbability(probability float64) {
	ke.acceptanceProbability = math.Min(explorer.Guaranteed, probability)
}
