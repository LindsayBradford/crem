// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

import "github.com/LindsayBradford/crem/internal/pkg/model/planningunit"

type PlanningUnitValueMap map[planningunit.Id]float64

type PlanningUnitDecisionVariable interface {
	ValuesPerPlanningUnit() PlanningUnitValueMap
}

type ContainedValuesPerPlanningUnit struct {
	planningUnitValues PlanningUnitValueMap
}

func (c *ContainedValuesPerPlanningUnit) ValuesPerPlanningUnit() PlanningUnitValueMap {
	return c.planningUnitValues
}

func (c *ContainedValuesPerPlanningUnit) PlanningUnitValue(planningUnit planningunit.Id) float64 {
	return c.planningUnitValues[planningUnit]
}

func (c *ContainedValuesPerPlanningUnit) SetPlanningUnitValue(planningUnit planningunit.Id, value float64) {
	c.planningUnitValues[planningUnit] = value
}
