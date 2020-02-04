// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const notImplementedValue float64 = 0

var _ variable.UndoableDecisionVariable = new(DumbObjective)

type DumbObjective struct {
	variable.PerPlanningUnitDecisionVariable

	command variable.ChangeCommand

	actionObserved action.ManagementAction
}

func (o *DumbObjective) Initialise() *DumbObjective {
	o.PerPlanningUnitDecisionVariable.Initialise()

	o.command = new(variable.NullChangeCommand)

	o.SetUnitOfMeasure(variable.NotApplicable)
	o.SetPrecision(2)
	return o
}

func (o *DumbObjective) WithName(variableName string) *DumbObjective {
	o.SetName(variableName)
	return o
}

func (o *DumbObjective) WithStartingValue(value float64) *DumbObjective {
	o.SetPlanningUnitValue(0, value)
	return o
}

func (o *DumbObjective) WithObservers(observers ...variable.Observer) *DumbObjective {
	o.Subscribe(observers...)
	return o
}

func (o *DumbObjective) deriveInitialValue() float64 {
	o.SetValue(notImplementedValue)
	return notImplementedValue
}

func (o *DumbObjective) ObserveAction(action action.ManagementAction) {
	o.observeAction(action)
}

func (o *DumbObjective) ObserveActionInitialising(action action.ManagementAction) {
	o.observeAction(action)
	o.command.Do()
}

func (o *DumbObjective) observeAction(action action.ManagementAction) {
	o.actionObserved = action
	switch o.actionObserved.Type() {
	case actions.DumbActionType:
		o.handleDumbAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
func (o *DumbObjective) NotifyObservers() {
	for _, observer := range o.Observers() {
		observer.ObserveDecisionVariable(o)
	}
}

func (o *DumbObjective) handleDumbAction() {
	variableName := action.ModelVariableName(o.Name())
	actionCost := o.actionObserved.ModelVariableValue(variableName)

	var change float64
	switch o.actionObserved.IsActive() {
	case true:
		change = actionCost
	case false:
		change = actionCost * -1
	}

	change = math.RoundFloat(change, int(o.Precision()))

	o.command = new(variable.ChangePerPlanningUnitDecisionVariableCommand).
		ForVariable(o).
		InPlanningUnit(o.actionObserved.PlanningUnit()).
		WithChange(change)
}

func (o *DumbObjective) UndoableValue() float64 {
	return o.command.Value()
}

func (o *DumbObjective) SetUndoableValue(value float64) {
	o.command.SetChange(value)
}

func (o *DumbObjective) DifferenceInValues() float64 {
	return o.command.Change()
}

func (o *DumbObjective) ApplyDoneValue() {
	o.command.Do()
}

func (o *DumbObjective) ApplyUndoneValue() {
	o.command.Undo()
}
