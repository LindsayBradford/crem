// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/pkg/errors"
)

const ImplementationCost2VariableName = "ImplementationCost2"
const notImplementedCost float64 = 0

type ImplementationCost2 struct {
	variable.BaseInductiveDecisionVariable
	actionObserved action.ManagementAction

	valuePerPlanningUnit map[planningunit.Id]float64
}

func (ic *ImplementationCost2) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *ImplementationCost2 {
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
	ic.valuePerPlanningUnit = make(map[planningunit.Id]float64, 0)
	return notImplementedCost
}

func (ic *ImplementationCost2) ObserveAction(action action.ManagementAction) {
	ic.actionObserved = action
	switch ic.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		ic.handleRiverBankRestorationAction()
	case actions.GullyRestorationType:
		ic.handleGullyRestorationAction()
	case actions.HillSlopeRestorationType:
		ic.handleHillSlopeRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (ic *ImplementationCost2) ObserveActionInitialising(action action.ManagementAction) {
	ic.actionObserved = action
	switch ic.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		ic.handleInitialisingRiverBankRestorationAction()
	case actions.GullyRestorationType:
		ic.handleInitialisingGullyRestorationAction()
	case actions.HillSlopeRestorationType:
		ic.handleInitialisingHillSlopeRestorationAction()
	default:
		panic(errors.New("Unhandled observation of initialising management action type [" + string(action.Type()) + "]"))
	}
	ic.NotifyObservers()
}

func (ic *ImplementationCost2) handleRiverBankRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()
		ic.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsCost + toBeCost)
		ic.acceptPlanningUnitChange(asIsCost, toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.RiverBankRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setTempVariable(notImplementedCost, implementationCost)
	case false:
		setTempVariable(implementationCost, notImplementedCost)
	}
}

func (ic *ImplementationCost2) handleInitialisingRiverBankRestorationAction() {
	setVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()
		ic.BaseInductiveDecisionVariable.SetValue(currentValue - asIsCost + toBeCost)
		ic.acceptPlanningUnitChange(asIsCost, toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.RiverBankRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setVariable(notImplementedCost, implementationCost)
	case false:
		setVariable(implementationCost, notImplementedCost)
	}
}

func (ic *ImplementationCost2) handleGullyRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()
		ic.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsCost + toBeCost)

		ic.acceptPlanningUnitChange(asIsCost, toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.GullyRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setTempVariable(notImplementedCost, implementationCost)
	case false:
		setTempVariable(implementationCost, notImplementedCost)
	}
}

func (ic *ImplementationCost2) handleInitialisingGullyRestorationAction() {
	setVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()
		ic.BaseInductiveDecisionVariable.SetValue(currentValue - asIsCost + toBeCost)
		ic.acceptPlanningUnitChange(asIsCost, toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.GullyRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setVariable(notImplementedCost, implementationCost)
	case false:
		setVariable(implementationCost, notImplementedCost)
	}
}

func (ic *ImplementationCost2) handleHillSlopeRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()
		ic.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsCost + toBeCost)
		ic.acceptPlanningUnitChange(asIsCost, toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.HillSlopeRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setTempVariable(notImplementedCost, implementationCost)
	case false:
		setTempVariable(implementationCost, notImplementedCost)
	}
}

func (ic *ImplementationCost2) handleInitialisingHillSlopeRestorationAction() {
	setVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()
		ic.BaseInductiveDecisionVariable.SetValue(currentValue - asIsCost + toBeCost)
		ic.acceptPlanningUnitChange(asIsCost, toBeCost)
	}

	implementationCost := ic.actionObserved.ModelVariableValue(actions.HillSlopeRestorationCost)

	switch ic.actionObserved.IsActive() {
	case true:
		setVariable(notImplementedCost, implementationCost)
	case false:
		setVariable(implementationCost, notImplementedCost)
	}
}

func (ic *ImplementationCost2) acceptPlanningUnitChange(asIsCost float64, toBeCost float64) {
	planningUnit := ic.actionObserved.PlanningUnit()
	ic.valuePerPlanningUnit[planningUnit] = ic.valuePerPlanningUnit[planningUnit] - asIsCost + toBeCost
}

func (ic *ImplementationCost2) ValuesPerPlanningUnit() variableNew.PlanningUnitValueMap {
	return ic.valuePerPlanningUnit
}

func (ic *ImplementationCost2) RejectInductiveValue() {
	ic.rejectPlanningUnitChange()
	ic.BaseInductiveDecisionVariable.RejectInductiveValue()
}

func (ic *ImplementationCost2) rejectPlanningUnitChange() {
	change := ic.BaseInductiveDecisionVariable.DifferenceInValues()
	planningUnit := ic.actionObserved.PlanningUnit()

	ic.valuePerPlanningUnit[planningUnit] = ic.valuePerPlanningUnit[planningUnit] - change
}
