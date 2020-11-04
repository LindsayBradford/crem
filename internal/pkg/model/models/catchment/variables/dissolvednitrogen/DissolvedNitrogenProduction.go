// Copyright (c) 2019 Australian Rivers Institute.

package dissolvednitrogen

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	catchmentActions "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	catchmentParameters "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const (
	VariableName = "DissolvedNitrogen"

	planningUnitIndex           = 0
	proportionOfVegetationIndex = 8

	RiverbankVegetationProportion = "RiverbankVegetationProportion"
	RiparianFineSediment          = "RiparianFineSediment"

	RiparianNitrogenContribution  = "RiparianNitrogenContribution"
	GullyNitrogenContribution     = "GullyNitrogenContribution"
	HillSlopeNitrogenContribution = "HillSlopeNitrogenContribution"

	conversionFactor = 0.01
)

var _ variable.UndoableDecisionVariable = new(DissolvedNitrogenProduction)

type DissolvedNitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	catchmentActions.Container

	command variable.ChangeCommand

	actionObserved action.ManagementAction

	sedimentProductionVariable *sedimentproduction.SedimentProduction

	numberOfSubCatchments uint

	hillSlopeDeliveryRatio float64

	hillSlopeNitrogenContribution float64
	bankNitrogenContribution      float64
	gullyNitrogenContribution     float64

	subCatchmentAttributes map[planningunit.Id]attributes.Attributes
}

func (np *DissolvedNitrogenProduction) Initialise(subCatchmentsTable tables.CsvTable, actionsTable tables.CsvTable, parameters catchmentParameters.Parameters) *DissolvedNitrogenProduction {
	np.PerPlanningUnitDecisionVariable.Initialise()
	np.Container.WithActionsTable(actionsTable)

	np.SetName(VariableName)
	np.SetUnitOfMeasure(variable.TonnesPerYear)
	np.SetPrecision(3)

	np.hillSlopeDeliveryRatio = parameters.GetFloat64(catchmentParameters.HillSlopeDeliveryRatio)

	np.command = new(variable.NullChangeCommand)

	np.deriveInitialState(subCatchmentsTable, parameters)

	return np
}

func (np *DissolvedNitrogenProduction) WithName(variableName string) *DissolvedNitrogenProduction {
	np.SetName(variableName)
	return np
}

func (np *DissolvedNitrogenProduction) WithStartingValue(value float64) *DissolvedNitrogenProduction {
	np.SetPlanningUnitValue(0, value)
	return np
}

func (np *DissolvedNitrogenProduction) WithObservers(observers ...variable.Observer) *DissolvedNitrogenProduction {
	np.Subscribe(observers...)
	return np
}

func (np *DissolvedNitrogenProduction) WithSedimentProductionVariable(variable *sedimentproduction.SedimentProduction) *DissolvedNitrogenProduction {
	np.sedimentProductionVariable = variable
	return np
}

func (np *DissolvedNitrogenProduction) deriveInitialState(subCatchmentsTable tables.CsvTable, parameters catchmentParameters.Parameters) {
	np.deriveNumberOfSubCatchments(subCatchmentsTable)
	np.initialiseSubCatchmentAttributes()
	np.deriveInitialNitrogen(subCatchmentsTable)
}

func (np *DissolvedNitrogenProduction) deriveNumberOfSubCatchments(subCatchmentsTable tables.CsvTable) {
	_, rowCount := subCatchmentsTable.ColumnAndRowSize()
	np.numberOfSubCatchments = rowCount
}

func (np *DissolvedNitrogenProduction) initialiseSubCatchmentAttributes() {
	np.subCatchmentAttributes = make(map[planningunit.Id]attributes.Attributes, np.numberOfSubCatchments)
	for index, _ := range np.subCatchmentAttributes {
		newAttributes := make(attributes.Attributes, 0)
		np.subCatchmentAttributes[index] = newAttributes
	}
}

func (np *DissolvedNitrogenProduction) deriveInitialNitrogen(subCatchmentsTable tables.CsvTable) {
	np.buildDefaultSubCatchmentAttributes(subCatchmentsTable)
	np.replaceDefaultAttributeValuesWithActionValues()
	np.calculateInitialParticulateNitrogenPerSubCatchment()
}

func (np *DissolvedNitrogenProduction) buildDefaultSubCatchmentAttributes(subCatchmentsTable tables.CsvTable) {
	for row := uint(0); row < np.numberOfSubCatchments; row++ {
		subCatchmentFloat64 := subCatchmentsTable.CellFloat64(planningUnitIndex, row)
		subCatchment := Float64ToSubCatchmentId(subCatchmentFloat64)

		riverBankVegetationProportion := subCatchmentsTable.CellFloat64(proportionOfVegetationIndex, row)

		np.subCatchmentAttributes[subCatchment] =
			np.subCatchmentAttributes[subCatchment].
				Add(RiverbankVegetationProportion, riverBankVegetationProportion).
				Add(RiparianFineSediment, float64(0)).
				Add(RiparianNitrogenContribution, float64(0)).
				Add(HillSlopeNitrogenContribution, float64(0)).
				Add(GullyNitrogenContribution, float64(0))
	}
}

func (np *DissolvedNitrogenProduction) replaceDefaultAttributeValuesWithActionValues() {
	// Order below matters. Riparian nitrogen contribution depends on base attribute values being pre-calculated.
	np.calculateBaseAttributes()
	np.calculateRiparianNitrogenContributionAttribute()
}

func (np *DissolvedNitrogenProduction) calculateBaseAttributes() {
	for key, value := range np.Map() {
		components := np.DeriveMapKeyComponents(key)
		if components == nil {
			continue
		}

		np.calculateGullyAndHillSlopeContributions(components, value)
		np.calculateRiparianFineSediment(components, value)
	}
}

func (np *DissolvedNitrogenProduction) calculateGullyAndHillSlopeContributions(components *catchmentActions.KeyComponents, value float64) {
	if components.ElementType != catchmentActions.ParticulateNitrogenOriginalAttribute {
		return
	}

	switch components.SourceType {
	case catchmentActions.HillSlopeSource:
		deliveryAdjustedValue := value * np.hillSlopeDeliveryRatio
		np.subCatchmentAttributes[components.SubCatchment] =
			np.subCatchmentAttributes[components.SubCatchment].Replace(HillSlopeNitrogenContribution, deliveryAdjustedValue)
	case catchmentActions.GullySource:
		np.subCatchmentAttributes[components.SubCatchment] =
			np.subCatchmentAttributes[components.SubCatchment].Replace(GullyNitrogenContribution, value)
	default: // Deliberately does nothing
	}
}

func (np *DissolvedNitrogenProduction) calculateRiparianFineSediment(components *catchmentActions.KeyComponents, value float64) {
	if components.ElementType != catchmentActions.FineSedimentOriginalAttribute {
		return
	}

	switch components.SourceType {
	case catchmentActions.RiparianSource:
		np.subCatchmentAttributes[components.SubCatchment] =
			np.subCatchmentAttributes[components.SubCatchment].Replace(RiparianFineSediment, value)
	default: // Deliberately does nothing
	}
}

func (np *DissolvedNitrogenProduction) calculateRiparianNitrogenContributionAttribute() {
	sedimentSubCatchmentValues := np.sedimentProductionVariable.PlanningUnitAttributes()
	for subCatchment, sedimentVariableAttributes := range sedimentSubCatchmentValues {
		np.calculateRiparianNitrogenContributionForSubCatchment(subCatchment, sedimentVariableAttributes)
	}
}

func (np *DissolvedNitrogenProduction) calculateRiparianNitrogenContributionForSubCatchment(subCatchment planningunit.Id, attributes attributes.Attributes) {
	riverbankSediment := attributes.Value(sedimentproduction.RiverbankSedimentContribution).(float64)

	localAttributes := np.subCatchmentAttributes[subCatchment]
	fineSediment := localAttributes.Value(RiparianFineSediment).(float64)

	riparianNitrogen := riverbankSediment * fineSediment * conversionFactor
	np.subCatchmentAttributes[subCatchment] =
		np.subCatchmentAttributes[subCatchment].Replace(RiparianNitrogenContribution, riparianNitrogen)
}

func (np *DissolvedNitrogenProduction) calculateInitialParticulateNitrogenPerSubCatchment() {
	for subCatchment, attributes := range np.subCatchmentAttributes {
		np.updateParticulateNitrogenFor(subCatchment, attributes)
	}
}

type nitrogenContext struct {
	riparianVegetationProportion float64

	riparianContribution  float64
	hillSlopeContribution float64
	gullyContribution     float64
}

func (np *DissolvedNitrogenProduction) updateParticulateNitrogenFor(subCatchment planningunit.Id, attributes attributes.Attributes) {

	context := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		riparianContribution:         attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:            attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution:        attributes.Value(HillSlopeNitrogenContribution).(float64),
	}

	nitrogenProduced := np.calculateNitrogenProduction(context)
	np.SetPlanningUnitValue(subCatchment, nitrogenProduced)
}

func (np *DissolvedNitrogenProduction) calculateNitrogenProduction(context nitrogenContext) float64 {
	filteredHillSlopeContribution := np.deriveHillSlopeNitrogenProduction(context)
	nitrogenProduced := context.riparianContribution + context.gullyContribution + filteredHillSlopeContribution

	roundedNitrogenProduced := math.RoundFloat(nitrogenProduced, int(np.Precision()))
	return roundedNitrogenProduced
}

func (np *DissolvedNitrogenProduction) deriveHillSlopeNitrogenProduction(context nitrogenContext) float64 {
	riparianFilter := riparianBufferFilter(context.riparianVegetationProportion)
	filteredHillSlopeContribution := context.hillSlopeContribution * riparianFilter

	return filteredHillSlopeContribution
}

func Float64ToSubCatchmentId(value float64) planningunit.Id {
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

func (np *DissolvedNitrogenProduction) ObserveAction(action action.ManagementAction) {
	np.observeAction(action)
}

func (np *DissolvedNitrogenProduction) ObserveActionInitialising(action action.ManagementAction) {
	np.observeAction(action)
	np.command.Do()
}

func (np *DissolvedNitrogenProduction) observeAction(action action.ManagementAction) {
	np.actionObserved = action
	switch np.actionObserved.Type() {
	case catchmentActions.RiverBankRestorationType:
		np.handleRiverBankRestorationAction()
	case catchmentActions.GullyRestorationType:
		np.handleGullyRestorationAction()
	case catchmentActions.HillSlopeRestorationType:
		np.handleHillSlopeRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (np *DissolvedNitrogenProduction) handleRiverBankRestorationAction() {
	var asIsVegetation, asIsRiparianSediment, asIsFineSediment,
		toBeVegetation, toBeRiparianSediment, toBeFineSediment float64

	switch np.actionObserved.IsActive() {
	case true:
		asIsRiparianSediment = np.actionObserved.ModelVariableValue(catchmentActions.OriginalRiparianSedimentProduction)
		asIsFineSediment = np.actionObserved.ModelVariableValue(catchmentActions.FineSedimentOriginalAttribute)
		asIsVegetation = np.actionObserved.ModelVariableValue(catchmentActions.OriginalBufferVegetation)

		toBeRiparianSediment = np.actionObserved.ModelVariableValue(catchmentActions.ActionedRiparianSedimentProduction)
		toBeFineSediment = np.actionObserved.ModelVariableValue(catchmentActions.FineSedimentActionedAttribute)
		toBeVegetation = np.actionObserved.ModelVariableValue(catchmentActions.ActionedBufferVegetation)
	case false:
		asIsRiparianSediment = np.actionObserved.ModelVariableValue(catchmentActions.ActionedRiparianSedimentProduction)
		asIsFineSediment = np.actionObserved.ModelVariableValue(catchmentActions.FineSedimentActionedAttribute)
		asIsVegetation = np.actionObserved.ModelVariableValue(catchmentActions.ActionedBufferVegetation)

		toBeRiparianSediment = np.actionObserved.ModelVariableValue(catchmentActions.OriginalRiparianSedimentProduction)
		toBeFineSediment = np.actionObserved.ModelVariableValue(catchmentActions.FineSedimentOriginalAttribute)
		toBeVegetation = np.actionObserved.ModelVariableValue(catchmentActions.OriginalBufferVegetation)
	}

	actionSubCatchment := np.actionObserved.PlanningUnit()
	attributes := np.subCatchmentAttributes[actionSubCatchment]

	asIsRiparianNitrogen := asIsRiparianSediment * asIsFineSediment * conversionFactor

	asIsContext := nitrogenContext{
		riparianVegetationProportion: asIsVegetation,

		riparianContribution:  asIsRiparianNitrogen,
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: attributes.Value(HillSlopeNitrogenContribution).(float64),
	}

	asIsNitrogen := np.calculateNitrogenProduction(asIsContext)

	toBeRiparianNitrogen := toBeRiparianSediment * toBeFineSediment * conversionFactor

	toBeContext := nitrogenContext{
		riparianVegetationProportion: toBeVegetation,

		riparianContribution:  toBeRiparianNitrogen,
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: attributes.Value(HillSlopeNitrogenContribution).(float64),
	}

	toBeNitrogen := np.calculateNitrogenProduction(toBeContext)

	np.command = new(RiverBankRestorationCommand).
		ForVariable(np).
		InPlanningUnit(actionSubCatchment).
		WithVegetationProportion(toBeVegetation).
		WithRiverBankNitrogenContribution(toBeRiparianNitrogen).
		WithChange(toBeNitrogen - asIsNitrogen)
}

func (np *DissolvedNitrogenProduction) handleGullyRestorationAction() {
	var asIsNitrogen, toBeNitrogen float64

	switch np.actionObserved.IsActive() {
	case true:
		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
		toBeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
	case false:
		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
		toBeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
	}

	actionSubCatchment := np.actionObserved.PlanningUnit()

	np.command = new(GullyRestorationCommand).
		ForVariable(np).
		InPlanningUnit(actionSubCatchment).
		WithNitrogenContribution(toBeNitrogen).
		WithChange(toBeNitrogen - asIsNitrogen)
}

func (np *DissolvedNitrogenProduction) handleHillSlopeRestorationAction() {
	var toBeHillSlopeNitrogen, asIsHillSlopeNitrogen float64

	switch np.actionObserved.IsActive() {
	case true:
		asIsHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
		toBeHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
	case false:
		asIsHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
		toBeHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
	}

	actionSubCatchment := np.actionObserved.PlanningUnit()
	attributes := np.subCatchmentAttributes[actionSubCatchment]

	asIsContext := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),

		riparianContribution:  attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: asIsHillSlopeNitrogen,
	}

	asIsNitrogen := np.calculateNitrogenProduction(asIsContext)

	toBeContext := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),

		riparianContribution:  attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: toBeHillSlopeNitrogen,
	}

	toBeNitrogen := np.calculateNitrogenProduction(toBeContext)

	np.command = new(HillSlopeRevegetationCommand).
		ForVariable(np).
		InPlanningUnit(actionSubCatchment).
		WithFilteredNitrogenContribution(toBeHillSlopeNitrogen).
		WithChange(toBeNitrogen - asIsNitrogen)
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
func (np *DissolvedNitrogenProduction) NotifyObservers() {
	for _, observer := range np.Observers() {
		observer.ObserveDecisionVariable(np)
	}
}

func (np *DissolvedNitrogenProduction) UndoableValue() float64 {
	return np.command.Value()
}

func (np *DissolvedNitrogenProduction) SetUndoableValue(value float64) {
	np.command.SetChange(value)
}

func (np *DissolvedNitrogenProduction) DifferenceInValues() float64 {
	return np.command.Change()
}

func (np *DissolvedNitrogenProduction) ApplyDoneValue() {
	np.command.Do()
}

func (np *DissolvedNitrogenProduction) ApplyUndoneValue() {
	np.command.Undo()
}
