// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

type WetlandsEstablishmentGroup struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	actionMap map[planningunit.Id]*WetlandsEstablishment
	Container
}

func (w *WetlandsEstablishmentGroup) WithPlanningUnitTable(planningUnitTable tables.CsvTable) *WetlandsEstablishmentGroup {
	w.planningUnitTable = planningUnitTable
	return w
}

func (w *WetlandsEstablishmentGroup) WithActionsTable(parentSoilsTable tables.CsvTable) *WetlandsEstablishmentGroup {
	w.Container.WithFilter(WetlandType).WithActionsTable(parentSoilsTable)
	return w
}

func (w *WetlandsEstablishmentGroup) WithParameters(parameters parameters.Parameters) *WetlandsEstablishmentGroup {
	w.parameters = parameters
	return w
}

func (w *WetlandsEstablishmentGroup) ManagementActions() []action.ManagementAction {
	w.createManagementActions()
	actions := make([]action.ManagementAction, 0)
	for _, value := range w.actionMap {
		actions = append(actions, value)
	}
	return actions
}

func (w *WetlandsEstablishmentGroup) createManagementActions() {
	// BUG: One per planning unit is wrong.  The filtered container contains a better count.

	_, rowCount := w.planningUnitTable.ColumnAndRowSize()
	w.actionMap = make(map[planningunit.Id]*WetlandsEstablishment, rowCount)

	for row := uint(0); row < rowCount; row++ {
		w.createManagementAction(row)
	}
}

func (w *WetlandsEstablishmentGroup) createManagementAction(rowNumber uint) {
	planningUnit := w.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	planningUnitAsId := planningunit.Float64ToId(planningUnit)

	if !w.mapsToPlanningUnit(planningUnitAsId) {
		return
	}

	opportunityCostInDollars := w.opportunityCost(planningUnitAsId)
	implementationCostInDollars := w.implementationCost(planningUnitAsId)

	dissolvedNitrogenRemovalEfficiency := w.dissolvedNitrogenRemovalEfficiency(planningUnitAsId)
	particulateNitrogenRemovalEfficiency := w.particulateNitrogenRemovalEfficiency(planningUnitAsId)
	sedimentRemovalEfficiency := w.sedimentNitrogenRemovalEfficiency(planningUnitAsId)

	w.actionMap[planningUnitAsId] =
		NewWetlandsEstablishment().
			WithPlanningUnit(planningUnitAsId).
			WithImplementationCost(implementationCostInDollars).
			WithOpportunityCost(opportunityCostInDollars).
			WithDissolvedNitrogenRemovalEfficiency(dissolvedNitrogenRemovalEfficiency).
			WithParticulateNitrogenRemovalEfficiency(particulateNitrogenRemovalEfficiency).
			WithSedimentRemovalEfficiency(sedimentRemovalEfficiency)
}
