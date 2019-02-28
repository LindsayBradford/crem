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

func (sl *ImplementationCost) Initialise(planningUnitTable *tables.CsvTable, parameters parameters.Parameters) *ImplementationCost {
	sl.SetName(ImplementationCostVariableName)
	sl.SetValue(sl.deriveInitialImplementationCost())
	return sl
}

func (sl *ImplementationCost) deriveInitialImplementationCost() float64 {
	return notImplementedCost
}

func (sl *ImplementationCost) ObserveAction(action action.ManagementAction) {
	sl.actionObserved = action
	switch sl.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		sl.handleRiverBankRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (sl *ImplementationCost) handleRiverBankRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := sl.VolatileDecisionVariable.Value()
		sl.VolatileDecisionVariable.SetTemporaryValue(currentValue - asIsCost + toBeCost)
	}

	implementationCost := sl.actionObserved.ModelVariableValue(actions.RiverBankRestorationCost)

	switch sl.actionObserved.IsActive() {
	case true:
		setTempVariable(notImplementedCost, implementationCost)
	case false:
		setTempVariable(implementationCost, notImplementedCost)
	}
}
