// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	catchmentActions "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction2"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/pkg/errors"
	math2 "math"
)

const VariableName = "NitrogenProduction"
const notImplementedValue float64 = 0

type particulateNitrogenVariables struct {
	sediment      float64
	totalCarbon   float64
	totalNitrogen float64
}

var _ variable.UndoableDecisionVariable = new(NitrogenProduction)

type NitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable

	sedimentProductionVariable *sedimentproduction2.SedimentProduction2
	catchmentActions.ParentSoilsContainer

	command variable.ChangeCommand

	actionObserved action.ManagementAction

	hillSlopeNitrogenContribution float64
	bankNitrogenContribution      float64
	gullyNitrogenContribution     float64
}

func (np *NitrogenProduction) Initialise(parentSoilsTable tables.CsvTable) *NitrogenProduction {
	np.PerPlanningUnitDecisionVariable.Initialise()
	np.ParentSoilsContainer.WithParentSoilsTable(parentSoilsTable)

	np.SetName(VariableName)
	np.SetUnitOfMeasure(variable.TonnesPerYear)
	np.SetPrecision(3)

	np.command = new(variable.NullChangeCommand)

	np.deriveInitialState()

	return np
}

func (np *NitrogenProduction) WithName(variableName string) *NitrogenProduction {
	np.SetName(variableName)
	return np
}

func (np *NitrogenProduction) WithStartingValue(value float64) *NitrogenProduction {
	np.SetPlanningUnitValue(0, value)
	return np
}

func (np *NitrogenProduction) WithSedimentProductionVariable(variable *sedimentproduction2.SedimentProduction2) *NitrogenProduction {
	np.sedimentProductionVariable = variable
	return np
}

func (np *NitrogenProduction) WithObservers(observers ...variable.Observer) *NitrogenProduction {
	np.Subscribe(observers...)
	return np
}

func (np *NitrogenProduction) deriveInitialState() {
	sedimentPlanningUnitValues := np.sedimentProductionVariable.PlanningUnitAttributes()

	for planningUnit, attributes := range sedimentPlanningUnitValues {
		initialHillSlopeSediment := attributes.Value(sedimentproduction2.HillSlopeSedimentContribution).(float64)
		riverbankSediment := attributes.Value(sedimentproduction2.RiverbankSedimentContribution).(float64)
		initialGullySediment := attributes.Value(sedimentproduction2.GullySedimentContribution).(float64)

		initialNitrogen := np.initialHillSlopeSediment(planningUnit, initialHillSlopeSediment) +
			np.initialRiverbankNitrogen(planningUnit, riverbankSediment) +
			np.initialGullyNitrogen(planningUnit, initialGullySediment)

		roundedNitrogen := math.RoundFloat(initialNitrogen, int(np.Precision()))

		np.SetPlanningUnitValue(planningUnit, roundedNitrogen)
	}
}

func (np *NitrogenProduction) initialHillSlopeSediment(planningUnit planningunit.Id, initialHillSlopeSediment float64) float64 {
	if initialHillSlopeSediment == 0 {
		return 0
	}

	carbonKey := np.DeriveMapKey(planningUnit, catchmentActions.HillSlopeSource, catchmentActions.CarbonAttribute)
	nitrogenKey := np.DeriveMapKey(planningUnit, catchmentActions.HillSlopeSource, catchmentActions.NitrogenAttribute)

	variables := particulateNitrogenVariables{
		sediment:      initialHillSlopeSediment,
		totalCarbon:   np.MapValue(carbonKey),
		totalNitrogen: np.MapValue(nitrogenKey),
	}

	calculatedParticulateNitrogen := calculateParticulateNitrogen(variables)
	return calculatedParticulateNitrogen
}

func (np *NitrogenProduction) initialRiverbankNitrogen(planningUnit planningunit.Id, initialRiverbankSediment float64) float64 {
	if initialRiverbankSediment == 0 {
		return 0
	}

	carbonKey := np.DeriveMapKey(planningUnit, catchmentActions.RiparianSource, catchmentActions.CarbonAttribute)
	nitrogenKey := np.DeriveMapKey(planningUnit, catchmentActions.RiparianSource, catchmentActions.NitrogenAttribute)

	variables := particulateNitrogenVariables{
		sediment:      initialRiverbankSediment,
		totalCarbon:   np.ParentSoilsContainer.MapValue(carbonKey),
		totalNitrogen: np.ParentSoilsContainer.MapValue(nitrogenKey),
	}

	calculatedParticulateNitrogen := calculateParticulateNitrogen(variables)
	return calculatedParticulateNitrogen
}

func (np *NitrogenProduction) initialGullyNitrogen(planningUnit planningunit.Id, initialGullySediment float64) float64 {
	if initialGullySediment == 0 {
		return 0
	}

	carbonKey := np.DeriveMapKey(planningUnit, catchmentActions.GullySource, catchmentActions.CarbonAttribute)
	nitrogenKey := np.DeriveMapKey(planningUnit, catchmentActions.GullySource, catchmentActions.NitrogenAttribute)

	variables := particulateNitrogenVariables{
		sediment:      initialGullySediment,
		totalCarbon:   np.ParentSoilsContainer.MapValue(carbonKey),
		totalNitrogen: np.ParentSoilsContainer.MapValue(nitrogenKey),
	}

	calculatedParticulateNitrogen := calculateParticulateNitrogen(variables)
	return calculatedParticulateNitrogen
}

func calculateParticulateNitrogen(variables particulateNitrogenVariables) float64 {
	const logEnrichmentRatio = 0.8
	totalNitrogenParentSoil := 0.08*variables.totalCarbon - 0.007*(variables.totalCarbon/variables.totalNitrogen) + 0.09
	scaleAdjustedTotalNitrogen := 0.01 * totalNitrogenParentSoil

	assert.That(scaleAdjustedTotalNitrogen > 0).WithFailureMessage("totalNitrogen not positive").Holds()

	particulateNitrogen := math2.Pow(scaleAdjustedTotalNitrogen, logEnrichmentRatio)
	return particulateNitrogen
}

func (np *NitrogenProduction) deriveInitialHillSlopeContribution() float64 {
	sedimentPlanningUnitValues := np.sedimentProductionVariable.PlanningUnitAttributes()

	initialHillSlopeContribution := float64(0)
	for _, attributes := range sedimentPlanningUnitValues {
		initialHillSlopeContribution += attributes.Value(sedimentproduction2.HillSlopeSedimentContribution).(float64)
	}

	return initialHillSlopeContribution
}

func (np *NitrogenProduction) deriveInitialBankContribution() float64 {
	sedimentPlanningUnitValues := np.sedimentProductionVariable.PlanningUnitAttributes()

	riverbankContribution := float64(0)
	for _, attributes := range sedimentPlanningUnitValues {
		riverbankContribution += attributes.Value(sedimentproduction2.RiverbankSedimentContribution).(float64)
	}

	return riverbankContribution
}

func (np *NitrogenProduction) deriveInitialGullyContribution() float64 {
	sedimentPlanningUnitValues := np.sedimentProductionVariable.PlanningUnitAttributes()

	initialGullyContribution := float64(0)
	for _, attributes := range sedimentPlanningUnitValues {
		initialGullyContribution += attributes.Value(sedimentproduction2.GullySedimentContribution).(float64)
	}

	return initialGullyContribution
}

func (np *NitrogenProduction) ObserveAction(action action.ManagementAction) {
	np.observeAction(action)
}

func (np *NitrogenProduction) ObserveActionInitialising(action action.ManagementAction) {
	np.observeAction(action)
	np.command.Do()
}

func (np *NitrogenProduction) observeAction(action action.ManagementAction) {
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

func (np *NitrogenProduction) handleRiverBankRestorationAction() {
	sedimentPlanningUnitValues := np.sedimentProductionVariable.PlanningUnitAttributes()
	attributes := sedimentPlanningUnitValues[np.actionObserved.PlanningUnit()]

	var asIsCarbon, toBeCarbon, asIsNitrogen, toBeNitrogen float64
	switch np.actionObserved.IsActive() {
	case true:
		asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
		toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)

		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
		toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	case false:
		asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
		toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)

		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
		toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	}

	asIsSediment := attributes.Value(sedimentproduction2.RiverbankSedimentContribution).(float64)
	toBeSediment := asIsSediment + np.sedimentProductionVariable.DifferenceInValues()

	asIsVariables := particulateNitrogenVariables{
		sediment:      asIsSediment,
		totalCarbon:   asIsCarbon,
		totalNitrogen: asIsNitrogen,
	}

	asIsParticulateNitrogen := calculateParticulateNitrogen(asIsVariables)

	toBeVariables := particulateNitrogenVariables{
		sediment:      toBeSediment,
		totalCarbon:   toBeCarbon,
		totalNitrogen: toBeNitrogen,
	}

	toBeParticulateNitrogen := calculateParticulateNitrogen(toBeVariables)

	np.command = new(RiverBankRestorationCommand).
		ForVariable(np).
		InPlanningUnit(np.actionObserved.PlanningUnit()).
		WithChange(toBeParticulateNitrogen - asIsParticulateNitrogen)
}

func (np *NitrogenProduction) handleGullyRestorationAction() {
	sedimentPlanningUnitValues := np.sedimentProductionVariable.PlanningUnitAttributes()
	attributes := sedimentPlanningUnitValues[np.actionObserved.PlanningUnit()]

	var asIsCarbon, toBeCarbon, asIsNitrogen, toBeNitrogen float64
	switch np.actionObserved.IsActive() {
	case true:
		asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
		toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)

		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
		toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	case false:
		asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
		toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)

		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
		toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	}

	asIsSediment := attributes.Value(sedimentproduction2.GullySedimentContribution).(float64)
	toBeSediment := asIsSediment + np.sedimentProductionVariable.DifferenceInValues()

	asIsVariables := particulateNitrogenVariables{
		sediment:      asIsSediment,
		totalCarbon:   asIsCarbon,
		totalNitrogen: asIsNitrogen,
	}

	asIsParticulateNitrogen := calculateParticulateNitrogen(asIsVariables)

	toBeVariables := particulateNitrogenVariables{
		sediment:      toBeSediment,
		totalCarbon:   toBeCarbon,
		totalNitrogen: toBeNitrogen,
	}

	toBeParticulateNitrogen := calculateParticulateNitrogen(toBeVariables)

	np.command = new(GullyRestorationCommand).
		ForVariable(np).
		InPlanningUnit(np.actionObserved.PlanningUnit()).
		WithChange(toBeParticulateNitrogen - asIsParticulateNitrogen)
}

func (np *NitrogenProduction) handleHillSlopeRestorationAction() {
	sedimentPlanningUnitValues := np.sedimentProductionVariable.PlanningUnitAttributes()
	attributes := sedimentPlanningUnitValues[np.actionObserved.PlanningUnit()]

	var asIsCarbon, toBeCarbon, asIsNitrogen, toBeNitrogen float64
	switch np.actionObserved.IsActive() {
	case true:
		asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
		toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)

		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
		toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	case false:
		asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
		toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)

		asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
		toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	}

	asIsSediment := attributes.Value(sedimentproduction2.HillSlopeSedimentContribution).(float64)
	toBeSediment := asIsSediment + np.sedimentProductionVariable.DifferenceInValues()

	asIsVariables := particulateNitrogenVariables{
		sediment:      asIsSediment,
		totalCarbon:   asIsCarbon,
		totalNitrogen: asIsNitrogen,
	}

	asIsParticulateNitrogen := calculateParticulateNitrogen(asIsVariables)

	toBeVariables := particulateNitrogenVariables{
		sediment:      toBeSediment,
		totalCarbon:   toBeCarbon,
		totalNitrogen: toBeNitrogen,
	}

	toBeParticulateNitrogen := calculateParticulateNitrogen(toBeVariables)

	np.command = new(HillSlopeRevegetationCommand).
		ForVariable(np).
		InPlanningUnit(np.actionObserved.PlanningUnit()).
		WithChange(toBeParticulateNitrogen - asIsParticulateNitrogen)
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
func (np *NitrogenProduction) NotifyObservers() {
	for _, observer := range np.Observers() {
		observer.ObserveDecisionVariable(np)
	}
}

func (np *NitrogenProduction) UndoableValue() float64 {
	return np.command.Value()
}

func (np *NitrogenProduction) SetUndoableValue(value float64) {
	np.command.SetChange(value)
}

func (np *NitrogenProduction) DifferenceInValues() float64 {
	return np.command.Change()
}

func (np *NitrogenProduction) ApplyDoneValue() {
	np.command.Do()
}

func (np *NitrogenProduction) ApplyUndoneValue() {
	np.command.Undo()
}
