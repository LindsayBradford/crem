// Copyright (c) 2019 Australian Rivers Institute.

package dissolvednitrogen

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
	VariableName = "DissolvedNitrogen"

	planningUnitIndex = 0

	RiparianNitrogenContribution  = "RiparianNitrogenContribution"
	GullyNitrogenContribution     = "GullyNitrogenContribution"
	HillSlopeNitrogenContribution = "HillSlopeNitrogenContribution"
)

var _ variable.UndoableDecisionVariable = new(DissolvedNitrogenProduction)

type DissolvedNitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	catchmentActions.Container

	command variable.ChangeCommand

	actionObserved action.ManagementAction

	numberOfSubCatchments uint

	subCatchmentAttributes map[planningunit.Id]attributes.Attributes
}

func (dn *DissolvedNitrogenProduction) Initialise(subCatchmentsTable tables.CsvTable, actionsTable tables.CsvTable, parameters catchmentParameters.Parameters) *DissolvedNitrogenProduction {
	dn.PerPlanningUnitDecisionVariable.Initialise()
	dn.Container.WithActionsTable(actionsTable)

	dn.SetName(VariableName)
	dn.SetUnitOfMeasure(variable.TonnesPerYear)
	dn.SetPrecision(3)

	dn.command = new(variable.NullChangeCommand)

	dn.deriveInitialState(subCatchmentsTable, parameters)

	return dn
}

func (dn *DissolvedNitrogenProduction) WithName(variableName string) *DissolvedNitrogenProduction {
	dn.SetName(variableName)
	return dn
}

func (dn *DissolvedNitrogenProduction) WithStartingValue(value float64) *DissolvedNitrogenProduction {
	dn.SetPlanningUnitValue(0, value)
	return dn
}

func (dn *DissolvedNitrogenProduction) WithObservers(observers ...variable.Observer) *DissolvedNitrogenProduction {
	dn.Subscribe(observers...)
	return dn
}

func (dn *DissolvedNitrogenProduction) deriveInitialState(subCatchmentsTable tables.CsvTable, parameters catchmentParameters.Parameters) {
	dn.deriveNumberOfSubCatchments(subCatchmentsTable)
	dn.initialiseSubCatchmentAttributes()
	dn.deriveInitialNitrogen(subCatchmentsTable)
}

func (dn *DissolvedNitrogenProduction) deriveNumberOfSubCatchments(subCatchmentsTable tables.CsvTable) {
	_, rowCount := subCatchmentsTable.ColumnAndRowSize()
	dn.numberOfSubCatchments = rowCount
}

func (dn *DissolvedNitrogenProduction) initialiseSubCatchmentAttributes() {
	dn.subCatchmentAttributes = make(map[planningunit.Id]attributes.Attributes, dn.numberOfSubCatchments)
	for index, _ := range dn.subCatchmentAttributes {
		newAttributes := make(attributes.Attributes, 0)
		dn.subCatchmentAttributes[index] = newAttributes
	}
}

func (dn *DissolvedNitrogenProduction) deriveInitialNitrogen(subCatchmentsTable tables.CsvTable) {
	dn.buildDefaultSubCatchmentAttributes(subCatchmentsTable)
	dn.replaceDefaultAttributeValuesWithActionOriginalValues()
	dn.calculateInitialNitrogenPerSubCatchment()
}

func (dn *DissolvedNitrogenProduction) buildDefaultSubCatchmentAttributes(subCatchmentsTable tables.CsvTable) {
	for row := uint(0); row < dn.numberOfSubCatchments; row++ {
		subCatchmentFloat64 := subCatchmentsTable.CellFloat64(planningUnitIndex, row)
		subCatchment := Float64ToSubCatchmentId(subCatchmentFloat64)

		dn.subCatchmentAttributes[subCatchment] =
			dn.subCatchmentAttributes[subCatchment].
				Add(RiparianNitrogenContribution, float64(0)).
				Add(HillSlopeNitrogenContribution, float64(0)).
				Add(GullyNitrogenContribution, float64(0))
	}
}

func (dn *DissolvedNitrogenProduction) replaceDefaultAttributeValuesWithActionOriginalValues() {
	for key, value := range dn.Map() {
		components := dn.DeriveMapKeyComponents(key)
		if components == nil {
			continue
		}

		dn.calculateOriginalDissolvedNitrogenContributions(components, value)
	}
}

func (dn *DissolvedNitrogenProduction) calculateOriginalDissolvedNitrogenContributions(components *catchmentActions.KeyComponents, value float64) {
	if components.ElementType != catchmentActions.DissolvedNitrogenOriginalAttribute {
		return
	}

	switch components.Action {
	case catchmentActions.RiparianType:
		dn.subCatchmentAttributes[components.SubCatchment] =
			dn.subCatchmentAttributes[components.SubCatchment].Replace(RiparianNitrogenContribution, value)
	case catchmentActions.HillSlopeType:
		dn.subCatchmentAttributes[components.SubCatchment] =
			dn.subCatchmentAttributes[components.SubCatchment].Replace(HillSlopeNitrogenContribution, value)
	case catchmentActions.GullyType:
		dn.subCatchmentAttributes[components.SubCatchment] =
			dn.subCatchmentAttributes[components.SubCatchment].Replace(GullyNitrogenContribution, value)
	default: // Deliberately does nothing
	}
}

func (dn *DissolvedNitrogenProduction) calculateInitialNitrogenPerSubCatchment() {
	for subCatchment, attributes := range dn.subCatchmentAttributes {
		dn.updateParticulateNitrogenFor(subCatchment, attributes)
	}
}

type nitrogenContext struct {
	riparianContribution  float64
	hillSlopeContribution float64
	gullyContribution     float64
}

func (dn *DissolvedNitrogenProduction) updateParticulateNitrogenFor(subCatchment planningunit.Id, attributes attributes.Attributes) {

	context := nitrogenContext{
		riparianContribution:  attributes.Value(RiparianNitrogenContribution).(float64),
		gullyContribution:     attributes.Value(GullyNitrogenContribution).(float64),
		hillSlopeContribution: attributes.Value(HillSlopeNitrogenContribution).(float64),
	}

	nitrogenProduced := dn.calculateNitrogenProduction(context)
	dn.SetPlanningUnitValue(subCatchment, nitrogenProduced)
}

func (dn *DissolvedNitrogenProduction) calculateNitrogenProduction(context nitrogenContext) float64 {
	nitrogenProduced := context.riparianContribution + context.gullyContribution + context.hillSlopeContribution
	roundedNitrogenProduced := math.RoundFloat(nitrogenProduced, int(dn.Precision()))
	return roundedNitrogenProduced
}

func Float64ToSubCatchmentId(value float64) planningunit.Id {
	return planningunit.Id(value)
}

func (dn *DissolvedNitrogenProduction) ObserveAction(action action.ManagementAction) {
	dn.observeAction(action)
}

func (dn *DissolvedNitrogenProduction) ObserveActionInitialising(action action.ManagementAction) {
	dn.observeAction(action)
	dn.command.Do()
}

func (dn *DissolvedNitrogenProduction) observeAction(action action.ManagementAction) {
	dn.actionObserved = action
	switch dn.actionObserved.Type() {
	case catchmentActions.RiverBankRestorationType:
		dn.handleRiverBankRestorationAction()
	case catchmentActions.GullyRestorationType:
		dn.handleGullyRestorationAction()
	case catchmentActions.HillSlopeRestorationType:
		dn.handleHillSlopeRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (dn *DissolvedNitrogenProduction) handleRiverBankRestorationAction() {
	var asIsNitrogen, toBeNitrogen float64

	switch dn.actionObserved.IsActive() {
	case true:
		asIsNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenOriginalAttribute)
		toBeNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenActionedAttribute)
	case false:
		asIsNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenActionedAttribute)
		toBeNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenOriginalAttribute)
	}

	actionSubCatchment := dn.actionObserved.PlanningUnit()

	dn.command = new(RiverBankRestorationCommand).
		ForVariable(dn).
		InPlanningUnit(actionSubCatchment).
		WithNitrogenContribution(toBeNitrogen).
		WithChange(toBeNitrogen - asIsNitrogen)
}

func (dn *DissolvedNitrogenProduction) handleGullyRestorationAction() {
	var asIsNitrogen, toBeNitrogen float64

	switch dn.actionObserved.IsActive() {
	case true:
		asIsNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenOriginalAttribute)
		toBeNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenActionedAttribute)
	case false:
		asIsNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenActionedAttribute)
		toBeNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenOriginalAttribute)
	}

	actionSubCatchment := dn.actionObserved.PlanningUnit()

	dn.command = new(GullyRestorationCommand).
		ForVariable(dn).
		InPlanningUnit(actionSubCatchment).
		WithNitrogenContribution(toBeNitrogen).
		WithChange(toBeNitrogen - asIsNitrogen)
}

func (dn *DissolvedNitrogenProduction) handleHillSlopeRestorationAction() {
	var asIsNitrogen, toBeNitrogen float64

	switch dn.actionObserved.IsActive() {
	case true:
		asIsNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenOriginalAttribute)
		toBeNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenActionedAttribute)
	case false:
		asIsNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenActionedAttribute)
		toBeNitrogen = dn.actionObserved.ModelVariableValue(catchmentActions.DissolvedNitrogenOriginalAttribute)
	}

	actionSubCatchment := dn.actionObserved.PlanningUnit()

	dn.command = new(HillSlopeRevegetationCommand).
		ForVariable(dn).
		InPlanningUnit(actionSubCatchment).
		WithNitrogenContribution(toBeNitrogen).
		WithChange(toBeNitrogen - asIsNitrogen)
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
func (dn *DissolvedNitrogenProduction) NotifyObservers() {
	for _, observer := range dn.Observers() {
		observer.ObserveDecisionVariable(dn)
	}
}

func (dn *DissolvedNitrogenProduction) UndoableValue() float64 {
	return dn.command.Value()
}

func (dn *DissolvedNitrogenProduction) SetUndoableValue(value float64) {
	dn.command.SetChange(value)
}

func (dn *DissolvedNitrogenProduction) DifferenceInValues() float64 {
	return dn.command.Change()
}

func (dn *DissolvedNitrogenProduction) ApplyDoneValue() {
	dn.command.Do()
}

func (dn *DissolvedNitrogenProduction) ApplyUndoneValue() {
	dn.command.Undo()
}
