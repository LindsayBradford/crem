// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction2

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const VariableName = "SedimentProduction2"
const RiverbankVegetationProportion = "RiverbankVegetationProportion"
const HillSlopeVegetationProportion = "HillSlopeVegetationProportion"
const RiverbankSedimentContribution = "RiverbankSedimentContribution"
const GullySedimentContribution = "GullySedimentContribution"
const HillSlopeSedimentContribution = "HillSlopeSedimentContribution"
const SedimentProduced = "SedimentProduced"

var _ variable.DecisionVariable = new(SedimentProduction2)

const planningUnitIndex = 0

type SedimentProduction2 struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	actionObserved action.ManagementAction

	command variable.ChangeCommand

	bankSedimentContribution      actions.BankSedimentContribution
	gullySedimentContribution     actions.GullySedimentContribution
	hillSlopeSedimentContribution actions.HillSlopeSedimentContribution

	numberOfPlanningUnits      uint
	cachedPlanningUnitSediment float64

	// hillSlopeVegetationProportionPerPlanningUnit map[planningunit.Id]float64

	planningUnitAttributes map[planningunit.Id]attributes.Attributes
}

func (sl *SedimentProduction2) Initialise(planningUnitTable tables.CsvTable, gulliesTable tables.CsvTable, parameters parameters.Parameters) *SedimentProduction2 {
	sl.PerPlanningUnitDecisionVariable.Initialise()

	sl.SetName(VariableName)
	sl.SetUnitOfMeasure(variable.TonnesPerYear)
	sl.SetPrecision(3)

	sl.command = new(variable.NullChangeCommand)

	sl.deriveNumberOfPlanningUnits(planningUnitTable)

	sl.initialisePlanningUnitAttributes()

	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.gullySedimentContribution.Initialise(gulliesTable, parameters)
	sl.hillSlopeSedimentContribution.Initialise(planningUnitTable, parameters)

	sl.deriveInitialSedimentProduction(planningUnitTable)

	return sl
}

func (sl *SedimentProduction2) initialisePlanningUnitAttributes() {
	sl.planningUnitAttributes = make(map[planningunit.Id]attributes.Attributes, sl.numberOfPlanningUnits)
	for index, _ := range sl.planningUnitAttributes {
		newAttributes := make(attributes.Attributes, 0)
		sl.planningUnitAttributes[index] = newAttributes
	}
}

func (sl *SedimentProduction2) deriveNumberOfPlanningUnits(planningUnitTable tables.CsvTable) {
	_, rowCount := planningUnitTable.ColumnAndRowSize()
	sl.numberOfPlanningUnits = rowCount
}

func (sl *SedimentProduction2) NumberOfPlanningUnits() uint {
	return sl.numberOfPlanningUnits
}

func (sl *SedimentProduction2) WithObservers(observers ...variable.Observer) *SedimentProduction2 {
	sl.Subscribe(observers...)
	return sl
}

func (sl *SedimentProduction2) deriveInitialSedimentProduction(planningUnitTable tables.CsvTable) {
	for row := uint(0); row < sl.numberOfPlanningUnits; row++ {
		planningUnitFloat64 := planningUnitTable.CellFloat64(planningUnitIndex, row)
		planningUnit := Float64ToPlanningUnitId(planningUnitFloat64)

		riverBankVegetationProportion := sl.originalRiverbankVegetationProportion(planningUnit)
		hillSlopeVegetationProportion := sl.originalHillSlopeVegetationProportion(planningUnit)
		riverbankSedimentContribution := sl.bankSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit)
		gullySedimentContribution := sl.gullySedimentContribution.SedimentContribution(planningUnit)

		riparianFilter := riparianBufferFilter(riverBankVegetationProportion)
		hillSlopeSedimentContribution := sl.hillSlopeSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit) * riparianFilter

		sedimentProduced :=
			riverbankSedimentContribution +
				gullySedimentContribution +
				hillSlopeSedimentContribution

		roundedSedimentProduced := math.RoundFloat(sedimentProduced, int(sl.Precision()))

		sl.planningUnitAttributes[planningUnit] = sl.planningUnitAttributes[planningUnit].
			Add(RiverbankVegetationProportion, riverBankVegetationProportion).
			Add(HillSlopeVegetationProportion, hillSlopeVegetationProportion).
			Add(RiverbankSedimentContribution, riverbankSedimentContribution).
			Add(GullySedimentContribution, gullySedimentContribution).
			Add(HillSlopeSedimentContribution, hillSlopeSedimentContribution).
			Add(SedimentProduced, roundedSedimentProduced)

		sl.SetPlanningUnitValue(planningUnit, roundedSedimentProduced)
	}
}

func (sl *SedimentProduction2) originalRiverbankVegetationProportion(id planningunit.Id) float64 {
	return sl.bankSedimentContribution.OriginalPlanningUnitVegetationProportion(id)
}

func (sl *SedimentProduction2) originalHillSlopeVegetationProportion(id planningunit.Id) float64 {
	return sl.hillSlopeSedimentContribution.OriginalPlanningUnitVegetationProportion(id)
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

func (sl *SedimentProduction2) ObserveAction(action action.ManagementAction) {
	sl.observeAction(action)
}

func (sl *SedimentProduction2) ObserveActionInitialising(action action.ManagementAction) {
	sl.observeAction(action)
	sl.command.Do()
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
	attribs := sl.planningUnitAttributes[planningUnit]
	vegetationProportion := attribs.Value(HillSlopeVegetationProportion).(float64)
	filteredHillSlopeSediment := sl.filteredHillSlopeSediment(planningUnit, vegetationProportion)
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

	sl.command = new(variable.ChangePerPlanningUnitDecisionVariableCommand).
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

	attribs := sl.planningUnitAttributes[planningUnit]
	filter := attribs.Value(RiverbankVegetationProportion).(float64)
	filteredHillSlopeSediment := hillSlopeSediment * filter

	return filteredHillSlopeSediment
}

func (sl *SedimentProduction2) UndoableValue() float64 {
	return sl.Value() + sl.command.Value()
}

func (sl *SedimentProduction2) SetUndoableValue(value float64) {
	sl.command.SetChange(value)
}

func (sl *SedimentProduction2) DifferenceInValues() float64 {
	return sl.command.Change()
}

func (sl *SedimentProduction2) ApplyDoneValue() {
	sl.command.Do()
}

func (sl *SedimentProduction2) ApplyUndoneValue() {
	sl.command.Undo()
}
