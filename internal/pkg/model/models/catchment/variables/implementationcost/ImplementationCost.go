// Copyright (c) 2019 Australian Rivers Institute.

package implementationcost

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const ImplementationCostVariableName = "ImplementationCost"

type ImplementationCost struct {
	variable.BaseInductiveDecisionVariable
	actionObserved action.ManagementAction

	planningUnitValues variableNew.PlanningUnitValueMap
}

func (ic *ImplementationCost) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *ImplementationCost {
	ic.SetName(ImplementationCostVariableName)
	ic.SetValue(ic.deriveInitialImplementationCost())
	ic.SetUnitOfMeasure(variableNew.Dollars)
	ic.SetPrecision(2)
	return ic
}

func (ic *ImplementationCost) WithObservers(observers ...variableNew.Observer) *ImplementationCost {
	ic.Subscribe(observers...)
	return ic
}

func (ic *ImplementationCost) deriveInitialImplementationCost() float64 {
	ic.planningUnitValues = make(map[planningunit.Id]float64, 0)
	return notImplementedCost
}

func (ic *ImplementationCost) ObserveAction(action action.ManagementAction) {
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

func (ic *ImplementationCost) ObserveActionInitialising(action action.ManagementAction) {
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

func (ic *ImplementationCost) handleRiverBankRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()

		newValue := currentValue - asIsCost + toBeCost
		newValue = math.RoundFloat(newValue, int(ic.Precision()))

		ic.BaseInductiveDecisionVariable.SetInductiveValue(newValue)
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

func (ic *ImplementationCost) handleInitialisingRiverBankRestorationAction() {
	setVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()

		newValue := currentValue - asIsCost + toBeCost
		newValue = math.RoundFloat(newValue, int(ic.Precision()))

		ic.BaseInductiveDecisionVariable.SetValue(newValue)
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

func (ic *ImplementationCost) handleGullyRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()

		newValue := currentValue - asIsCost + toBeCost
		newValue = math.RoundFloat(newValue, int(ic.Precision()))

		ic.BaseInductiveDecisionVariable.SetInductiveValue(newValue)

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

func (ic *ImplementationCost) handleInitialisingGullyRestorationAction() {
	setVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()

		newValue := currentValue - asIsCost + toBeCost
		newValue = math.RoundFloat(newValue, int(ic.Precision()))

		ic.BaseInductiveDecisionVariable.SetValue(newValue)
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

func (ic *ImplementationCost) handleHillSlopeRestorationAction() {
	setTempVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()

		newValue := currentValue - asIsCost + toBeCost
		newValue = math.RoundFloat(newValue, int(ic.Precision()))

		ic.BaseInductiveDecisionVariable.SetInductiveValue(newValue)
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

func (ic *ImplementationCost) handleInitialisingHillSlopeRestorationAction() {
	setVariable := func(asIsCost float64, toBeCost float64) {
		currentValue := ic.BaseInductiveDecisionVariable.Value()

		newValue := currentValue - asIsCost + toBeCost
		newValue = math.RoundFloat(newValue, int(ic.Precision()))
		ic.BaseInductiveDecisionVariable.SetValue(newValue)

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

func (ic *ImplementationCost) acceptPlanningUnitChange(asIsCost float64, toBeCost float64) {
	planningUnit := ic.actionObserved.PlanningUnit()
	ic.planningUnitValues[planningUnit] = ic.planningUnitValues[planningUnit] - asIsCost + toBeCost
	ic.planningUnitValues[planningUnit] = math.RoundFloat(ic.planningUnitValues[planningUnit], int(ic.Precision()))
}

func (ic *ImplementationCost) ValuesPerPlanningUnit() variableNew.PlanningUnitValueMap {
	return ic.planningUnitValues
}

func (ic *ImplementationCost) RejectInductiveValue() {
	ic.rejectPlanningUnitChange()
	ic.BaseInductiveDecisionVariable.RejectInductiveValue()
}

func (ic *ImplementationCost) rejectPlanningUnitChange() {
	change := ic.BaseInductiveDecisionVariable.DifferenceInValues()
	planningUnit := ic.actionObserved.PlanningUnit()

	ic.planningUnitValues[planningUnit] = ic.planningUnitValues[planningUnit] - change
}

func (ic *ImplementationCost) SetPlanningUnitValue(planningUnit planningunit.Id, newValue float64) {
	ic.planningUnitValues[planningUnit] = newValue
}

func (ic *ImplementationCost) PlanningUnitValue(planningUnit planningunit.Id) float64 {
	return ic.planningUnitValues[planningUnit]
}
