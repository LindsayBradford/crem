// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	actions2 "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction2"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/pkg/errors"
)

const VariableName = "NitrogenProduction"
const notImplementedValue float64 = 0

var _ variable.UndoableDecisionVariable = new(NitrogenProduction)

type NitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable

	sedimentProductionVariable *sedimentproduction2.SedimentProduction2

	command variable.ChangeCommand

	actionObserved action.ManagementAction
}

func (np *NitrogenProduction) Initialise() *NitrogenProduction {
	np.PerPlanningUnitDecisionVariable.Initialise()

	np.SetName(VariableName)
	np.SetUnitOfMeasure(variable.TonnesPerYear)
	np.SetPrecision(3)

	np.command = new(variable.NullChangeCommand)

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

func (np *NitrogenProduction) deriveInitialValue() float64 {
	np.SetValue(notImplementedValue)
	return notImplementedValue
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
	case actions2.RiverBankRestorationType:
		np.handleRiverBankRestorationAction()
	case actions2.GullyRestorationType:
		np.handleGullyRestorationAction()
	case actions2.HillSlopeRestorationType:
		np.handleHillSlopeRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (np *NitrogenProduction) handleRiverBankRestorationAction() {
	// TODO: Implement
}

func (np *NitrogenProduction) handleGullyRestorationAction() {
	// TODO: Implement
}

func (np *NitrogenProduction) handleHillSlopeRestorationAction() {
	// TODO: Implement
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
