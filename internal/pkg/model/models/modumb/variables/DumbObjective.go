// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/pkg/errors"
)

const notImplementedValue float64 = 0

type DumbObjective struct {
	variable.BaseInductiveDecisionVariable
	actionObserved       action.ManagementAction
	valuePerPlanningUnit map[planningunit.Id]float64
}

func (do *DumbObjective) Initialise() *DumbObjective {
	do.SetUnitOfMeasure(variable.NotApplicable)
	do.SetPrecision(2)
	do.valuePerPlanningUnit = make(map[planningunit.Id]float64, 0)
	return do
}

func (do *DumbObjective) WithName(variableName string) *DumbObjective {
	do.SetName(variableName)
	return do
}

func (do *DumbObjective) WithStartingValue(value float64) *DumbObjective {
	do.SetValue(value)
	return do
}

func (do *DumbObjective) WithObservers(observers ...variable.Observer) *DumbObjective {
	do.Subscribe(observers...)
	return do
}

func (do *DumbObjective) deriveInitialValue() float64 {
	do.valuePerPlanningUnit = make(map[planningunit.Id]float64, 0)
	return notImplementedValue
}

func (do *DumbObjective) ObserveAction(action action.ManagementAction) {
	do.actionObserved = action
	switch do.actionObserved.Type() {
	case actions.DumbActionType:
		do.handleDumbAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (do *DumbObjective) ObserveActionInitialising(action action.ManagementAction) {
	do.actionObserved = action
	switch do.actionObserved.Type() {
	case actions.DumbActionType:
		do.handleInitialisingDumbAction()
	default:
		panic(errors.New("Unhandled observation of initialising management action type [" + string(action.Type()) + "]"))
	}
	do.NotifyObservers()
}

func (do *DumbObjective) handleDumbAction() {
	setTempVariable := func(asIsValue float64, toBeValue float64) {
		if asIsValue == toBeValue {
			return
		}
		currentValue := do.BaseInductiveDecisionVariable.Value()
		do.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsValue + toBeValue)
		do.acceptPlanningUnitChange(asIsValue, toBeValue)
	}

	variableName := action.ModelVariableName(do.Name())
	variableValue := do.actionObserved.ModelVariableValue(variableName)

	switch do.actionObserved.IsActive() {
	case true:
		setTempVariable(notImplementedValue, variableValue)
	case false:
		setTempVariable(variableValue, notImplementedValue)
	}
}

func (do *DumbObjective) handleInitialisingDumbAction() {
	setVariable := func(asIsValue float64, toBeValue float64) {
		if asIsValue == toBeValue {
			return
		}
		currentValue := do.BaseInductiveDecisionVariable.Value()
		do.BaseInductiveDecisionVariable.SetValue(currentValue - asIsValue + toBeValue)
		do.acceptPlanningUnitChange(asIsValue, toBeValue)
	}

	variableName := action.ModelVariableName(do.Name())
	variableValue := do.actionObserved.ModelVariableValue(variableName)

	switch do.actionObserved.IsActive() {
	case true:
		setVariable(notImplementedValue, variableValue)
	case false:
		setVariable(variableValue, notImplementedValue)
	}
}

func (do *DumbObjective) acceptPlanningUnitChange(asIsSedimentContribution float64, toBeSedimentContribution float64) {
	planningUnit := do.actionObserved.PlanningUnit()
	do.valuePerPlanningUnit[planningUnit] = do.valuePerPlanningUnit[planningUnit] - asIsSedimentContribution + toBeSedimentContribution
}

func (do *DumbObjective) ValuesPerPlanningUnit() map[planningunit.Id]float64 {
	return do.valuePerPlanningUnit
}

func (do *DumbObjective) RejectInductiveValue() {
	do.rejectPlanningUnitChange()
	do.BaseInductiveDecisionVariable.RejectInductiveValue()
}

func (do *DumbObjective) rejectPlanningUnitChange() {
	change := do.BaseInductiveDecisionVariable.DifferenceInValues()
	planningUnit := do.actionObserved.PlanningUnit()

	do.valuePerPlanningUnit[planningUnit] = do.valuePerPlanningUnit[planningUnit] - change
}
