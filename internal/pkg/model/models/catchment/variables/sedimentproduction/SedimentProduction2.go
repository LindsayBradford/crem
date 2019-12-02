// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const SedimentProduction2VariableName = "SedimentProduction2"

var _ variableNew.DecisionVariable = new(SedimentProduction2)

const planningUnitIndex = 0

type SedimentProduction2 struct {
	variableNew.PerPlanningUnitDecisionVariable

	actionObserved action.ManagementAction

	command variableNew.ChangeCommand

	bankSedimentContribution      actions.BankSedimentContribution
	gullySedimentContribution     actions.GullySedimentContribution
	hillSlopeSedimentContribution actions.HillSlopeSedimentContribution

	cachedPlanningUnitSediment float64

	riparianVegetationProportionPerPlanningUnit  map[planningunit.Id]float64
	hillSlopeVegetationProportionPerPlanningUnit map[planningunit.Id]float64
}

func (sl *SedimentProduction2) Initialise(planningUnitTable tables.CsvTable, gulliesTable tables.CsvTable, parameters parameters.Parameters) *SedimentProduction2 {
	sl.PerPlanningUnitDecisionVariable.Initialise()

	sl.SetName(SedimentProduction2VariableName)
	sl.SetUnitOfMeasure(variableNew.TonnesPerYear)
	sl.SetPrecision(3)

	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.gullySedimentContribution.Initialise(gulliesTable, parameters)
	sl.hillSlopeSedimentContribution.Initialise(planningUnitTable, parameters)

	sl.deriveInitialPerPlanningUnitSedimentLoad(planningUnitTable)

	return sl
}

func (sl *SedimentProduction2) WithObservers(observers ...variableNew.Observer) *SedimentProduction2 {
	sl.Subscribe(observers...)
	return sl
}

func (sl *SedimentProduction2) deriveInitialPerPlanningUnitSedimentLoad(planningUnitTable tables.CsvTable) {
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

func (sl *SedimentProduction2) ObserveAction(action action.ManagementAction) {
	sl.observeAction(action)
}

func (sl *SedimentProduction2) ObserveActionInitialising(action action.ManagementAction) {
	sl.observeAction(action)
	sl.command.Do()
	sl.NotifyObservers() // TODO: Needed?
}

func (sl *SedimentProduction2) observeAction(action action.ManagementAction) {
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

func (sl *SedimentProduction2) handleRiverBankRestorationAction() {
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

func (sl *SedimentProduction2) planningUnitSediment(riparianVegetationBufferName action.ModelVariableName) float64 {
	planningUnit := sl.actionObserved.PlanningUnit()

	riparianSediment := sl.riparianSediment(riparianVegetationBufferName, planningUnit)
	hillSlopeSediment := sl.hillSlopeSediment(planningUnit)

	return riparianSediment + hillSlopeSediment
}

func (sl *SedimentProduction2) riparianSediment(vegetationBufferName action.ModelVariableName, planningUnit planningunit.Id) float64 {
	riparianVegetation := sl.actionObserved.ModelVariableValue(vegetationBufferName)
	riparianSediment := sl.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnit, riparianVegetation)
	return riparianSediment
}

func (sl *SedimentProduction2) hillSlopeSediment(planningUnit planningunit.Id) float64 {
	hillSlopeVegetation := sl.hillSlopeVegetationProportionPerPlanningUnit[planningUnit]
	filteredHillSlopeSediment := sl.filteredHillSlopeSediment(planningUnit, hillSlopeVegetation)
	return filteredHillSlopeSediment
}

func (sl *SedimentProduction2) handleGullyRestorationAction() {
	var toBeSediment, asIsSediment float64

	switch sl.actionObserved.IsActive() {
	case true:
		toBeSediment = sl.actionObserved.ModelVariableValue(actions.ActionedGullySediment)
		asIsSediment = sl.actionObserved.ModelVariableValue(actions.OriginalGullySediment)
	case false:
		toBeSediment = sl.actionObserved.ModelVariableValue(actions.OriginalGullySediment)
		asIsSediment = sl.actionObserved.ModelVariableValue(actions.ActionedGullySediment)
	}

	sl.command = new(variableNew.ChangePerPlanningUnitDecisionVariableCommand).
		ForVariable(sl).
		InPlanningUnit(sl.actionObserved.PlanningUnit()).
		WithChange(toBeSediment - asIsSediment)
}

func (sl *SedimentProduction2) handleHillSlopeRestorationAction() {
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

func (sl *SedimentProduction2) hillSlopeSedimentForVariable(vegetationBufferName action.ModelVariableName) float64 {
	hillSlopeVegetation := sl.actionObserved.ModelVariableValue(vegetationBufferName)
	filteredHillSlopeSediment := sl.filteredHillSlopeSediment(sl.actionObserved.PlanningUnit(), hillSlopeVegetation)
	return filteredHillSlopeSediment
}

func (sl *SedimentProduction2) filteredHillSlopeSediment(planningUnit planningunit.Id, hillSlopeVegetation float64) float64 {
	hillSlopeSediment := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, hillSlopeVegetation)
	filter := riparianBufferFilter(sl.riparianVegetationProportionPerPlanningUnit[planningUnit])
	filteredHillSlopeSediment := hillSlopeSediment * filter

	return filteredHillSlopeSediment
}

func (sl *SedimentProduction2) InductiveValue() float64 {
	return sl.command.Value()
}

func (sl *SedimentProduction2) SetInductiveValue(value float64) {
	sl.command.SetChange(value)
}

func (sl *SedimentProduction2) DifferenceInValues() float64 {
	return sl.command.Change()
}

func (sl *SedimentProduction2) AcceptInductiveValue() {
	sl.command.Do()
	sl.NotifyObservers() // TODO: Needed?
}

func (sl *SedimentProduction2) RejectInductiveValue() {
	sl.command.Undo()
	sl.NotifyObservers() // TODO: Needed?
}
