// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/pkg/name"
)

const SedimentVsCostName = "SedimentVsCost"

var _ variable.InductiveDecisionVariable = new(SedimentVsCost)

type SedimentVsCost struct {
	name.ContainedName

	sedimentLoad       *SedimentLoad
	implementationCost *ImplementationCost

	variable.ContainedDecisionVariableObservers
}

func (sc *SedimentVsCost) Initialise() *SedimentVsCost {
	sc.SetName(SedimentVsCostName)
	return sc
}

func (sc *SedimentVsCost) WithSedimentLoad(sedimentLoad *SedimentLoad) *SedimentVsCost {
	sc.sedimentLoad = sedimentLoad
	return sc
}

func (sc *SedimentVsCost) WithImplementationCost(implementationCost *ImplementationCost) *SedimentVsCost {
	sc.implementationCost = implementationCost
	return sc
}

func (sc *SedimentVsCost) Value() float64 {
	return sc.sedimentLoad.Value() + sc.implementationCost.Value()
}

func (sc *SedimentVsCost) SetValue(value float64) {
	// Deliberately does nothing
}

func (sc *SedimentVsCost) InductiveValue() float64 {
	return sc.sedimentLoad.InductiveValue() + sc.implementationCost.InductiveValue()
}

func (sc *SedimentVsCost) SetInductiveValue(value float64) {
	// Deliberately does nothing
}

func (sc *SedimentVsCost) AcceptInductiveValue() {
	if sc.Value() == sc.InductiveValue() {
		return
	}

	sc.sedimentLoad.AcceptInductiveValue()
	sc.implementationCost.AcceptInductiveValue()

	sc.NotifyObservers()
}

func (sc *SedimentVsCost) RejectInductiveValue() {
	if sc.Value() == sc.InductiveValue() {
		return
	}

	sc.sedimentLoad.RejectInductiveValue()
	sc.implementationCost.RejectInductiveValue()

	sc.NotifyObservers()
}

func (sc *SedimentVsCost) DifferenceInValues() float64 {
	return sc.InductiveValue() - sc.Value()
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variable.
func (sc *SedimentVsCost) NotifyObservers() {
	for _, observer := range sc.Observers() {
		observer.ObserveDecisionVariable(sc)
	}
}
