// Copyright (c) 2019 Australian Rivers Institute.

package implementationcost

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/math"
)

const ImplementationCost2VariableName = "ImplementationCost2"
const notImplementedCost float64 = 0

var _ variable.InductiveDecisionVariable = new(ImplementationCost2)

type ImplementationCost2 struct {
	variableNew.PerPlanningUnitDecisionVariable

	actionObserved action.ManagementAction

	command *variableNew.ChangePerPlanningUnitDecisionVariableCommand
}

func (ic *ImplementationCost2) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *ImplementationCost2 {
	ic.PerPlanningUnitDecisionVariable.Initialise()

	ic.SetName(ImplementationCost2VariableName)
	ic.SetValue(ic.deriveInitialImplementationCost())
	ic.SetUnitOfMeasure(variableNew.Dollars)
	ic.SetPrecision(2)

	return ic
}

func (ic *ImplementationCost2) WithObservers(observers ...variableNew.Observer) *ImplementationCost2 {
	ic.Subscribe(observers...)
	return ic
}

func (ic *ImplementationCost2) deriveInitialImplementationCost() float64 {
	return notImplementedCost
}

func (ic *ImplementationCost2) ObserveAction(action action.ManagementAction) {
	ic.observeAction(action)
}

func (ic *ImplementationCost2) ObserveActionInitialising(action action.ManagementAction) {
	ic.observeAction(action)
	ic.command.Do()
	ic.NotifyObservers() // TODO: Needed?
}

func (ic *ImplementationCost2) observeAction(action action.ManagementAction) {
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

func (ic *ImplementationCost2) handleActionForModelVariable(name action.ModelVariableName) {
	actionCost := ic.actionObserved.ModelVariableValue(name)

	var newValue float64
	switch ic.actionObserved.IsActive() {
	case true:
		newValue = actionCost
	case false:
		newValue = -1 * actionCost
	}

	newValue = math.RoundFloat(newValue, int(ic.Precision()))

	ic.command = new(variableNew.ChangePerPlanningUnitDecisionVariableCommand).
		ForVariable(ic).
		InPlanningUnit(ic.actionObserved.PlanningUnit()).
		WithChange(newValue)
}

// TODO: This still feels janky!!

func (ic *ImplementationCost2) InductiveValue() float64 {
	return ic.command.Value()
}

func (ic *ImplementationCost2) SetInductiveValue(value float64) {
	ic.command.WithChange(value)
}

func (ic *ImplementationCost2) DifferenceInValues() float64 {
	return ic.command.Change()
}

func (ic *ImplementationCost2) AcceptInductiveValue() {
	ic.command.Do()
	ic.NotifyObservers() // TODO: Needed?
}

func (ic *ImplementationCost2) RejectInductiveValue() {
	ic.command.Undo()
	ic.NotifyObservers() // TODO: Needed?
}
