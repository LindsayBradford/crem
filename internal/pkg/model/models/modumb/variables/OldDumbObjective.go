// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableOld"
	"github.com/pkg/errors"
)

const notImplementedValue float64 = 0

type OldDumbObjective struct {
	variableOld.BaseInductiveDecisionVariable
	actionObserved       action.ManagementAction
	valuePerPlanningUnit map[planningunit.Id]float64
}

func (o *OldDumbObjective) Initialise() *OldDumbObjective {
	o.SetUnitOfMeasure(variable.NotApplicable)
	o.SetPrecision(2)
	o.valuePerPlanningUnit = make(map[planningunit.Id]float64, 0)
	return o
}

func (o *OldDumbObjective) WithName(variableName string) *OldDumbObjective {
	o.SetName(variableName)
	return o
}

func (o *OldDumbObjective) WithStartingValue(value float64) *OldDumbObjective {
	o.SetValue(value)
	return o
}

func (o *OldDumbObjective) WithObservers(observers ...variable.Observer) *OldDumbObjective {
	o.Subscribe(observers...)
	return o
}

func (o *OldDumbObjective) deriveInitialValue() float64 {
	o.valuePerPlanningUnit = make(map[planningunit.Id]float64, 0)
	return notImplementedValue
}

func (o *OldDumbObjective) ObserveAction(action action.ManagementAction) {
	o.actionObserved = action
	switch o.actionObserved.Type() {
	case actions.DumbActionType:
		o.handleDumbAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (o *OldDumbObjective) ObserveActionInitialising(action action.ManagementAction) {
	o.actionObserved = action
	switch o.actionObserved.Type() {
	case actions.DumbActionType:
		o.handleInitialisingDumbAction()
	default:
		panic(errors.New("Unhandled observation of initialising management action type [" + string(action.Type()) + "]"))
	}
	o.NotifyObservers()
}

func (o *OldDumbObjective) handleDumbAction() {
	setTempVariable := func(asIsValue float64, toBeValue float64) {
		if asIsValue == toBeValue {
			return
		}
		currentValue := o.BaseInductiveDecisionVariable.Value()
		o.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsValue + toBeValue)
		o.acceptPlanningUnitChange(asIsValue, toBeValue)
	}

	variableName := action.ModelVariableName(o.Name())
	variableValue := o.actionObserved.ModelVariableValue(variableName)

	switch o.actionObserved.IsActive() {
	case true:
		setTempVariable(notImplementedValue, variableValue)
	case false:
		setTempVariable(variableValue, notImplementedValue)
	}
}

func (o *OldDumbObjective) handleInitialisingDumbAction() {
	setVariable := func(asIsValue float64, toBeValue float64) {
		if asIsValue == toBeValue {
			return
		}
		currentValue := o.BaseInductiveDecisionVariable.Value()
		o.BaseInductiveDecisionVariable.SetValue(currentValue - asIsValue + toBeValue)
		o.acceptPlanningUnitChange(asIsValue, toBeValue)
	}

	variableName := action.ModelVariableName(o.Name())
	variableValue := o.actionObserved.ModelVariableValue(variableName)

	switch o.actionObserved.IsActive() {
	case true:
		setVariable(notImplementedValue, variableValue)
	case false:
		setVariable(variableValue, notImplementedValue)
	}
}

func (o *OldDumbObjective) acceptPlanningUnitChange(asIsSedimentContribution float64, toBeSedimentContribution float64) {
	planningUnit := o.actionObserved.PlanningUnit()
	o.valuePerPlanningUnit[planningUnit] = o.valuePerPlanningUnit[planningUnit] - asIsSedimentContribution + toBeSedimentContribution
}

func (o *OldDumbObjective) ValuesPerPlanningUnit() map[planningunit.Id]float64 {
	return o.valuePerPlanningUnit
}

func (o *OldDumbObjective) RejectInductiveValue() {
	o.rejectPlanningUnitChange()
	o.BaseInductiveDecisionVariable.RejectInductiveValue()
}

func (o *OldDumbObjective) rejectPlanningUnitChange() {
	change := o.BaseInductiveDecisionVariable.DifferenceInValues()
	planningUnit := o.actionObserved.PlanningUnit()

	o.valuePerPlanningUnit[planningUnit] = o.valuePerPlanningUnit[planningUnit] - change
}
