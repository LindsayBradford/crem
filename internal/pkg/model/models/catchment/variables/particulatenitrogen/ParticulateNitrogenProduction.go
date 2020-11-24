// Copyright (c) 2019 Australian Rivers Institute.

package particulatenitrogen

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
	VariableName = "ParticulateNitrogen"

	planningUnitIndex           = 0
	proportionOfVegetationIndex = 8

	RiverbankVegetationProportion = "RiverbankVegetationProportion"
	RiparianFineSediment          = "RiparianFineSediment"

	RiparianNitrogenContribution = "RiparianNitrogenContribution"
	GullyNitrogenContribution    = "GullyNitrogenContribution"

	WetlandRemovalEfficiency      = "WetlandRemovalEfficiency"
	HillSlopeNitrogenContribution = "HillSlopeNitrogenContribution"

	conversionFactor = 0.01
)

var _ variable.UndoableDecisionVariable = new(ParticulateNitrogenProduction)

type ParticulateNitrogenProduction struct {
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

func (np *ParticulateNitrogenProduction) Initialise(subCatchmentsTable tables.CsvTable, actionsTable tables.CsvTable, parameters catchmentParameters.Parameters) *ParticulateNitrogenProduction {
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

func (np *ParticulateNitrogenProduction) WithName(variableName string) *ParticulateNitrogenProduction {
	np.SetName(variableName)
	return np
}

func (np *ParticulateNitrogenProduction) WithStartingValue(value float64) *ParticulateNitrogenProduction {
	np.SetPlanningUnitValue(0, value)
	return np
}

func (np *ParticulateNitrogenProduction) WithObservers(observers ...variable.Observer) *ParticulateNitrogenProduction {
	np.Subscribe(observers...)
	return np
}

func (np *ParticulateNitrogenProduction) WithSedimentProductionVariable(variable *sedimentproduction.SedimentProduction) *ParticulateNitrogenProduction {
	np.sedimentProductionVariable = variable
	return np
}

func (np *ParticulateNitrogenProduction) deriveInitialState(subCatchmentsTable tables.CsvTable, parameters catchmentParameters.Parameters) {
	np.deriveNumberOfSubCatchments(subCatchmentsTable)
	np.initialiseSubCatchmentAttributes()
	np.deriveInitialNitrogen(subCatchmentsTable)
}

func (np *ParticulateNitrogenProduction) deriveNumberOfSubCatchments(subCatchmentsTable tables.CsvTable) {
	_, rowCount := subCatchmentsTable.ColumnAndRowSize()
	np.numberOfSubCatchments = rowCount
}

func (np *ParticulateNitrogenProduction) initialiseSubCatchmentAttributes() {
	np.subCatchmentAttributes = make(map[planningunit.Id]attributes.Attributes, np.numberOfSubCatchments)
	for index, _ := range np.subCatchmentAttributes {
		newAttributes := make(attributes.Attributes, 0)
		np.subCatchmentAttributes[index] = newAttributes
	}
}

func (np *ParticulateNitrogenProduction) deriveInitialNitrogen(subCatchmentsTable tables.CsvTable) {
	np.buildDefaultSubCatchmentAttributes(subCatchmentsTable)
	np.replaceDefaultAttributeValuesWithActionValues()
	np.calculateInitialParticulateNitrogenPerSubCatchment()
}

func (np *ParticulateNitrogenProduction) buildDefaultSubCatchmentAttributes(subCatchmentsTable tables.CsvTable) {
	for row := uint(0); row < np.numberOfSubCatchments; row++ {
		subCatchmentFloat64 := subCatchmentsTable.CellFloat64(planningUnitIndex, row)
		subCatchment := Float64ToSubCatchmentId(subCatchmentFloat64)

		riverBankVegetationProportion := subCatchmentsTable.CellFloat64(proportionOfVegetationIndex, row)

		np.subCatchmentAttributes[subCatchment] =
			np.subCatchmentAttributes[subCatchment].
				Add(RiverbankVegetationProportion, riverBankVegetationProportion).
				Add(RiparianFineSediment, float64(0)).
				Add(RiparianNitrogenContribution, float64(0)).
				Add(WetlandRemovalEfficiency, float64(0)).
				Add(HillSlopeNitrogenContribution, float64(0)).
				Add(GullyNitrogenContribution, float64(0))
	}
}

func (np *ParticulateNitrogenProduction) replaceDefaultAttributeValuesWithActionValues() {
	// Order below matters. Riparian nitrogen contribution depends on base attribute values being pre-calculated.
	np.calculateBaseAttributes()
	np.calculateRiparianNitrogenContributionAttribute()
}

func (np *ParticulateNitrogenProduction) calculateBaseAttributes() {
	for key, value := range np.Map() {
		components := np.DeriveMapKeyComponents(key)
		if components == nil {
			continue
		}

		np.calculateGullyAndHillSlopeContributions(components, value)
		np.calculateRiparianFineSediment(components, value)
	}
}

func (np *ParticulateNitrogenProduction) calculateGullyAndHillSlopeContributions(components *catchmentActions.KeyComponents, value float64) {
	if components.ElementType != catchmentActions.ParticulateNitrogenOriginalAttribute {
		return
	}

	switch components.Action {
	case catchmentActions.HillSlopeType:
		deliveryAdjustedValue := value * np.hillSlopeDeliveryRatio
		np.subCatchmentAttributes[components.SubCatchment] =
			np.subCatchmentAttributes[components.SubCatchment].Replace(HillSlopeNitrogenContribution, deliveryAdjustedValue)
	case catchmentActions.GullyType:
		np.subCatchmentAttributes[components.SubCatchment] =
			np.subCatchmentAttributes[components.SubCatchment].Replace(GullyNitrogenContribution, value)
	default: // Deliberately does nothing
	}
}

func (np *ParticulateNitrogenProduction) calculateRiparianFineSediment(components *catchmentActions.KeyComponents, value float64) {
	if components.ElementType != catchmentActions.FineSedimentOriginalAttribute {
		return
	}

	switch components.Action {
	case catchmentActions.RiparianType:
		np.subCatchmentAttributes[components.SubCatchment] =
			np.subCatchmentAttributes[components.SubCatchment].Replace(RiparianFineSediment, value)
	default: // Deliberately does nothing
	}
}

func (np *ParticulateNitrogenProduction) calculateRiparianNitrogenContributionAttribute() {
	sedimentSubCatchmentValues := np.sedimentProductionVariable.PlanningUnitAttributes()
	for subCatchment, sedimentVariableAttributes := range sedimentSubCatchmentValues {
		np.calculateRiparianNitrogenContributionForSubCatchment(subCatchment, sedimentVariableAttributes)
	}
}

func (np *ParticulateNitrogenProduction) calculateRiparianNitrogenContributionForSubCatchment(subCatchment planningunit.Id, attributes attributes.Attributes) {
	riverbankSediment := attributes.Value(sedimentproduction.RiverbankSedimentContribution).(float64)

	localAttributes := np.subCatchmentAttributes[subCatchment]
	fineSediment := localAttributes.Value(RiparianFineSediment).(float64)

	riparianNitrogen := riverbankSediment * fineSediment * conversionFactor
	np.subCatchmentAttributes[subCatchment] =
		np.subCatchmentAttributes[subCatchment].Replace(RiparianNitrogenContribution, riparianNitrogen)
}

func (np *ParticulateNitrogenProduction) calculateInitialParticulateNitrogenPerSubCatchment() {
	for subCatchment, attributes := range np.subCatchmentAttributes {
		np.updateParticulateNitrogenFor(subCatchment, attributes)
	}
}

type nitrogenContext struct {
	riparianVegetationProportion float64

	riparianContribution float64

	wetlandRemovalEfficiency float64
	hillSlopeContribution    float64

	gullyContribution float64
}

func (np *ParticulateNitrogenProduction) updateParticulateNitrogenFor(subCatchment planningunit.Id, attributes attributes.Attributes) {

	context := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		riparianContribution:         attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:            attributes.Value(GullyNitrogenContribution).(float64),
		wetlandRemovalEfficiency:     attributes.Value(WetlandRemovalEfficiency).(float64),
		hillSlopeContribution:        attributes.Value(HillSlopeNitrogenContribution).(float64),
	}

	nitrogenProduced := np.calculateNitrogenProduction(context)
	np.SetPlanningUnitValue(subCatchment, nitrogenProduced)
}

func (np *ParticulateNitrogenProduction) calculateNitrogenProduction(context nitrogenContext) float64 {
	filteredHillSlopeContribution := np.deriveHillSlopeNitrogenProduction(context)
	nitrogenProduced := context.riparianContribution + context.gullyContribution + filteredHillSlopeContribution

	roundedNitrogenProduced := math.RoundFloat(nitrogenProduced, int(np.Precision()))
	return roundedNitrogenProduced
}

func (np *ParticulateNitrogenProduction) deriveHillSlopeNitrogenProduction(context nitrogenContext) float64 {
	wetlandMediatedHillSlopeContribution := (1 - context.wetlandRemovalEfficiency) * context.hillSlopeContribution

	riparianFilter := riparianBufferFilter(context.riparianVegetationProportion)
	filteredHillSlopeContribution := wetlandMediatedHillSlopeContribution * riparianFilter

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

func (np *ParticulateNitrogenProduction) ObserveAction(action action.ManagementAction) {
	np.observeAction(action)
}

func (np *ParticulateNitrogenProduction) ObserveActionInitialising(action action.ManagementAction) {
	np.observeAction(action)
	np.command.Do()
}

func (np *ParticulateNitrogenProduction) observeAction(action action.ManagementAction) {
	np.actionObserved = action
	switch np.actionObserved.Type() {
	case catchmentActions.RiverBankRestorationType:
		np.handleRiverBankRestorationAction()
	case catchmentActions.GullyRestorationType:
		np.handleGullyRestorationAction()
	case catchmentActions.HillSlopeRestorationType:
		np.handleHillSlopeRestorationAction()
	case catchmentActions.WetlandsEstablishmentType:
		np.handleWetlandsEstablishmentAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (np *ParticulateNitrogenProduction) handleRiverBankRestorationAction() {
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
		wetlandRemovalEfficiency:     attributes.Value(WetlandRemovalEfficiency).(float64),

		riparianContribution:  asIsRiparianNitrogen,
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: attributes.Value(HillSlopeNitrogenContribution).(float64),
	}

	asIsNitrogen := np.calculateNitrogenProduction(asIsContext)

	toBeRiparianNitrogen := toBeRiparianSediment * toBeFineSediment * conversionFactor

	toBeContext := nitrogenContext{
		riparianVegetationProportion: toBeVegetation,
		wetlandRemovalEfficiency:     attributes.Value(WetlandRemovalEfficiency).(float64),

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

func (np *ParticulateNitrogenProduction) handleGullyRestorationAction() {
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

func (np *ParticulateNitrogenProduction) handleHillSlopeRestorationAction() {
	var asIsHillSlopeNitrogen, toBeHillSlopeNitrogen float64

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
		wetlandRemovalEfficiency:     attributes.Value(WetlandRemovalEfficiency).(float64),

		riparianContribution:  attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: asIsHillSlopeNitrogen,
	}

	asIsNitrogen := np.calculateNitrogenProduction(asIsContext)

	toBeContext := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		wetlandRemovalEfficiency:     attributes.Value(WetlandRemovalEfficiency).(float64),

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

func (np *ParticulateNitrogenProduction) handleWetlandsEstablishmentAction() {
	var asIsRemovalEfficiency, toBeRemovalEfficiency float64

	switch np.actionObserved.IsActive() {
	case true:
		asIsRemovalEfficiency = 0
		toBeRemovalEfficiency = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenRemovalEfficiency)
	case false:
		asIsRemovalEfficiency = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenRemovalEfficiency)
		toBeRemovalEfficiency = 0
	}

	actionSubCatchment := np.actionObserved.PlanningUnit()
	attributes := np.subCatchmentAttributes[actionSubCatchment]

	asIsContext := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		wetlandRemovalEfficiency:     asIsRemovalEfficiency,

		riparianContribution:  attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: asIsRemovalEfficiency,
	}

	asIsNitrogen := np.calculateNitrogenProduction(asIsContext)

	toBeContext := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		wetlandRemovalEfficiency:     toBeRemovalEfficiency,

		riparianContribution:  attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: toBeRemovalEfficiency,
	}

	toBeNitrogen := np.calculateNitrogenProduction(toBeContext)

	np.command = new(WetlandsEstablishmentCommand).
		ForVariable(np).
		InPlanningUnit(actionSubCatchment).
		WithRemovalEfficiency(toBeRemovalEfficiency).
		WithChange(toBeNitrogen - asIsNitrogen)
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
func (np *ParticulateNitrogenProduction) NotifyObservers() {
	for _, observer := range np.Observers() {
		observer.ObserveDecisionVariable(np)
	}
}

func (np *ParticulateNitrogenProduction) UndoableValue() float64 {
	return np.command.Value()
}

func (np *ParticulateNitrogenProduction) SetUndoableValue(value float64) {
	np.command.SetChange(value)
}

func (np *ParticulateNitrogenProduction) DifferenceInValues() float64 {
	return np.command.Change()
}

func (np *ParticulateNitrogenProduction) ApplyDoneValue() {
	np.command.Do()
}

func (np *ParticulateNitrogenProduction) ApplyUndoneValue() {
	np.command.Undo()
}
