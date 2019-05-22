// Copyright (c) 2019 Australian Rivers Institute.

package variable

type PlanningUnitDecisionVariable interface {
	ValuesPerPlanningUnit() map[string]float64
}
