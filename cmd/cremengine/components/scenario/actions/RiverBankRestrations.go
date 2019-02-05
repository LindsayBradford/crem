// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"strconv"

	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

type RiverBankRestorations struct {
	planningUnitTable *tables.CsvTable
	parameters        parameters.Parameters

	actionMap map[planningUnitId]*RiverBankRestoration
}

func (r *RiverBankRestorations) Initialise(planningUnitTable *tables.CsvTable, parameters parameters.Parameters) *RiverBankRestorations {
	r.planningUnitTable = planningUnitTable
	r.parameters = parameters
	r.createManagementActions()

	return r
}

func (r *RiverBankRestorations) ManagementActions() map[planningUnitId]*RiverBankRestoration {
	return r.actionMap
}

func (r *RiverBankRestorations) createManagementActions() {
	_, rowCount := r.planningUnitTable.Size()
	r.actionMap = make(map[planningUnitId]*RiverBankRestoration, rowCount)

	for row := uint(0); row < rowCount; row++ {
		r.createManagementAction(row)
	}
}

func (r *RiverBankRestorations) createManagementAction(rowNumber uint) {
	planningUnit := r.planningUnitTable.CellInt64(planningUnitIndex, rowNumber)
	mapKey := planningUnitId(planningUnit)

	planningUnitAsString := strconv.FormatInt(planningUnit, 10)

	originalBufferVegetation := r.originalBufferVegetation(rowNumber)
	changeInBufferVegetation := r.calculateChangeInBufferVegetation(rowNumber)

	costInDollars := r.calculateImplementationCost(rowNumber)

	r.actionMap[mapKey] =
		new(RiverBankRestoration).
			WithPlanningUnit(planningUnitAsString).
			WithRiverBankRestorationType().
			WithOriginalBufferVegetation(originalBufferVegetation).
			WithChangeInBufferVegetation(changeInBufferVegetation).
			WithImplementationCost(costInDollars)
}

func (r *RiverBankRestorations) originalBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.planningUnitTable.CellFloat64(riparianVegetationIndex, rowNumber)
	return proportionOfRiparianVegetation
}

func (r *RiverBankRestorations) calculateChangeInBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.originalBufferVegetation(rowNumber)
	changeInRiparianVegetation := 1 - proportionOfRiparianVegetation // TODO: Assumes full riparian revegation - revisit later.
	return changeInRiparianVegetation
}

func (r *RiverBankRestorations) calculateImplementationCost(rowNumber uint) float64 {
	riparianRevegetationCostPerKilometer := r.parameters.GetFloat64(parameters.RiparianRevegetationCostPerKilometer)
	riverLength := r.planningUnitTable.CellFloat64(riverLengthIndex, rowNumber)
	vegetationChange := r.calculateChangeInBufferVegetation(rowNumber)

	vegetationChangeLengthInKilometres := vegetationChange * riverLength

	implementationCost := vegetationChangeLengthInKilometres * riparianRevegetationCostPerKilometer

	return implementationCost
}