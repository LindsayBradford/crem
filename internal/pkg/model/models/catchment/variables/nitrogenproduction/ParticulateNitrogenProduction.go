// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	catchmentActions "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	catchmentParameters "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
)

const (
	VariableName                = "ParticulateNitrogen"
	notImplementedValue float64 = 0

	planningUnitIndex           = 0
	proportionOfVegetationIndex = 8

	RiverbankNitrogenContribution = "RiverbankNitrogenContribution"
	GullyNitrogenContribution     = "GullyNitrogenContribution"
	HillSlopeNitrogenContribution = "HillSlopeNitrogenContribution"
	RiverbankVegetationProportion = "RiverbankVegetationProportion"
)

var _ variable.UndoableDecisionVariable = new(ParticulateNitrogenProduction)

type ParticulateNitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	catchmentActions.Container

	command variable.ChangeCommand

	actionObserved action.ManagementAction

	numberOfPlanningUnits uint

	hillSlopeDeliveryRatio float64

	hillSlopeNitrogenContribution float64
	bankNitrogenContribution      float64
	gullyNitrogenContribution     float64

	planningUnitAttributes map[planningunit.Id]attributes.Attributes
}

func (np *ParticulateNitrogenProduction) Initialise(planningUnitTable tables.CsvTable, actionsTable tables.CsvTable, parameters catchmentParameters.Parameters) *ParticulateNitrogenProduction {
	np.PerPlanningUnitDecisionVariable.Initialise()
	np.Container.WithActionsTable(actionsTable)

	np.SetName(VariableName)
	np.SetUnitOfMeasure(variable.TonnesPerYear)
	np.SetPrecision(3)

	np.hillSlopeDeliveryRatio = parameters.GetFloat64(catchmentParameters.HillSlopeDeliveryRatio)

	np.command = new(variable.NullChangeCommand)

	np.deriveInitialState(planningUnitTable, parameters)

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

func (np *ParticulateNitrogenProduction) deriveInitialState(planningUnitTable tables.CsvTable, parameters catchmentParameters.Parameters) {
	np.deriveNumberOfPlanningUnits(planningUnitTable)

	np.initialisePlanningUnitAttributes()

	np.deriveInitialNitrogen(planningUnitTable, parameters)
}

func (np *ParticulateNitrogenProduction) deriveNumberOfPlanningUnits(planningUnitTable tables.CsvTable) {
	_, rowCount := planningUnitTable.ColumnAndRowSize()
	np.numberOfPlanningUnits = rowCount
}

func (np *ParticulateNitrogenProduction) initialisePlanningUnitAttributes() {
	np.planningUnitAttributes = make(map[planningunit.Id]attributes.Attributes, np.numberOfPlanningUnits)
	for index, _ := range np.planningUnitAttributes {
		newAttributes := make(attributes.Attributes, 0)
		np.planningUnitAttributes[index] = newAttributes
	}
}

func (np *ParticulateNitrogenProduction) deriveInitialNitrogen(planningUnitTable tables.CsvTable, parameters catchmentParameters.Parameters) {
	np.buildDefaultPlanningUnitAttributes(planningUnitTable)
	np.replaceDefaultAttributeValuesWithActionVales()
	np.calculateInitialParticulateNitrogenPerPlanningUnit()
}

func (np *ParticulateNitrogenProduction) buildDefaultPlanningUnitAttributes(planningUnitTable tables.CsvTable) {
	for row := uint(0); row < np.numberOfPlanningUnits; row++ {
		planningUnitFloat64 := planningUnitTable.CellFloat64(planningUnitIndex, row)
		planningUnit := Float64ToPlanningUnitId(planningUnitFloat64)

		riverBankVegetationProportion := planningUnitTable.CellFloat64(proportionOfVegetationIndex, row)
		np.planningUnitAttributes[planningUnit] = np.planningUnitAttributes[planningUnit].Add(RiverbankVegetationProportion, riverBankVegetationProportion)
		np.planningUnitAttributes[planningUnit] = np.planningUnitAttributes[planningUnit].Add(HillSlopeNitrogenContribution, float64(0))
		np.planningUnitAttributes[planningUnit] = np.planningUnitAttributes[planningUnit].Add(GullyNitrogenContribution, float64(0))
	}
}

func (np *ParticulateNitrogenProduction) replaceDefaultAttributeValuesWithActionVales() {
	for key, value := range np.Map() {
		components := np.DeriveMapKeyComponents(key)

		if components == nil || components.ElementType != catchmentActions.ParticulateNitrogenOriginalAttribute {
			continue
		}

		switch components.SourceType {
		case catchmentActions.HillSlopeSource:
			np.planningUnitAttributes[components.SubCatchment] =
				np.planningUnitAttributes[components.SubCatchment].Replace(HillSlopeNitrogenContribution, value)
		case catchmentActions.GullySource:
			np.planningUnitAttributes[components.SubCatchment] =
				np.planningUnitAttributes[components.SubCatchment].Replace(GullyNitrogenContribution, value)
		}
	}
}

func (np *ParticulateNitrogenProduction) calculateInitialParticulateNitrogenPerPlanningUnit() {
	for subCatchment, attributes := range np.planningUnitAttributes {
		np.updateParticulateNitrogenFor(subCatchment, attributes)
	}
}

type nitrogenContext struct {
	riparianVegetationProportion float64

	hillSlopeContribution float64
	gullyContribution     float64
}

func (np *ParticulateNitrogenProduction) updateParticulateNitrogenFor(subCatchment planningunit.Id, attributes attributes.Attributes) {

	context := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		gullyContribution:            attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution:        attributes.Value(HillSlopeNitrogenContribution).(float64),
	}

	nitrogenProduced := np.calculateNitrogenProduction(context)
	np.SetPlanningUnitValue(subCatchment, nitrogenProduced)
}

func (np *ParticulateNitrogenProduction) calculateNitrogenProduction(context nitrogenContext) float64 {
	deliveryAdjustedHillSlopeContribution := context.hillSlopeContribution * np.hillSlopeDeliveryRatio

	riparianFilter := riparianBufferFilter(context.riparianVegetationProportion)
	filteredHillSlopeContribution := deliveryAdjustedHillSlopeContribution * riparianFilter

	nitrogenProduced := context.gullyContribution + filteredHillSlopeContribution

	roundedNitrogenProduced := math.RoundFloat(nitrogenProduced, int(np.Precision()))
	return roundedNitrogenProduced
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
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (np *ParticulateNitrogenProduction) handleRiverBankRestorationAction() {
	var toBeVegetation float64

	switch np.actionObserved.IsActive() {
	case true:
		toBeVegetation = np.actionObserved.ModelVariableValue(catchmentActions.ActionedBufferVegetation)
	case false:
		toBeVegetation = np.actionObserved.ModelVariableValue(catchmentActions.OriginalBufferVegetation)
	}

	np.command = new(RiverBankRestorationCommand).
		ForVariable(np).
		InPlanningUnit(np.actionObserved.PlanningUnit()).
		WithVegetationProportion(toBeVegetation).
		WithChange(0)
}

func (np *ParticulateNitrogenProduction) handleGullyRestorationAction() {
	var toBeNitrogen, asIsNitrogen float64

	switch np.actionObserved.IsActive() {
	case true:
		toBeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
	case false:
		toBeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
	}

	np.command = new(GullyRestorationCommand).
		ForVariable(np).
		InPlanningUnit(np.actionObserved.PlanningUnit()).
		WithChange(toBeNitrogen - asIsNitrogen)
}

func (np *ParticulateNitrogenProduction) handleHillSlopeRestorationAction() {
	var toBeHillSlopeNitrogen, asIsHillSlopeNitrogen float64

	switch np.actionObserved.IsActive() {
	case true:
		toBeHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
		asIsHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
	case false:
		toBeHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenOriginalAttribute)
		asIsHillSlopeNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.ParticulateNitrogenActionedAttribute)
	}

	attributes := np.planningUnitAttributes[np.actionObserved.PlanningUnit()]

	asIsContext := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		gullyContribution:            attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution:        asIsHillSlopeNitrogen,
	}

	asIsNitrogen := np.calculateNitrogenProduction(asIsContext)

	toBeContext := nitrogenContext{
		riparianVegetationProportion: attributes.Value(RiverbankVegetationProportion).(float64),
		gullyContribution:            attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution:        toBeHillSlopeNitrogen,
	}

	toBeNitrogen := np.calculateNitrogenProduction(toBeContext)

	np.command = new(HillSlopeRevegetationCommand).
		ForVariable(np).
		InPlanningUnit(np.actionObserved.PlanningUnit()).
		WithNitrogenContribution(toBeHillSlopeNitrogen).
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
