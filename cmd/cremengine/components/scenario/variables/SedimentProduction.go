// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/actions"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/pkg/assert/release"
	"github.com/pkg/errors"
)

const SedimentProductionVariableName = "SedimentProduction"

var _ variable.DecisionVariable = new(SedimentProduction)

type SedimentProduction struct {
	variable.BaseInductiveDecisionVariable

	bankSedimentContribution  actions.BankSedimentContribution
	gullySedimentContribution actions.GullySedimentContribution

	actionObserved action.ManagementAction
}

func (sl *SedimentProduction) Initialise(planningUnitTable tables.CsvTable, gulliesTable tables.CsvTable, parameters parameters.Parameters) *SedimentProduction {
	sl.SetName(SedimentProductionVariableName)
	sl.SetUnitOfMeasure(variable.TonnesPerYear)
	sl.SetPrecision(3)
	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.gullySedimentContribution.Initialise(gulliesTable, parameters)
	sl.SetValue(sl.deriveInitialSedimentLoad())
	return sl
}

func (sl *SedimentProduction) WithObservers(observers ...variable.Observer) *SedimentProduction {
	sl.Subscribe(observers...)
	return sl
}

func (sl *SedimentProduction) deriveInitialSedimentLoad() float64 {
	return sl.bankSedimentContribution.OriginalSedimentContribution() +
		sl.gullySedimentContribution.OriginalSedimentContribution() +
		sl.hillSlopeSedimentContribution()
}

func (sl *SedimentProduction) hillSlopeSedimentContribution() float64 {
	return 0 // TODO: implement
}

func (sl *SedimentProduction) ObserveAction(action action.ManagementAction) {
	sl.actionObserved = action
	switch sl.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		sl.handleRiverBankRestorationAction()
	case actions.GullyRestorationType:
		sl.handleGullyRestorationAction()
	default:
		panic(errors.New("Unhandled observation of management action type [" + string(action.Type()) + "]"))
	}
}

func (sl *SedimentProduction) ObserveActionInitialising(action action.ManagementAction) {
	sl.actionObserved = action
	switch sl.actionObserved.Type() {
	case actions.RiverBankRestorationType:
		sl.handleInitialisingRiverBankRestorationAction()
	case actions.GullyRestorationType:
		sl.handleInitialisingGullyRestorationAction()
	default:
		panic(errors.New("Unhandled observation of initialising management action type [" + string(action.Type()) + "]"))
	}
	sl.NotifyObservers()
}

func (sl *SedimentProduction) handleRiverBankRestorationAction() {
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

func (sl *SedimentProduction) handleInitialisingRiverBankRestorationAction() {
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVegetation := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVegetation := sl.actionObserved.ModelVariableValue(toBeName)

		asIsSedimentContribution := sl.bankSedimentContribution.SedimentImpactOfRiparianVegetation(asIsVegetation)
		toBeSedimentContribution := sl.bankSedimentContribution.SedimentImpactOfRiparianVegetation(toBeVegetation)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)
	}

	assert.That(sl.actionObserved.IsActive()).WithFailureMessage("initialising action should always be active").Holds()
	setVariable(actions.OriginalBufferVegetation, actions.ActionedBufferVegetation)
}

func (sl *SedimentProduction) handleGullyRestorationAction() {
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVolume := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVolume := sl.actionObserved.ModelVariableValue(toBeName)

		asIsSedimentContribution := sl.gullySedimentContribution.SedimentFromVolume(asIsVolume)
		toBeSedimentContribution := sl.gullySedimentContribution.SedimentFromVolume(toBeVolume)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetInductiveValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)
	}

	switch sl.actionObserved.IsActive() {
	case true:
		setVariable(actions.OriginalGullyVolume, actions.ActionedGullyVolume)
	case false:
		setVariable(actions.ActionedGullyVolume, actions.OriginalGullyVolume)
	}
}

func (sl *SedimentProduction) handleInitialisingGullyRestorationAction() {
	setVariable := func(asIsName action.ModelVariableName, toBeName action.ModelVariableName) {
		asIsVolume := sl.actionObserved.ModelVariableValue(asIsName)
		toBeVolume := sl.actionObserved.ModelVariableValue(toBeName)

		asIsSedimentContribution := sl.gullySedimentContribution.SedimentFromVolume(asIsVolume)
		toBeSedimentContribution := sl.gullySedimentContribution.SedimentFromVolume(toBeVolume)

		currentValue := sl.BaseInductiveDecisionVariable.Value()
		sl.BaseInductiveDecisionVariable.SetValue(currentValue - asIsSedimentContribution + toBeSedimentContribution)
	}

	assert.That(sl.actionObserved.IsActive()).WithFailureMessage("initialising action should always be active").Holds()
	setVariable(actions.OriginalGullyVolume, actions.ActionedGullyVolume)
}
