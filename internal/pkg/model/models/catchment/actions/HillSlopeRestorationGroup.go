// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

type HillSlopeRestorationGroup struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	actionMap map[planningunit.Id]*HillSlopeRestoration
}

func (r *HillSlopeRestorationGroup) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *HillSlopeRestorationGroup {
	r.planningUnitTable = planningUnitTable
	r.parameters = parameters
	r.createManagementActions()

	return r
}

func (r *HillSlopeRestorationGroup) ManagementActions() map[planningunit.Id]*HillSlopeRestoration {
	return r.actionMap
}

func (r *HillSlopeRestorationGroup) createManagementActions() {
	_, rowCount := r.planningUnitTable.ColumnAndRowSize()
	r.actionMap = make(map[planningunit.Id]*HillSlopeRestoration, rowCount)

	for row := uint(0); row < rowCount; row++ {
		r.createManagementAction(row)
	}
}

func (r *HillSlopeRestorationGroup) createManagementAction(rowNumber uint) {
	planningUnit := r.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	planningUnitAsId := planningunit.Float64ToId(planningUnit)

	originalHillSlopeVegetation := r.originalHillSlopeVegetation(rowNumber)

	costInDollars := r.calculateImplementationCost(rowNumber)

	vegetationTarget := r.parameters.GetFloat64(parameters.HillSlopeBevegetationProportionTarget)

	r.actionMap[planningUnitAsId] =
		NewHillSlopeRestoration().
			WithPlanningUnit(planningUnitAsId).
			WithOriginalHillSlopeVegetation(originalHillSlopeVegetation).
			WithActionedHillSlopeVegetation(vegetationTarget).
			WithImplementationCost(costInDollars)
}

func (r *HillSlopeRestorationGroup) originalHillSlopeVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.planningUnitTable.CellFloat64(proportionOfHillSlopeVegetationIndex, rowNumber)
	return proportionOfRiparianVegetation
}

func (r *HillSlopeRestorationGroup) calculateChangeInHillSlopeVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.originalHillSlopeVegetation(rowNumber)
	vegetationTarget := r.parameters.GetFloat64(parameters.HillSlopeBevegetationProportionTarget)
	changeInRiparianVegetation := vegetationTarget - proportionOfRiparianVegetation
	return changeInRiparianVegetation
}

func (r *HillSlopeRestorationGroup) calculateImplementationCost(rowNumber uint) float64 {
	implementationCostPerKmSquared := r.parameters.GetFloat64(parameters.HillSlopeRestorationCostPerKilometerSquared)
	hillSlopeAreaInMetresSquared := r.planningUnitTable.CellFloat64(hillSlopeAreaIndex, rowNumber)
	hillSlopeAreaInKilometresSquared := hillSlopeAreaInMetresSquared / 1000

	vegetationChange := r.calculateChangeInHillSlopeVegetation(rowNumber)

	vegetationChangeInKilometresSquared := vegetationChange * hillSlopeAreaInKilometresSquared

	implementationCost := vegetationChangeInKilometresSquared * implementationCostPerKmSquared

	return implementationCost
}
