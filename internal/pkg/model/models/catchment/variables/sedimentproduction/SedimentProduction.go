// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	catchmentParameters "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const (
	VariableName = "SedimentProduction"

	RiverbankVegetationProportion = "RiverbankVegetationProportion"
	HillSlopeVegetationProportion = "HillSlopeVegetationProportion"
	RiverbankSedimentContribution = "RiverbankSedimentContribution"
	GullySedimentContribution     = "GullySedimentContribution"
	HillSlopeSedimentContribution = "HillSlopeSedimentContribution"
)

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

	numberOfPlanningUnits      uint
	cachedPlanningUnitSediment float64
	hillSlopeDeliveryRatio     float64

	planningUnitAttributes map[planningunit.Id]attributes.Attributes
}

func (sl *SedimentProduction) Initialise(planningUnitTable tables.CsvTable, gulliesTable tables.CsvTable, parameters catchmentParameters.Parameters) *SedimentProduction {
	sl.PerPlanningUnitDecisionVariable.Initialise()

	sl.SetName(VariableName)
	sl.SetUnitOfMeasure(variable.TonnesPerYear)
	sl.SetPrecision(3)

	sl.hillSlopeDeliveryRatio = parameters.GetFloat64(catchmentParameters.HillSlopeDeliveryRatio)

	sl.command = new(variable.NullChangeCommand)

	sl.deriveInitialState(planningUnitTable, gulliesTable, parameters)

	return sl
}

func (sl *SedimentProduction) deriveInitialState(planningUnitTable tables.CsvTable, gulliesTable tables.CsvTable, parameters catchmentParameters.Parameters) {
	sl.deriveNumberOfPlanningUnits(planningUnitTable)

	sl.initialisePlanningUnitAttributes()

	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.gullySedimentContribution.Initialise(gulliesTable, parameters)
	sl.hillSlopeSedimentContribution.Initialise(planningUnitTable, parameters)

	sl.deriveInitialSedimentProduction(planningUnitTable)
}

func (sl *SedimentProduction) initialisePlanningUnitAttributes() {
	sl.planningUnitAttributes = make(map[planningunit.Id]attributes.Attributes, sl.numberOfPlanningUnits)
	for index, _ := range sl.planningUnitAttributes {
		newAttributes := make(attributes.Attributes, 0)
		sl.planningUnitAttributes[index] = newAttributes
	}
}

func (sl *SedimentProduction) deriveNumberOfPlanningUnits(planningUnitTable tables.CsvTable) {
	_, rowCount := planningUnitTable.ColumnAndRowSize()
	sl.numberOfPlanningUnits = rowCount
}

func (sl *SedimentProduction) NumberOfPlanningUnits() uint {
	return sl.numberOfPlanningUnits
}

func (sl *SedimentProduction) WithObservers(observers ...variable.Observer) *SedimentProduction {
	sl.Subscribe(observers...)
	return sl
}

func (sl *SedimentProduction) deriveInitialSedimentProduction(planningUnitTable tables.CsvTable) {
	for row := uint(0); row < sl.numberOfPlanningUnits; row++ {
		planningUnitFloat64 := planningUnitTable.CellFloat64(planningUnitIndex, row)
		planningUnit := Float64ToPlanningUnitId(planningUnitFloat64)

		riverbankSedimentContribution := sl.bankSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit)
		gullySedimentContribution := sl.gullySedimentContribution.SedimentContribution(planningUnit)

		riverBankVegetationProportion := sl.originalRiverbankVegetationProportion(planningUnit)
		riparianFilter := riparianBufferFilter(riverBankVegetationProportion)
		hillSlopeSedimentContribution := sl.hillSlopeSedimentContribution.OriginalPlanningUnitSedimentContribution(planningUnit) * riparianFilter

		sl.planningUnitAttributes[planningUnit] = new(attributes.Attributes).
			Add(RiverbankVegetationProportion, riverBankVegetationProportion).
			Add(RiverbankSedimentContribution, riverbankSedimentContribution).
			Add(GullySedimentContribution, gullySedimentContribution).
			Add(HillSlopeSedimentContribution, hillSlopeSedimentContribution)

		sedimentProduced := riverbankSedimentContribution + gullySedimentContribution + hillSlopeSedimentContribution
		roundedSedimentProduced := math.RoundFloat(sedimentProduced, int(sl.Precision()))

		sl.SetPlanningUnitValue(planningUnit, roundedSedimentProduced)
	}
}

func (sl *SedimentProduction) originalRiverbankVegetationProportion(id planningunit.Id) float64 {
	return sl.bankSedimentContribution.OriginalPlanningUnitVegetationProportion(id)
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
	var toBeRiverBankSediment, asIsRiverBankSediment, toBeVegetation, asIsVegetation float64

	switch sl.actionObserved.IsActive() {
	case true:
		asIsVegetation = sl.actionObserved.ModelVariableValue(actions.OriginalBufferVegetation)
		asIsRiverBankSediment = sl.actionObserved.ModelVariableValue(actions.ActionedRiparianSedimentProduction)

		toBeVegetation = sl.actionObserved.ModelVariableValue(actions.ActionedBufferVegetation)
		toBeRiverBankSediment = sl.actionObserved.ModelVariableValue(actions.ActionedRiparianSedimentProduction)
	case false:
		asIsVegetation = sl.actionObserved.ModelVariableValue(actions.ActionedBufferVegetation)
		asIsRiverBankSediment = sl.actionObserved.ModelVariableValue(actions.ActionedRiparianSedimentProduction)

		toBeVegetation = sl.actionObserved.ModelVariableValue(actions.OriginalBufferVegetation)
		toBeRiverBankSediment = sl.actionObserved.ModelVariableValue(actions.ActionedRiparianSedimentProduction)
	}

	attributes := sl.planningUnitAttributes[sl.actionObserved.PlanningUnit()]

	asIsContext := sedimentContext{
		riparianVVegetationProportion: asIsVegetation,
		riparianContribution:          asIsRiverBankSediment,
		gullyContribution:             attributes.Value(GullySedimentContribution).(float64),
		hillSlopeContribution:         attributes.Value(HillSlopeSedimentContribution).(float64),
	}

	asIsSediment := sl.calculateSedimentProduction(asIsContext)

	toBeContext := sedimentContext{
		riparianVVegetationProportion: toBeVegetation,
		riparianContribution:          toBeRiverBankSediment,
		gullyContribution:             attributes.Value(GullySedimentContribution).(float64),
		hillSlopeContribution:         attributes.Value(HillSlopeSedimentContribution).(float64),
	}

	toBeSediment := sl.calculateSedimentProduction(toBeContext)

	sl.command = new(RiverBankRestorationCommand).
		ForVariable(sl).
		InPlanningUnit(sl.actionObserved.PlanningUnit()).
		WithVegetationProportion(toBeVegetation).
		WithRiverBankContribution(toBeRiverBankSediment).
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
	attribs := sl.planningUnitAttributes[planningUnit]
	hillSlopeSediment := attribs.Value(HillSlopeSedimentContribution).(float64)
	filteredHillSlopeSediment := sl.filteredHillSlopeSediment(planningUnit, hillSlopeSediment)
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

	sl.command = new(GullyRestorationCommand).
		ForVariable(sl).
		InPlanningUnit(sl.actionObserved.PlanningUnit()).
		WithChange(toBeSediment - asIsSediment)
}

func (sl *SedimentProduction) handleHillSlopeRestorationAction() {
	var toBeHillSlopeSediment, asIsHillSlopeSediment float64

	switch sl.actionObserved.IsActive() {
	case true:
		toBeHillSlopeSediment = sl.actionObserved.ModelVariableValue(actions.HillSlopeErosionActionedAttribute)
		asIsHillSlopeSediment = sl.actionObserved.ModelVariableValue(actions.HillSlopeErosionOriginalAttribute)
	case false:
		asIsHillSlopeSediment = sl.actionObserved.ModelVariableValue(actions.HillSlopeErosionActionedAttribute)
		toBeHillSlopeSediment = sl.actionObserved.ModelVariableValue(actions.HillSlopeErosionOriginalAttribute)
	}

	attributes := sl.planningUnitAttributes[sl.actionObserved.PlanningUnit()]

	asIsContext := sedimentContext{
		riparianVVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		riparianContribution:          attributes.Value(RiverbankSedimentContribution).(float64),
		gullyContribution:             attributes.Value(GullySedimentContribution).(float64),
		hillSlopeContribution:         asIsHillSlopeSediment,
	}

	asIsSediment := sl.calculateSedimentProduction(asIsContext)

	toBeContext := sedimentContext{
		riparianVVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		riparianContribution:          attributes.Value(RiverbankSedimentContribution).(float64),
		gullyContribution:             attributes.Value(GullySedimentContribution).(float64),
		hillSlopeContribution:         toBeHillSlopeSediment,
	}

	toBeSediment := sl.calculateSedimentProduction(toBeContext)

	sl.command = new(HillSlopeRevegetationCommand).
		ForVariable(sl).
		InPlanningUnit(sl.actionObserved.PlanningUnit()).
		WithSedimentContribution(toBeHillSlopeSediment).
		WithChange(toBeSediment - asIsSediment)
}

func (sl *SedimentProduction) filteredHillSlopeSediment(planningUnit planningunit.Id, hillSlopeVegetation float64) float64 {
	hillSlopeSediment := sl.hillSlopeSedimentContribution.PlanningUnitSedimentContribution(planningUnit, hillSlopeVegetation)

	attribs := sl.planningUnitAttributes[planningUnit]
	filter := attribs.Value(RiverbankVegetationProportion).(float64)
	filteredHillSlopeSediment := hillSlopeSediment * filter

	return filteredHillSlopeSediment
}

func (sl *SedimentProduction) UndoableValue() float64 {
	return sl.Value() + sl.command.Value()
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

func (sl *SedimentProduction) PlanningUnitAttributes() map[planningunit.Id]attributes.Attributes {
	return sl.planningUnitAttributes
}

func (sl *SedimentProduction) Command() variable.ChangeCommand {
	return sl.command
}

type sedimentContext struct {
	riparianContribution          float64
	riparianVVegetationProportion float64

	hillSlopeContribution float64
	gullyContribution     float64
}

func (sl *SedimentProduction) calculateSedimentProduction(context sedimentContext) float64 {
	deliveryAdjustedHillSlopeContribution := context.hillSlopeContribution * sl.hillSlopeDeliveryRatio

	riparianFilter := riparianBufferFilter(context.riparianVVegetationProportion)
	filteredHillSlopeContribution := deliveryAdjustedHillSlopeContribution * riparianFilter

	nitrogenProduced := context.riparianContribution + context.gullyContribution + filteredHillSlopeContribution

	roundedSedimentProduced := math.RoundFloat(nitrogenProduced, int(sl.Precision()))
	return roundedSedimentProduced
}
