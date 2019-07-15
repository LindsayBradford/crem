// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const defaultRevegetationProportion = float64(0.75)

type RiverBankRestorations struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	actionMap map[planningunit.Id]*RiverBankRestoration
}

func (r *RiverBankRestorations) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *RiverBankRestorations {
	r.planningUnitTable = planningUnitTable
	r.parameters = parameters
	r.createManagementActions()

	return r
}

func (r *RiverBankRestorations) ManagementActions() map[planningunit.Id]*RiverBankRestoration {
	return r.actionMap
}

func (r *RiverBankRestorations) createManagementActions() {
	_, rowCount := r.planningUnitTable.ColumnAndRowSize()
	r.actionMap = make(map[planningunit.Id]*RiverBankRestoration, rowCount)

	for row := uint(0); row < rowCount; row++ {
		r.createManagementAction(row)
	}
}

func (r *RiverBankRestorations) createManagementAction(rowNumber uint) {
	planningUnit := r.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	planningUnitAsId := Float64ToPlanningUnitId(planningUnit)

	originalBufferVegetation := r.originalBufferVegetation(rowNumber)

	costInDollars := r.calculateImplementationCost(rowNumber)

	r.actionMap[planningUnitAsId] =
		new(RiverBankRestoration).
			WithPlanningUnit(planningUnitAsId).
			WithRiverBankRestorationType().
			WithUnActionedBufferVegetation(originalBufferVegetation).
			WithActionedBufferVegetation(defaultRevegetationProportion).
			WithImplementationCost(costInDollars)
}

func (r *RiverBankRestorations) originalBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.planningUnitTable.CellFloat64(riparianVegetationIndex, rowNumber)
	return proportionOfRiparianVegetation
}

func (r *RiverBankRestorations) calculateChangeInBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.originalBufferVegetation(rowNumber)
	changeInRiparianVegetation := defaultRevegetationProportion - proportionOfRiparianVegetation // TODO: Assumes full riparian revegation - revisit later.
	return changeInRiparianVegetation
}

func (r *RiverBankRestorations) calculateImplementationCost(rowNumber uint) float64 {
	riparianRevegetationCostPerKilometer := r.parameters.GetFloat64(parameters.RiparianRevegetationCostPerKilometer)
	riverLengthInMetres := r.planningUnitTable.CellFloat64(riverLengthIndex, rowNumber)
	riverLengthInKilometres := riverLengthInMetres / 1000

	vegetationChange := r.calculateChangeInBufferVegetation(rowNumber)

	vegetationChangeLengthInKilometres := vegetationChange * riverLengthInKilometres

	implementationCost := vegetationChangeLengthInKilometres * riparianRevegetationCostPerKilometer

	return implementationCost
}
