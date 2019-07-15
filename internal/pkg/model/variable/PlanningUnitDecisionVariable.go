// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/LindsayBradford/crem/internal/pkg/model/planningunit"

type PlanningUnitDecisionVariable interface {
	ValuesPerPlanningUnit() map[planningunit.Id]float64
}
