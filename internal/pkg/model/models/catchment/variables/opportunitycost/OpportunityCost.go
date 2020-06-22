// Copyright (c) 2019 Australian Rivers Institute.

package opportunitycost

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/math"
)

const VariableName = "OpportunityCost"
const notImplementedCost float64 = 0

var _ variable.UndoableDecisionVariable = new(OpportunityCost)

type OpportunityCost struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	actionObserved action.ManagementAction

	command variable.ChangeCommand
}

func (ic *OpportunityCost) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *OpportunityCost {
	ic.PerPlanningUnitDecisionVariable.Initialise()

	ic.command = new(variable.NullChangeCommand)

	ic.SetName(VariableName)
	ic.SetValue(ic.deriveInitialCost())
	ic.SetUnitOfMeasure(variable.Dollars)
	ic.SetPrecision(2)

	return ic
}

func (ic *OpportunityCost) WithObservers(observers ...variable.Observer) *OpportunityCost {
	ic.Subscribe(observers...)
	return ic
}

func (ic *OpportunityCost) deriveInitialCost() float64 {
	return notImplementedCost
}

func (ic *OpportunityCost) ObserveAction(action action.ManagementAction) {
	ic.observeAction(action)
}

func (ic *OpportunityCost) ObserveActionInitialising(action action.ManagementAction) {
	ic.observeAction(action)
	ic.command.Do()
}

func (ic *OpportunityCost) observeAction(action action.ManagementAction) {
	ic.actionObserved = action
	switch ic.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		ic.handleActionForModelVariable(actions.RiverBankRestorationCost)
	case actions.GullyRestorationType:
		ic.handleActionForModelVariable(actions.GullyRestorationCost)
	case actions.HillSlopeRestorationType:
		ic.handleActionForModelVariable(actions.HillSlopeRestorationCost)
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (ic *OpportunityCost) handleActionForModelVariable(name action.ModelVariableName) {
	actionCost := ic.actionObserved.ModelVariableValue(name)

	var newValue float64
	switch ic.actionObserved.IsActive() {
	case true:
		newValue = actionCost
	case false:
		newValue = -1 * actionCost
	}

	newValue = math.RoundFloat(newValue, int(ic.Precision()))

	ic.command = new(variable.ChangePerPlanningUnitDecisionVariableCommand).
		ForVariable(ic).
		InPlanningUnit(ic.actionObserved.PlanningUnit()).
		WithChange(newValue)
}

func (ic *OpportunityCost) UndoableValue() float64 {
	return ic.Value() + ic.command.Value()
}

func (ic *OpportunityCost) SetUndoableValue(value float64) {
	ic.command.SetChange(value)
}

func (ic *OpportunityCost) DifferenceInValues() float64 {
	return ic.command.Change()
}

func (ic *OpportunityCost) ApplyDoneValue() {
	ic.command.Do()
}

func (ic *OpportunityCost) ApplyUndoneValue() {
	ic.command.Undo()
}
