// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/actions"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/pkg/errors"
)

const ImplementationCostVariableName = "ImplementationCost"
const notImplementedCost float64 = 0

type ImplementationCost struct {
	variable.VolatileDecisionVariable
	actionObserved action.ManagementAction
}

func (ic *ImplementationCost) Initialise(planningUnitTable *tables.CsvTable, parameters parameters.Parameters) *ImplementationCost {
	ic.SetName(ImplementationCostVariableName)
	ic.SetValue(ic.deriveInitialImplementationCost())
	return ic
}

func (ic *ImplementationCost) deriveInitialImplementationCost() float64 {
	return notImplementedCost
}

func (ic *ImplementationCost) ObserveAction(action action.ManagementAction) {
	ic.actionObserved = action
	switch ic.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		ic.handleRiverBankRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (ic *ImplementationCost) ObserveInitialisationAction(action action.ManagementAction) {
	ic.actionObserved = action
	switch ic.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		ic.handleInitialisingRiverBankRestorationAction()
	default:
		panic(errors.New("Unhandled observation of initialising management action type [" + string(action.Type()) + "]"))
	}
	ic.NotifyObservers()
}

func (ic *ImplementationCost) handleRiverBankRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.VolatileDecisionVariable.Value()
		ic.VolatileDecisionVariable.SetTemporaryValue(currentValue - asIsCost + toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.RiverBankRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setTempVariable(notImplementedCost, implementationCost)
	case false:
		setTempVariable(implementationCost, notImplementedCost)
	}
}

func (ic *ImplementationCost) handleInitialisingRiverBankRestorationAction() {
	setVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.VolatileDecisionVariable.Value()
		ic.VolatileDecisionVariable.SetValue(currentValue - asIsCost + toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.RiverBankRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setVariable(notImplementedCost, implementationCost)
	case false:
		setVariable(implementationCost, notImplementedCost)
	}
}
