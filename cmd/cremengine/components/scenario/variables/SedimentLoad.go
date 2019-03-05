// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/actions"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/pkg/errors"
)

const SedimentLoadVariableName = "SedimentLoad"

var _ variable.DecisionVariable = new(SedimentLoad)

type SedimentLoad struct {
	variable.BaseInductiveDecisionVariable
	bankSedimentContribution actions.BankSedimentContribution
	actionObserved           action.ManagementAction
}

func (sl *SedimentLoad) Initialise(planningUnitTable *tables.CsvTable, parameters parameters.Parameters) *SedimentLoad {
	sl.SetName(SedimentLoadVariableName)
	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.SetValue(sl.deriveInitialSedimentLoad())
	return sl
}

func (sl *SedimentLoad) deriveInitialSedimentLoad() float64 {
	return sl.bankSedimentContribution.OriginalSedimentContribution() +
		sl.gullySedimentContribution() +
		sl.hillSlopeSedimentContribution()
}

func (sl *SedimentLoad) gullySedimentContribution() float64 {
	return 0 // TODO: implement
}

func (sl *SedimentLoad) hillSlopeSedimentContribution() float64 {
	return 0 // TODO: implement
}

func (sl *SedimentLoad) ObserveAction(action action.ManagementAction) {
	sl.actionObserved = action
	switch sl.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		sl.handleRiverBankRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (sl *SedimentLoad) ObserveActionInitialising(action action.ManagementAction) {
	sl.actionObserved = action
	switch sl.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		sl.handleInitialisingRiverBankRestorationAction()
	default:
		panic(errors.New("Unhandled observation of initialising management action type [" + string(action.Type()) + "]"))
	}
	sl.NotifyObservers()
}

func (sl *SedimentLoad) handleRiverBankRestorationAction() {
	setTempVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		asIsSedimentContribution := sl.bankSedimentContribution.SedimentImpactOfRiparianVegetation(asIsVegetation)
		toBeSedimentContribution := sl.bankSedimentContribution.SedimentImpactOfRiparianVegetation(toBeVegetation)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)
	}

	switch sl.actionObserved.IsActive() {
	case true:
		setTempVariable(actions.OriginalBufferVegetation, actions.ActionedBufferVegetation)
	case false:
		setTempVariable(actions.ActionedBufferVegetation, actions.OriginalBufferVegetation)
	}
}

func (sl *SedimentLoad) handleInitialisingRiverBankRestorationAction() {
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		asIsSedimentContribution := sl.bankSedimentContribution.SedimentImpactOfRiparianVegetation(asIsVegetation)
		toBeSedimentContribution := sl.bankSedimentContribution.SedimentImpactOfRiparianVegetation(toBeVegetation)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)
	}

	switch sl.actionObserved.IsActive() {
	case true:
		setVariable(actions.OriginalBufferVegetation, actions.ActionedBufferVegetation)
	case false:
		setVariable(actions.ActionedBufferVegetation, actions.OriginalBufferVegetation)
	}
}
