// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	math2 "math"

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

	sedimentPerPlanningUnit    map[planningunit.Id]float64
	cachedPlanningUnitSediment float64

	riparianVegetationProportionPerPlanningUnit  map[planningunit.Id]float64
	hillSlopeVegetationProportionPerPlanningUnit map[planningunit.Id]float64
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
	sl.sedimentPerPlanningUnit = make(map[planningunit.Id]float64, rowCount)
	sl.riparianVegetationProportionPerPlanningUnit = make(map[planningunit.Id]float64, rowCount)
	sl.hillSlopeVegetationProportionPerPlanningUnit = make(map[planningunit.Id]float64, rowCount)

	for row := uint(0); row < rowCount; row++ {
		planningUnitFloat64 := planningUnitTable.CellFloat64(planningUnitIndex, row)
		planningUnit := Float64ToPlanningUnitId(planningUnitFloat64)

		sl.riparianVegetationProportionPerPlanningUnit[planningUnit] =
			sl.bankSedimentContribution.OriginalPlanningUnitVegetationProportion(planningUnit)

		sl.hillSlopeVegetationProportionPerPlanningUnit[planningUnit] =
			sl.hillSlopeSedimentContribution.OriginalPlanningUnitVegetationProportion(planningUnit)

		bankSedimentContribution := sl.bankSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit)
		gullySedimentContribution := sl.gullySedimentContribution.SedimentContribution(planningUnit)

		riparianFilter := riparianBufferFilter(sl.riparianVegetationProportionPerPlanningUnit[planningUnit])
		hillSlopeSedimentContribution := sl.hillSlopeSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit) * riparianFilter

		sl.sedimentPerPlanningUnit[planningUnit] =
			bankSedimentContribution +
				gullySedimentContribution +
				hillSlopeSedimentContribution

		sl.sedimentPerPlanningUnit[planningUnit] = math.RoundFloat(sl.sedimentPerPlanningUnit[planningUnit], int(sl.Precision()))
	}
}

func (sl *SedimentProduction) deriveInitialSedimentLoad() float64 {
	initialSedimentLoad := float64(0)
	for _, sedimentAtPlanningUnit := range sl.sedimentPerPlanningUnit {
		initialSedimentLoad += sedimentAtPlanningUnit
	}

	initialSedimentLoad = math.RoundFloat(initialSedimentLoad, int(sl.Precision()))

	return initialSedimentLoad
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
		asIsRiparianVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeRiparianVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		asIsRiparianFilter := riparianBufferFilter(sl.actionObserved.ModelVariableValue(asIsName))

		sl.riparianVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] =
			sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()

		asIsRiparianSediment := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsRiparianVegetation)
		toBeRiparianSediment := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeRiparianVegetation)

		hillSlopeVegetation := sl.hillSlopeVegetationProportionPerPlanningUnit[planningUnit]
		toBeRiparianFilter := riparianBufferFilter(sl.riparianVegetationProportionPerPlanningUnit[planningUnit])

		rawHillSlopeSediment := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, hillSlopeVegetation)

		asIsHillSlopeSediment := rawHillSlopeSediment * asIsRiparianFilter
		toBeHillSlopeSediment := rawHillSlopeSediment * toBeRiparianFilter

		currentSediment := sl.BaseInductiveDecisionVariable.Value()
		asIsPlanningUnitSediment := asIsRiparianSediment + asIsHillSlopeSediment
		toBePlanningUnitSediment := toBeRiparianSediment + toBeHillSlopeSediment
		newSediment := currentSediment - asIsPlanningUnitSediment + toBePlanningUnitSediment

		sl.BaseInductiveDecisionVariable.SetInductiveValue(newSediment)

		sl.acceptPlanningUnitChange(asIsPlanningUnitSediment, toBePlanningUnitSediment)
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
		asIsRiparianVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeRiparianVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()

		asIsRiparianFilter := riparianBufferFilter(asIsRiparianVegetation)

		sl.riparianVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] = sl.actionObserved.ModelVariableValue(toBeName)

		asIsRiparianSediment := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsRiparianVegetation)
		toBeRiparianSediment := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeRiparianVegetation)

		hillSlopeVegetation := sl.hillSlopeVegetationProportionPerPlanningUnit[planningUnit]
		toBeRiparianFilter := riparianBufferFilter(toBeRiparianVegetation)

		rawHillSlopeSediment := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, hillSlopeVegetation)

		asIsHillSlopeSediment := rawHillSlopeSediment * asIsRiparianFilter
		toBeHillSlopeSediment := rawHillSlopeSediment * toBeRiparianFilter

		currentSediment := sl.BaseInductiveDecisionVariable.Value()

		asIsSediment := asIsRiparianSediment + asIsHillSlopeSediment
		toBeSediment := toBeRiparianSediment + toBeHillSlopeSediment

		newSediment := currentSediment - asIsSediment + toBeSediment

		sl.BaseInductiveDecisionVariable.SetValue(newSediment)

		sl.acceptPlanningUnitChange(asIsSediment, toBeSediment)
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
	setTempVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()
		riparianFilter := riparianBufferFilter(sl.riparianVegetationProportionPerPlanningUnit[planningUnit])

		rawAsIsSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsVegetation)
		asIsSedimentContribution := rawAsIsSedimentContribution * riparianFilter
		rawToBeSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeVegetation)
		toBeSedimentContribution := rawToBeSedimentContribution * riparianFilter

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)

		sl.hillSlopeVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] = toBeVegetation

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
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		planningUnit := sl.actionObserved.PlanningUnit()
		riparianFilter := riparianBufferFilter(sl.riparianVegetationProportionPerPlanningUnit[planningUnit])

		rawAsIsSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, asIsVegetation)
		asIsSedimentContribution := rawAsIsSedimentContribution * riparianFilter
		rawToBeSedimentContribution := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, toBeVegetation)
		toBeSedimentContribution := rawToBeSedimentContribution * riparianFilter

		previousSediment := sl.BaseInductiveDecisionVariable.Value()
		newSediment := previousSediment - asIsSedimentContribution + toBeSedimentContribution
		sl.BaseInductiveDecisionVariable.SetValue(newSediment)

		sl.hillSlopeVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] = toBeVegetation

		sl.acceptPlanningUnitChange(asIsSedimentContribution, toBeSedimentContribution)
	}

	assert.That(sl.actionObserved.IsActive()).WithFailureMessage("initialising action should always be active").Holds()
	setVariable(actions.OriginalHillSlopeVegetation, actions.ActionedHillSlopeVegetation)
}

func riparianBufferFilter(proportionOfRiparianBufferVegetation float64) float64 {
	if proportionOfRiparianBufferVegetation < 0.25 {
		return 1
	}
	if proportionOfRiparianBufferVegetation > 0.75 {
		return 0.25
	}
	return 1 - proportionOfRiparianBufferVegetation
}

func (sl *SedimentProduction) acceptPlanningUnitChange(asIsSediment float64, toBeSediment float64) {
	planningUnit := sl.actionObserved.PlanningUnit()

	newPlanningUnitSediment := sl.sedimentPerPlanningUnit[planningUnit] + toBeSediment - asIsSediment
	newPlanningUnitSediment = math.RoundFloat(newPlanningUnitSediment, int(sl.Precision()))

	// max against 0 here to stop tiny rounding errors resulting in negative sediment for a planning unit.
	sl.cachedPlanningUnitSediment = sl.sedimentPerPlanningUnit[planningUnit]
	sl.sedimentPerPlanningUnit[planningUnit] = math2.Max(0, newPlanningUnitSediment)
}

func (sl *SedimentProduction) ValuesPerPlanningUnit() map[planningunit.Id]float64 {
	return sl.sedimentPerPlanningUnit
}

func (sl *SedimentProduction) RejectInductiveValue() {
	sl.rejectPlanningUnitChange()
	sl.BaseInductiveDecisionVariable.RejectInductiveValue()
}

func (sl *SedimentProduction) rejectPlanningUnitChange() {
	planningUnit := sl.actionObserved.PlanningUnit()
	sl.sedimentPerPlanningUnit[planningUnit] = sl.cachedPlanningUnitSediment

	if sl.actionObserved.Type() == actions.RiverBankRestorationType {
		switch sl.actionObserved.IsActive() {
		case true:
			{
				sl.riparianVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] =
					sl.actionObserved.ModelVariableValue(actions.OriginalBufferVegetation)
			}
		default:
			{
				sl.riparianVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] =
					sl.actionObserved.ModelVariableValue(actions.ActionedBufferVegetation)
			}
		}
	}

	if sl.actionObserved.Type() == actions.HillSlopeRestorationType {
		switch sl.actionObserved.IsActive() {
		case true:
			{
				sl.hillSlopeVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] =
					sl.actionObserved.ModelVariableValue(actions.OriginalHillSlopeVegetation)
			}
		default:
			{
				sl.hillSlopeVegetationProportionPerPlanningUnit[sl.actionObserved.PlanningUnit()] =
					sl.actionObserved.ModelVariableValue(actions.ActionedHillSlopeVegetation)
			}
		}
	}
}
