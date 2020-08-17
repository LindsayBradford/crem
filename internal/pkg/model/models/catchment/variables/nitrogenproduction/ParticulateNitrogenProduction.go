// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	catchmentActions "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/pkg/errors"
)

const VariableName = "ParticulateNitrogen"
const notImplementedValue float64 = 0

var _ variable.UndoableDecisionVariable = new(ParticulateNitrogenProduction)

type ParticulateNitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	sedimentProductionVariable *sedimentproduction.SedimentProduction
	catchmentActions.Container

	command variable.ChangeCommand

	actionObserved action.ManagementAction

	hillSlopeNitrogenContribution float64
	bankNitrogenContribution      float64
	gullyNitrogenContribution     float64
}

func (np *ParticulateNitrogenProduction) Initialise(actionsTable tables.CsvTable) *ParticulateNitrogenProduction {
	np.PerPlanningUnitDecisionVariable.Initialise()
	np.Container.WithActionsTable(actionsTable)

	np.SetName(VariableName)
	np.SetUnitOfMeasure(variable.TonnesPerYear)
	np.SetPrecision(3)

	np.command = new(variable.NullChangeCommand)

	np.deriveInitialState()

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

// TODO: Deprecated?
func (np *ParticulateNitrogenProduction) WithSedimentProductionVariable(variable *sedimentproduction.SedimentProduction) *ParticulateNitrogenProduction {
	np.sedimentProductionVariable = variable
	return np
}

func (np *ParticulateNitrogenProduction) WithObservers(observers ...variable.Observer) *ParticulateNitrogenProduction {
	np.Subscribe(observers...)
	return np
}

func (np *ParticulateNitrogenProduction) deriveInitialState() {
	// TODO:   drive off particulate nitrogen values in actions table instead.
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
	// TODO: Implement
}

func (np *ParticulateNitrogenProduction) handleRiverBankRestorationAction_deprecated() {
	//planningUnit := np.actionObserved.PlanningUnit()
	//
	//sedimentVariableCommand := np.sedimentProductionVariable.Command()
	//
	//var asIsSediment, toBeSediment float64
	//if riverCommand, isRiverCommand := sedimentVariableCommand.(*sedimentproduction.RiverBankRestorationCommand); isRiverCommand {
	//	asIsSediment = riverCommand.UndoneRiverbankContribution()
	//	toBeSediment = riverCommand.DoneRiverbankContribution()
	//}
	//
	//var asIsCarbon, toBeCarbon, asIsNitrogen, toBeNitrogen float64
	//isActive := np.actionObserved.IsActive()
	//switch isActive {
	//case true:
	//	asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
	//	toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
	//
	//	asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
	//	toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	//case false:
	//	asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
	//	toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
	//
	//	asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
	//	toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	//}
	//
	//asIsVariables := particulateNitrogenVariables{
	//	sediment:      asIsSediment,
	//	totalCarbon:   asIsCarbon,
	//	totalNitrogen: asIsNitrogen,
	//}
	//
	//asIsParticulateNitrogen := calculateParticulateNitrogen(asIsVariables)
	//
	//toBeVariables := particulateNitrogenVariables{
	//	sediment:      toBeSediment,
	//	totalCarbon:   toBeCarbon,
	//	totalNitrogen: toBeNitrogen,
	//}
	//
	//toBeParticulateNitrogen := calculateParticulateNitrogen(toBeVariables)
	//
	//np.command = new(RiverBankRestorationCommand).
	//	ForVariable(np).
	//	InPlanningUnit(planningUnit).
	//	WithChange(toBeParticulateNitrogen - asIsParticulateNitrogen)
}

func (np *ParticulateNitrogenProduction) handleGullyRestorationActionMirroredSediment() {
	//actionPlanningUnit := np.actionObserved.PlanningUnit()
	//change := np.sedimentProductionVariable.DifferenceInValues()
	//
	//np.command = new(GullyRestorationCommand).
	//	ForVariable(np).
	//	InPlanningUnit(actionPlanningUnit).
	//	WithChange(change)
}

func (np *ParticulateNitrogenProduction) handleGullyRestorationAction() {
	// TODO: Implement
}

func (np *ParticulateNitrogenProduction) handleGullyRestorationAction_deprecated() {
	//planningUnit := np.actionObserved.PlanningUnit()
	//
	//sedimentVariableCommand := np.sedimentProductionVariable.Command()
	//
	//var asIsSediment, toBeSediment float64
	//if gullyCommand, isGullyCommand := sedimentVariableCommand.(*sedimentproduction.GullyRestorationCommand); isGullyCommand {
	//	asIsSediment = gullyCommand.UndoneGullyContribution()
	//	toBeSediment = gullyCommand.DoneGullyContribution()
	//}
	//
	//var asIsCarbon, toBeCarbon, asIsNitrogen, toBeNitrogen float64
	//isActive := np.actionObserved.IsActive()
	//switch isActive {
	//case true:
	//	asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
	//	toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
	//
	//	asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
	//	toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	//case false:
	//	asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
	//	toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
	//
	//	asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
	//	toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	//}
	//
	//asIsVariables := particulateNitrogenVariables{
	//	sediment:      asIsSediment,
	//	totalCarbon:   asIsCarbon,
	//	totalNitrogen: asIsNitrogen,
	//}
	//
	//asIsParticulateNitrogen := calculateParticulateNitrogen(asIsVariables)
	//
	//toBeVariables := particulateNitrogenVariables{
	//	sediment:      toBeSediment,
	//	totalCarbon:   toBeCarbon,
	//	totalNitrogen: toBeNitrogen,
	//}
	//
	//toBeParticulateNitrogen := calculateParticulateNitrogen(toBeVariables)
	//
	//np.command = new(GullyRestorationCommand).
	//	ForVariable(np).
	//	InPlanningUnit(planningUnit).
	//	WithChange(toBeParticulateNitrogen - asIsParticulateNitrogen)
}

func (np *ParticulateNitrogenProduction) handleHillSlopeRestorationAction() {
	// TODO: Implement
}

func (np *ParticulateNitrogenProduction) handleHillSlopeRestorationAction_deprecated() {
	//planningUnit := np.actionObserved.PlanningUnit()
	//
	//sedimentVariableCommand := np.sedimentProductionVariable.Command()
	//
	//var asIsSediment, toBeSediment float64
	//if hillSlopeCommand, isHillSlopeCommand := sedimentVariableCommand.(*sedimentproduction.HillSlopeRevegetationCommand); isHillSlopeCommand {
	//	asIsSediment = hillSlopeCommand.UndoneHillSlopeContribution()
	//	toBeSediment = hillSlopeCommand.DoneHillSlopeContribution()
	//}
	//
	//var asIsCarbon, toBeCarbon, asIsNitrogen, toBeNitrogen float64
	//isActive := np.actionObserved.IsActive()
	//switch isActive {
	//case true:
	//	asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
	//	toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
	//
	//	asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
	//	toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	//case false:
	//	asIsCarbon = np.actionObserved.ModelVariableValue(catchmentActions.ActionedTotalCarbon)
	//	toBeCarbon = np.actionObserved.ModelVariableValue(catchmentActions.OriginalTotalCarbon)
	//
	//	asIsNitrogen = np.actionObserved.ModelVariableValue(catchmentActions.TotalNitrogen)
	//	toBeNitrogen = asIsNitrogen // TODO:  Check in with Jing's revisiting this being constant across actions
	//}
	//
	//asIsVariables := particulateNitrogenVariables{
	//	sediment:      asIsSediment,
	//	totalCarbon:   asIsCarbon,
	//	totalNitrogen: asIsNitrogen,
	//}
	//
	//asIsParticulateNitrogen := calculateParticulateNitrogen(asIsVariables)
	//
	//toBeVariables := particulateNitrogenVariables{
	//	sediment:      toBeSediment,
	//	totalCarbon:   toBeCarbon,
	//	totalNitrogen: toBeNitrogen,
	//}
	//
	//toBeParticulateNitrogen := calculateParticulateNitrogen(toBeVariables)
	//
	//np.command = new(HillSlopeRevegetationCommand).
	//	ForVariable(np).
	//	InPlanningUnit(planningUnit).
	//	WithChange(toBeParticulateNitrogen - asIsParticulateNitrogen)
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
