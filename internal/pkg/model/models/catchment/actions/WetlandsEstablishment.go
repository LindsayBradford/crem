// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const WetlandsEstablishmentType action.ManagementActionType = "WetlandsEstablishment"

func NewWetlandsEstablishment() *WetlandsEstablishment {
	return new(WetlandsEstablishment).WithType(WetlandsEstablishmentType)
}

type WetlandsEstablishment struct {
	action.SimpleManagementAction
}

func (w *WetlandsEstablishment) WithType(actionType action.ManagementActionType) *WetlandsEstablishment {
	w.SimpleManagementAction.WithType(actionType)
	return w
}

func (w *WetlandsEstablishment) WithPlanningUnit(planningUnit planningunit.Id) *WetlandsEstablishment {
	w.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return w
}

const WetlandsEstablishmentCost action.ModelVariableName = "WetlandsEstablishmentCost"

func (w *WetlandsEstablishment) WithImplementationCost(costInDollars float64) *WetlandsEstablishment {
	return w.WithVariable(WetlandsEstablishmentCost, costInDollars)
}

const WetlandsEstablishmentOpportunityCost action.ModelVariableName = "WetlandsEstablishmentOpportunityCost"

func (w *WetlandsEstablishment) WithOpportunityCost(costInDollars float64) *WetlandsEstablishment {
	return w.WithVariable(WetlandsEstablishmentOpportunityCost, costInDollars)
}

func (w *WetlandsEstablishment) WithDissolvedNitrogenRemovalEfficiency(removalEfficiency float64) *WetlandsEstablishment {
	return w.WithVariable(DissolvedNitrogenRemovalEfficiency, removalEfficiency)
}

func (w *WetlandsEstablishment) WithParticulateNitrogenRemovalEfficiency(removalEfficiency float64) *WetlandsEstablishment {
	return w.WithVariable(ParticulateNitrogenRemovalEfficiency, removalEfficiency)
}

func (w *WetlandsEstablishment) WithSedimentRemovalEfficiency(removalEfficiency float64) *WetlandsEstablishment {
	return w.WithVariable(SedimentRemovalEfficiency, removalEfficiency)
}

func (w *WetlandsEstablishment) WithVariable(variableName action.ModelVariableName, value float64) *WetlandsEstablishment {
	w.SimpleManagementAction.WithVariable(variableName, value)
	return w
}
