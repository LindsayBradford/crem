// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

import "C"
import "github.com/LindsayBradford/crem/internal/pkg/model/planningunit"

type PlanningUnitValueMap map[planningunit.Id]float64

type PlanningUnitDecisionVariable interface {
	DecisionVariable
	ValuesPerPlanningUnit() PlanningUnitValueMap
}

func NewPerPlanningUnitDecisionVariable() *PerPlanningUnitDecisionVariable {
	variable := new(PerPlanningUnitDecisionVariable)
	variable.planningUnitValues = make(PlanningUnitValueMap, 0)
	return variable
}

type PerPlanningUnitDecisionVariable struct {
	SimpleDecisionVariable
	planningUnitValues PlanningUnitValueMap
}

func (c *PerPlanningUnitDecisionVariable) SetPlanningUnitValue(planningUnit planningunit.Id, newValue float64) {
	oldValue := c.planningUnitValues[planningUnit]
	c.planningUnitValues[planningUnit] = newValue

	difference := newValue - oldValue
	c.value += difference
}

func (c *PerPlanningUnitDecisionVariable) PlanningUnitValue(planningUnit planningunit.Id) float64 {
	return c.planningUnitValues[planningUnit]
}

func (c *PerPlanningUnitDecisionVariable) ValuesPerPlanningUnit() PlanningUnitValueMap {
	return c.planningUnitValues
}
