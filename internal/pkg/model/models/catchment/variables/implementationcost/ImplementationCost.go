// Copyright (c) 2019 Australian Rivers Institute.

package implementationcost

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/math"
)

const VariableName = "ImplementationCost"
const notImplementedCost float64 = 0

var _ variable.UndoableDecisionVariable = new(ImplementationCost)

type ImplementationCost struct {
	variable.PerPlanningUnitDecisionVariable

	actionObserved action.ManagementAction

	command *variable.ChangePerPlanningUnitDecisionVariableCommand
}

func (ic *ImplementationCost) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *ImplementationCost {
	ic.PerPlanningUnitDecisionVariable.Initialise()
	ic.SetName(VariableName)
	ic.SetValue(ic.deriveInitialImplementationCost())
	ic.SetUnitOfMeasure(variable.Dollars)
	ic.SetPrecision(2)

	return ic
}

func (ic *ImplementationCost) WithObservers(observers ...variable.Observer) *ImplementationCost {
	ic.Subscribe(observers...)
	return ic
}

func (ic *ImplementationCost) deriveInitialImplementationCost() float64 {
	return notImplementedCost
}

func (ic *ImplementationCost) ObserveAction(action action.ManagementAction) {
	ic.observeAction(action)
}

func (ic *ImplementationCost) ObserveActionInitialising(action action.ManagementAction) {
	ic.observeAction(action)
	ic.command.Do()
}

func (ic *ImplementationCost) observeAction(action action.ManagementAction) {
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

func (ic *ImplementationCost) handleActionForModelVariable(name action.ModelVariableName) {
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

func (ic *ImplementationCost) UndoableValue() float64 {
	return ic.command.Value()
}

func (ic *ImplementationCost) SetUndoableValue(value float64) {
	ic.command.WithChange(value)
}

func (ic *ImplementationCost) DifferenceInValues() float64 {
	return ic.command.Change()
}

func (ic *ImplementationCost) ApplyDoneValue() {
	ic.command.Do()
}

func (ic *ImplementationCost) ApplyUndoneValue() {
	ic.command.Undo()
}
