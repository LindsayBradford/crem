// Copyright (c) 2019 Australian Rivers Institute.

package totalnitrogen

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	catchmentActions "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	catchmentParameters "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/dissolvednitrogen"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/particulatenitrogen"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/math"
)

const (
	VariableName = "TotalNitrogen"

	planningUnitIndex           = 0
	proportionOfVegetationIndex = 8

	ProportionOfRiparianVegetation             = "ProportionOfRiparianVegetation"
	RiparianDissolvedNitrogenRemovalEfficiency = "RiparianDissolvedNitrogenRemovalEfficiency"
	WetlandsDissolvedNitrogenRemovalEfficiency = "WetlandsDissolvedNitrogenRemovalEfficiency"
	RiparianNitrogenContribution               = "RiparianNitrogenContribution"
	GullyNitrogenContribution                  = "GullyNitrogenContribution"
	HillSlopeNitrogenContribution              = "HillSlopeNitrogenContribution"
)

var _ variable.UndoableDecisionVariable = new(TotalNitrogenProduction)

type TotalNitrogenProduction struct {
	variable.PerPlanningUnitDecisionVariable
	variable.Bounds

	catchmentActions.Container

	command variable.ChangeCommand

	actionObserved action.ManagementAction

	numberOfSubCatchments uint

	particulateNitrogen *particulatenitrogen.ParticulateNitrogenProduction
	dissolvedNitrogen   *dissolvednitrogen.DissolvedNitrogenProduction

	LastUpdated string
}

func (tn *TotalNitrogenProduction) WithBaseNitrogenVariables(particulate *particulatenitrogen.ParticulateNitrogenProduction, dissolved *dissolvednitrogen.DissolvedNitrogenProduction) *TotalNitrogenProduction {
	tn.particulateNitrogen = particulate
	tn.dissolvedNitrogen = dissolved
	return tn
}

func (tn *TotalNitrogenProduction) Initialise(subCatchmentsTable tables.CsvTable, actionsTable tables.CsvTable, parameters catchmentParameters.Parameters) *TotalNitrogenProduction {
	tn.PerPlanningUnitDecisionVariable.Initialise()
	tn.Container.WithActionsTable(actionsTable)

	tn.SetName(VariableName)
	tn.SetUnitOfMeasure(variable.TonnesPerYear)
	tn.SetPrecision(3)

	tn.deriveInitialState(subCatchmentsTable)

	tn.command = new(variable.NullChangeCommand)

	return tn
}

func (tn *TotalNitrogenProduction) deriveInitialState(subCatchmentsTable tables.CsvTable) {
	tn.deriveNumberOfSubCatchments(subCatchmentsTable)
	tn.deriveInitialNitrogen(subCatchmentsTable)
}

func (tn *TotalNitrogenProduction) deriveNumberOfSubCatchments(subCatchmentsTable tables.CsvTable) {
	_, rowCount := subCatchmentsTable.ColumnAndRowSize()
	tn.numberOfSubCatchments = rowCount
}

func (tn *TotalNitrogenProduction) deriveInitialNitrogen(subCatchmentsTable tables.CsvTable) {
	for row := uint(0); row < tn.numberOfSubCatchments; row++ {
		subCatchmentFloat64 := subCatchmentsTable.CellFloat64(planningUnitIndex, row)
		subCatchment := Float64ToSubCatchmentId(subCatchmentFloat64)
		tn.calculateTotalNitrogenForPlanningUnit(subCatchment)
	}
}

func Float64ToSubCatchmentId(value float64) planningunit.Id {
	return planningunit.Id(value)
}

func (tn *TotalNitrogenProduction) calculateTotalNitrogenForPlanningUnit(pu planningunit.Id) {
	particulateValue := math.RoundFloat(tn.particulateNitrogen.PlanningUnitValue(pu), int(tn.Precision()))
	dissolvedValue := math.RoundFloat(tn.dissolvedNitrogen.PlanningUnitValue(pu), int(tn.Precision()))
	roundedTotal := math.RoundFloat(particulateValue+dissolvedValue, int(tn.Precision()))
	tn.SetPlanningUnitValue(pu, roundedTotal)
}

func (tn *TotalNitrogenProduction) WithName(variableName string) *TotalNitrogenProduction {
	tn.SetName(variableName)
	return tn
}

func (tn *TotalNitrogenProduction) WithStartingValue(value float64) *TotalNitrogenProduction {
	tn.SetPlanningUnitValue(0, value)
	return tn
}

func (tn *TotalNitrogenProduction) WithObservers(observers ...variable.Observer) *TotalNitrogenProduction {
	tn.Subscribe(observers...)
	return tn
}

func (tn *TotalNitrogenProduction) ObserveAction(action action.ManagementAction) {
	tn.observeAction(action)
}

func (tn *TotalNitrogenProduction) ObserveActionInitialising(action action.ManagementAction) {
	tn.observeAction(action)
	tn.command.Do()
}

func (tn *TotalNitrogenProduction) observeAction(action action.ManagementAction) {
	tn.actionObserved = action

	// why is this only intermittently working?
	particulateChange := math.RoundFloat(tn.particulateNitrogen.DifferenceInValues(), int(tn.Precision()))
	dissolvedChange := math.RoundFloat(tn.dissolvedNitrogen.DifferenceInValues(), int(tn.Precision()))

	roundedChange := math.RoundFloat(particulateChange+dissolvedChange, int(tn.Precision()))

	tn.command = new(variable.ChangePerPlanningUnitDecisionVariableCommand).
		ForVariable(tn).
		InPlanningUnit(tn.actionObserved.PlanningUnit()).
		WithChange(roundedChange)
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
func (tn *TotalNitrogenProduction) NotifyObservers() {
	for _, observer := range tn.Observers() {
		observer.ObserveDecisionVariable(tn)
	}
}

func (tn *TotalNitrogenProduction) UndoableValue() float64 {
	return tn.Value() + tn.command.Value()
}

func (tn *TotalNitrogenProduction) SetUndoableValue(value float64) {
	tn.command.SetChange(value)
}

func (tn *TotalNitrogenProduction) DifferenceInValues() float64 {
	return tn.command.Change()
}

func (tn *TotalNitrogenProduction) ApplyDoneValue() {
	tn.command.Do()
}

func (tn *TotalNitrogenProduction) ApplyUndoneValue() {
	tn.command.Undo()
}
