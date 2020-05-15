// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	catchmentActions "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
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

	hillSlopeNitrogenContribution float64
	bankNitrogenContribution      float64
	gullyNitrogenContribution     float64
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
	np.hillSlopeNitrogenContribution = np.deriveInitialHillSlopeContribution()
	np.bankNitrogenContribution = np.deriveInitialBankContribution()
	np.gullyNitrogenContribution = np.deriveInitialGullyContribution()

	totalNitrogenContribution :=
		np.hillSlopeNitrogenContribution + np.bankNitrogenContribution + np.gullyNitrogenContribution

	return totalNitrogenContribution
}

func (np *NitrogenProduction) deriveInitialHillSlopeContribution() float64 {
	return 0
}

func (np *NitrogenProduction) deriveInitialBankContribution() float64 {
	return 0
}

func (np *NitrogenProduction) deriveInitialGullyContribution() float64 {
	return 0
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
