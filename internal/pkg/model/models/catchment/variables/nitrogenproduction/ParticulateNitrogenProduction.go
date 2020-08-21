// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	catchmentActions "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	catchmentParameters "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/pkg/errors"
)

const VariableName = "ParticulateNitrogen"
const notImplementedValue float64 = 0

var _ variable.UndoableDecisionVariable = new(ParticulateNitrogenProduction)

type ParticulateNitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	catchmentActions.Container

	command variable.ChangeCommand

	actionObserved action.ManagementAction

	hillSlopeNitrogenContribution float64
	bankNitrogenContribution      float64
	gullyNitrogenContribution     float64
}

func (np *ParticulateNitrogenProduction) Initialise(actionsTable tables.CsvTable, parameters catchmentParameters.Parameters) *ParticulateNitrogenProduction {
	np.PerPlanningUnitDecisionVariable.Initialise()
	np.Container.WithActionsTable(actionsTable)

	np.SetName(VariableName)
	np.SetUnitOfMeasure(variable.TonnesPerYear)
	np.SetPrecision(3)

	np.command = new(variable.NullChangeCommand)

	np.deriveInitialState(parameters)

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

func (np *ParticulateNitrogenProduction) deriveInitialState(parameters catchmentParameters.Parameters) {
	hillSlopeDeliveryRatio := parameters.GetFloat64(catchmentParameters.HillSlopeDeliveryRatio)
	for key, value := range np.Map() {
		components := np.DeriveMapKeyComponents(key)
		if components == nil || components.ElementType != catchmentActions.ParticulateNitrogenOriginalAttribute {
			continue
		}

		if components.SourceType == catchmentActions.HillSlopeSource {
			value = value * hillSlopeDeliveryRatio
		}

		// TODO: this approach doesn't cater to riparian filtering
		currentValue := np.PlanningUnitValue(components.SubCatchment)
		newValue := currentValue + value
		np.SetPlanningUnitValue(components.SubCatchment, newValue)
	}
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
	//TODO: This doesn't handle riparian buffer filtering dependency.
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
	//TODO: This doesn't handle riparian buffer filtering dependency.
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
