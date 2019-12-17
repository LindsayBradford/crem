// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableOld"
)

const SedimentVsCostName = "SedimentVsCost"

var _ variable.UndoableDecisionVariable = new(SedimentVsCost)

type SedimentVsCost struct {
	variableOld.CompositeInductiveDecisionVariable
}

func (sc *SedimentVsCost) Initialise() *SedimentVsCost {
	sc.SetName(SedimentVsCostName)
	sc.SetUnitOfMeasure(variable.NotApplicable)
	sc.SetPrecision(6)
	sc.CompositeInductiveDecisionVariable.Initialise()
	return sc
}

func (sc *SedimentVsCost) WithObservers(observers ...variable.Observer) *SedimentVsCost {
	sc.Subscribe(observers...)
	return sc
}

func (sc *SedimentVsCost) ObserveActionInitialising(action action.ManagementAction) {
	// deliberately does nothing
}

func (sc *SedimentVsCost) ObserveAction(action action.ManagementAction) {
	// deliberately does nothing
}
