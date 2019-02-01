// Copyright (c) 2019 Australian Rivers Institute.

package scenario

import (
	"strconv"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

type RiverBankRestorations struct {
	planningUnitTable *tables.CsvTable
	parameters        Parameters

	actionMap map[planningUnitId]*RiverBankRestoration
}

func (rbrs *RiverBankRestorations) Initialise(planningUnitTable *tables.CsvTable, parameters Parameters) {
	rbrs.planningUnitTable = planningUnitTable
	rbrs.parameters = parameters
	rbrs.createManagementActions()
}

func (rbrs *RiverBankRestorations) createManagementActions() {
	_, rowCount := rbrs.planningUnitTable.Size()
	rbrs.actionMap = make(map[planningUnitId]*RiverBankRestoration, rowCount)

	for row := uint(0); row < rowCount; row++ {
		rbrs.createManagementAction(row)
	}
}

func (rbrs *RiverBankRestorations) createManagementAction(rowNumber uint) {
	planningUnit := rbrs.planningUnitTable.CellInt64(planningUnitIndex, rowNumber)
	mapKey := planningUnitId(planningUnit)

	planningUnitAsString := strconv.FormatInt(planningUnit, 10)

	changeInBufferVegetation := rbrs.calculateChangeInBufferVegetation(rowNumber)
	costInDollars := rbrs.calculateImplementationCost(rowNumber)

	rbrs.actionMap[mapKey] =
		new(RiverBankRestoration).
			WithPlanningUnit(planningUnitAsString).
			WithHardcodedType().
			WithChangeInBufferVegetation(changeInBufferVegetation).
			WithImplementationCost(costInDollars)
}

func (rbrs *RiverBankRestorations) calculateChangeInBufferVegetation(rowNumber uint) float64 {
	return 1
}

func (rbrs *RiverBankRestorations) calculateImplementationCost(rowNumber uint) float64 {
	riparianRevegetationCostPerKilometer := rbrs.parameters.GetFloat64(RiparianRevegetationCostPerKilometer)
	riverLength := rbrs.planningUnitTable.CellFloat64(riverLengthIndex, rowNumber)
	vegetationChange := rbrs.calculateChangeInBufferVegetation(rowNumber)

	vegetationChangeLengthInKilometres := vegetationChange * riverLength

	implementationCost := vegetationChangeLengthInKilometres * riparianRevegetationCostPerKilometer

	return implementationCost
}

type RiverBankRestoration struct {
	action.SimpleManagementAction
}

func (rbr *RiverBankRestoration) WithPlanningUnit(planningUnit string) *RiverBankRestoration {
	rbr.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return rbr
}

const RiverBankRestorationType action.ManagementActionType = "RiverBankRestoration"

func (rbr *RiverBankRestoration) WithHardcodedType() *RiverBankRestoration {
	rbr.SimpleManagementAction.WithType(RiverBankRestorationType)
	return rbr
}

const RiverBankRestorationCost action.ModelVariableName = "RiverBankRestorationCost"

func (rbr *RiverBankRestoration) WithImplementationCost(costInDollars float64) *RiverBankRestoration {
	rbr.SimpleManagementAction.WithVariable(RiverBankRestorationCost, costInDollars)
	return rbr
}

const ChangeInBufferVegetation action.ModelVariableName = "ChangeInBufferVegetation"

func (rbr *RiverBankRestoration) WithChangeInBufferVegetation(changeAsBufferProportion float64) *RiverBankRestoration {
	rbr.SimpleManagementAction.WithVariable(ChangeInBufferVegetation, changeAsBufferProportion)
	return rbr
}
