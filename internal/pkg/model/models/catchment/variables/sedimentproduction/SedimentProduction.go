// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const VariableName = "SedimentProduction"

var _ variable.DecisionVariable = new(SedimentProduction)

const planningUnitIndex = 0

type SedimentProduction struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	actionObserved action.ManagementAction

	command variable.ChangeCommand

	bankSedimentContribution      actions.BankSedimentContribution
	gullySedimentContribution     actions.GullySedimentContribution
	hillSlopeSedimentContribution actions.HillSlopeSedimentContribution

	cachedPlanningUnitSediment float64

	riparianVegetationProportionPerPlanningUnit  map[planningunit.Id]float64
	hillSlopeVegetationProportionPerPlanningUnit map[planningunit.Id]float64
}

func (sl *SedimentProduction) Initialise(planningUnitTable tables.CsvTable, gulliesTable tables.CsvTable, parameters parameters.Parameters) *SedimentProduction {
	sl.PerPlanningUnitDecisionVariable.Initialise()

	sl.SetName(VariableName)
	sl.SetUnitOfMeasure(variable.TonnesPerYear)
	sl.SetPrecision(3)

	sl.command = new(variable.NullChangeCommand)

	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.gullySedimentContribution.Initialise(gulliesTable, parameters)
	sl.hillSlopeSedimentContribution.Initialise(planningUnitTable, parameters)

	sl.deriveInitialPerPlanningUnitSedimentLoad(planningUnitTable)

	return sl
}

func (sl *SedimentProduction) WithObservers(observers ...variable.Observer) *SedimentProduction {
	sl.Subscribe(observers...)
	return sl
}

func (sl *SedimentProduction) deriveInitialPerPlanningUnitSedimentLoad(planningUnitTable tables.CsvTable) {
	_, rowCount := planningUnitTable.ColumnAndRowSize()

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

		sedimentProduced :=
			bankSedimentContribution +
				gullySedimentContribution +
				hillSlopeSedimentContribution

		roundedSedimentProduced := math.RoundFloat(sedimentProduced, int(sl.Precision()))

		sl.SetPlanningUnitValue(planningUnit, roundedSedimentProduced)
	}
}

func Float64ToPlanningUnitId(value float64) planningunit.Id {
	return planningunit.Id(value)
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

func (sl *SedimentProduction) ObserveAction(action action.ManagementAction) {
	sl.observeAction(action)
}

func (sl *SedimentProduction) ObserveActionInitialising(action action.ManagementAction) {
	sl.observeAction(action)
	sl.command.Do()
}

func (sl *SedimentProduction) observeAction(action action.ManagementAction) {
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

func (sl *SedimentProduction) handleRiverBankRestorationAction() {
	var asIsSediment, toBeSediment, vegetationBuffer float64
	switch sl.actionObserved.IsActive() {
	case true:
		vegetationBuffer = sl.actionObserved.ModelVariableValue(actions.ActionedBufferVegetation)
		toBeSediment = sl.planningUnitSediment(actions.ActionedBufferVegetation)
		asIsSediment = sl.planningUnitSediment(actions.OriginalBufferVegetation)
	case false:
		vegetationBuffer = sl.actionObserved.ModelVariableValue(actions.OriginalBufferVegetation)
		toBeSediment = sl.planningUnitSediment(actions.OriginalBufferVegetation)
		asIsSediment = sl.planningUnitSediment(actions.ActionedBufferVegetation)
	}

	sl.command = new(RiverBankRestorationCommand).
		ForVariable(sl).
		InPlanningUnit(sl.actionObserved.PlanningUnit()).
		WithVegetationBuffer(vegetationBuffer).
		WithChange(toBeSediment - asIsSediment)
}

func (sl *SedimentProduction) planningUnitSediment(riparianVegetationBufferName action.ModelVariableName) float64 {
	planningUnit := sl.actionObserved.PlanningUnit()

	riparianSediment := sl.riparianSediment(riparianVegetationBufferName, planningUnit)
	hillSlopeSediment := sl.hillSlopeSediment(planningUnit)

	return riparianSediment + hillSlopeSediment
}

func (sl *SedimentProduction) riparianSediment(vegetationBufferName action.ModelVariableName, planningUnit planningunit.Id) float64 {
	riparianVegetation := sl.actionObserved.ModelVariableValue(vegetationBufferName)
	riparianSediment := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, riparianVegetation)
	return riparianSediment
}

func (sl *SedimentProduction) hillSlopeSediment(planningUnit planningunit.Id) float64 {
	hillSlopeVegetation := sl.hillSlopeVegetationProportionPerPlanningUnit[planningUnit]
	filteredHillSlopeSediment := sl.filteredHillSlopeSediment(planningUnit, hillSlopeVegetation)
	return filteredHillSlopeSediment
}

func (sl *SedimentProduction) handleGullyRestorationAction() {
	var toBeSediment, asIsSediment float64

	switch sl.actionObserved.IsActive() {
	case true:
		toBeSediment = sl.actionObserved.ModelVariableValue(actions.ActionedGullySediment)
		asIsSediment = sl.actionObserved.ModelVariableValue(actions.OriginalGullySediment)
	case false:
		toBeSediment = sl.actionObserved.ModelVariableValue(actions.OriginalGullySediment)
		asIsSediment = sl.actionObserved.ModelVariableValue(actions.ActionedGullySediment)
	}

	sl.command = new(variable.ChangePerPlanningUnitDecisionVariableCommand).
		ForVariable(sl).
		InPlanningUnit(sl.actionObserved.PlanningUnit()).
		WithChange(toBeSediment - asIsSediment)
}

func (sl *SedimentProduction) handleHillSlopeRestorationAction() {
	var asIsSediment, toBeSediment, vegetationBuffer float64
	switch sl.actionObserved.IsActive() {
	case true:
		vegetationBuffer = sl.actionObserved.ModelVariableValue(actions.ActionedHillSlopeVegetation)
		toBeSediment = sl.hillSlopeSedimentForVariable(actions.ActionedHillSlopeVegetation)
		asIsSediment = sl.hillSlopeSedimentForVariable(actions.OriginalHillSlopeVegetation)
	case false:
		vegetationBuffer = sl.actionObserved.ModelVariableValue(actions.OriginalHillSlopeVegetation)
		toBeSediment = sl.hillSlopeSedimentForVariable(actions.OriginalHillSlopeVegetation)
		asIsSediment = sl.hillSlopeSedimentForVariable(actions.ActionedHillSlopeVegetation)
	}

	sl.command = new(HillSlopeRevegetationCommand).
		ForVariable(sl).
		InPlanningUnit(sl.actionObserved.PlanningUnit()).
		WithVegetationBuffer(vegetationBuffer).
		WithChange(toBeSediment - asIsSediment)
}

func (sl *SedimentProduction) hillSlopeSedimentForVariable(vegetationBufferName action.ModelVariableName) float64 {
	hillSlopeVegetation := sl.actionObserved.ModelVariableValue(vegetationBufferName)
	filteredHillSlopeSediment := sl.filteredHillSlopeSediment(sl.actionObserved.PlanningUnit(), hillSlopeVegetation)
	return filteredHillSlopeSediment
}

func (sl *SedimentProduction) filteredHillSlopeSediment(planningUnit planningunit.Id, hillSlopeVegetation float64) float64 {
	hillSlopeSediment := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, hillSlopeVegetation)
	filter := riparianBufferFilter(sl.riparianVegetationProportionPerPlanningUnit[planningUnit])
	filteredHillSlopeSediment := hillSlopeSediment * filter

	return filteredHillSlopeSediment
}

func (sl *SedimentProduction) UndoableValue() float64 {
	return sl.command.Value()
}

func (sl *SedimentProduction) SetUndoableValue(value float64) {
	sl.command.SetChange(value)
}

func (sl *SedimentProduction) DifferenceInValues() float64 {
	return sl.command.Change()
}

func (sl *SedimentProduction) ApplyDoneValue() {
	sl.command.Do()
}

func (sl *SedimentProduction) ApplyUndoneValue() {
	sl.command.Undo()
}
