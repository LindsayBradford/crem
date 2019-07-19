// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/assert/release"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const SedimentProductionVariableName = "SedimentProduction"

var _ variable.DecisionVariable = new(SedimentProduction)

const planningUnitIndex = 0

func Float64ToPlanningUnitId(value float64) planningunit.Id {
	return planningunit.Id(value)
}

type SedimentProduction struct {
	variable.BaseInductiveDecisionVariable

	bankSedimentContribution      actions.BankSedimentContribution
	gullySedimentContribution     actions.GullySedimentContribution
	hillSlopeSedimentContribution actions.HillSlopeSedimentContribution

	actionObserved action.ManagementAction

	valuePerPlanningUnit map[planningunit.Id]float64
}

func (sl *SedimentProduction) Initialise(planningUnitTable tables.CsvTable, gulliesTable tables.CsvTable, parameters parameters.Parameters) *SedimentProduction {
	sl.SetName(SedimentProductionVariableName)
	sl.SetUnitOfMeasure(variable.TonnesPerYear)
	sl.SetPrecision(3)

	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.gullySedimentContribution.Initialise(gulliesTable, parameters)
	sl.hillSlopeSedimentContribution.Initialise(planningUnitTable, parameters)

	sl.deriveInitialPerPlanningUnitSedimentLoad(planningUnitTable)
	sl.SetValue(sl.deriveInitialSedimentLoad())

	return sl
}

func (sl *SedimentProduction) WithObservers(observers ...variable.Observer) *SedimentProduction {
	sl.Subscribe(observers...)
	return sl
}

func (sl *SedimentProduction) deriveInitialPerPlanningUnitSedimentLoad(planningUnitTable tables.CsvTable) {
	_, rowCount := planningUnitTable.ColumnAndRowSize()
	sl.valuePerPlanningUnit = make(map[planningunit.Id]float64, rowCount)

	for row := uint(0); row < rowCount; row++ {
		planningUnitFloat64 := planningUnitTable.CellFloat64(planningUnitIndex, row)
		planningUnit := Float64ToPlanningUnitId(planningUnitFloat64)

		sl.valuePerPlanningUnit[planningUnit] =
			sl.bankSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit) +
				sl.gullySedimentContribution.SedimentContribution(planningUnit) +
				sl.hillSlopeSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit)
	}
}

func (sl *SedimentProduction) deriveInitialSedimentLoad() float64 {
	return sl.bankSedimentContribution.OriginalSedimentContribution() +
		sl.gullySedimentContribution.OriginalSedimentContribution() +
		sl.hillSlopeSedimentContribution.OriginalSedimentContribution()
}

func (sl *SedimentProduction) ObserveAction(action action.ManagementAction) {
	sl.actionObserved = action
	switch sl.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		sl.handleRiverBankRestorationAction()
	case actions.GullyRestorationType:
		sl.handleGullyRestorationAction()
	case actions.HillSlopeRestorationType:
		sl.handleHillSlopeRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (sl *SedimentProduction) ObserveActionInitialising(action action.ManagementAction) {
	sl.actionObserved = action
	switch sl.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		sl.handleInitialisingRiverBankRestorationAction()
	case actions.GullyRestorationType:
		sl.handleInitialisingGullyRestorationAction()
	case actions.HillSlopeRestorationType:
		sl.handleInitialisingHillSlopeRestorationAction()
	default:
		panic(errors.New("Unhandled observation of initialising management action type [" + string(action.Type()) + "]"))
	}
	sl.NotifyObservers()
}

func (sl *SedimentProduction) handleRiverBankRestorationAction() {
	setTempVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()

		asIsSedimentContribution := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsVegetation)
		toBeSedimentContribution := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeVegetation)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)

		sl.acceptPlanningUnitChange(asIsSedimentContribution, toBeSedimentContribution)
	}

	switch sl.actionObserved.IsActive() {
	case true:
		setTempVariable(actions.OriginalBufferVegetation, actions.ActionedBufferVegetation)
	case false:
		setTempVariable(actions.ActionedBufferVegetation, actions.OriginalBufferVegetation)
	}
}

func (sl *SedimentProduction) handleInitialisingRiverBankRestorationAction() {
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()

		asIsSedimentContribution := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsVegetation)
		toBeSedimentContribution := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeVegetation)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)

		sl.acceptPlanningUnitChange(asIsSedimentContribution, toBeSedimentContribution)
	}

	assert.That(sl.actionObserved.IsActive()).WithFailureMessage("initialising action should always be active").Holds()
	setVariable(actions.OriginalBufferVegetation, actions.ActionedBufferVegetation)
}

func (sl *SedimentProduction) handleGullyRestorationAction() {
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsSedimentContribution := sl.actionObserved.ModelVariableValue(asIsName)
		toBeSedimentContribution := sl.actionObserved.ModelVariableValue(toBeName)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)

		sl.acceptPlanningUnitChange(asIsSedimentContribution, toBeSedimentContribution)
	}

	switch sl.actionObserved.IsActive() {
	case true:
		setVariable(actions.OriginalGullySediment, actions.ActionedGullySediment)
	case false:
		setVariable(actions.ActionedGullySediment, actions.OriginalGullySediment)
	}
}

func (sl *SedimentProduction) handleInitialisingGullyRestorationAction() {
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsSedimentContribution := sl.actionObserved.ModelVariableValue(asIsName)
		toBeSedimentContribution := sl.actionObserved.ModelVariableValue(toBeName)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)

		sl.acceptPlanningUnitChange(asIsSedimentContribution, toBeSedimentContribution)
	}

	assert.That(sl.actionObserved.IsActive()).Holds()
	setVariable(actions.OriginalGullySediment, actions.ActionedGullySediment)
}

func (sl *SedimentProduction) handleHillSlopeRestorationAction() {
	// TODO:  Reduce if matching RiverBankingRestoration is active. CRP final report, table 3-4, sect 3.6.1, pg 31
	setTempVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()

		asIsSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsVegetation)
		toBeSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeVegetation)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)

		sl.acceptPlanningUnitChange(asIsSedimentContribution, toBeSedimentContribution)
	}

	switch sl.actionObserved.IsActive() {
	case true:
		setTempVariable(actions.OriginalHillSlopeVegetation, actions.ActionedHillSlopeVegetation)
	case false:
		setTempVariable(actions.ActionedHillSlopeVegetation, actions.OriginalHillSlopeVegetation)
	}
}

func (sl *SedimentProduction) handleInitialisingHillSlopeRestorationAction() {
	// TODO:  Need to reduce if matching RiverBankingRestoration is active. CRP final report, table 3-4, sect 3.6.1
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()

		asIsSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsVegetation)
		toBeSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeVegetation)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)

		sl.acceptPlanningUnitChange(asIsSedimentContribution, toBeSedimentContribution)
	}

	assert.That(sl.actionObserved.IsActive()).WithFailureMessage("initialising action should always be active").Holds()
	setVariable(actions.OriginalHillSlopeVegetation, actions.ActionedHillSlopeVegetation)
}

func (sl *SedimentProduction) acceptPlanningUnitChange(asIsSedimentContribution float64, toBeSedimentContribution float64) {
	planningUnit := sl.actionObserved.PlanningUnit()
	change := sl.valuePerPlanningUnit[planningUnit] - asIsSedimentContribution + toBeSedimentContribution
	sl.valuePerPlanningUnit[planningUnit] = math.RoundFloat(change, int(sl.Precision()))
}

func (sl *SedimentProduction) ValuesPerPlanningUnit() map[planningunit.Id]float64 {
	return sl.valuePerPlanningUnit
}

func (sl *SedimentProduction) RejectInductiveValue() {
	sl.rejectPlanningUnitChange()
	sl.BaseInductiveDecisionVariable.RejectInductiveValue()
}

func (sl *SedimentProduction) rejectPlanningUnitChange() {
	recordedChange := math.RoundFloat(sl.BaseInductiveDecisionVariable.DifferenceInValues(), int(sl.Precision()))
	planningUnit := sl.actionObserved.PlanningUnit()

	rejectChange := sl.valuePerPlanningUnit[planningUnit] - recordedChange
	sl.valuePerPlanningUnit[planningUnit] = math.RoundFloat(rejectChange, int(sl.Precision()))
}
