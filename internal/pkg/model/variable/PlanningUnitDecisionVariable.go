// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/pkg/math"
)

type PlanningUnitValueMap map[planningunit.Id]float64

type PlanningUnitDecisionVariable interface {
	DecisionVariable
	ValuesPerPlanningUnit() PlanningUnitValueMap
	SetPlanningUnitValue(planningUnit planningunit.Id, newValue float64)
	PlanningUnitValue(planningUnit planningunit.Id) float64
}

func NewPerPlanningUnitDecisionVariable() *PerPlanningUnitDecisionVariable {
	return new(PerPlanningUnitDecisionVariable).Initialise()
}

type PerPlanningUnitDecisionVariable struct {
	SimpleDecisionVariable
	ContainedDecisionVariableObservers

	planningUnitValues PlanningUnitValueMap
}

func (v *PerPlanningUnitDecisionVariable) Initialise() *PerPlanningUnitDecisionVariable {
	v.planningUnitValues = make(PlanningUnitValueMap, 0)
	return v
}

func (v *PerPlanningUnitDecisionVariable) SetPlanningUnitValue(planningUnit planningunit.Id, newPlanningUnitValue float64) {
	oldPlanningUnitValue := v.planningUnitValues[planningUnit]
	v.planningUnitValues[planningUnit] = newPlanningUnitValue
	v.planningUnitValues[planningUnit] = math.RoundFloat(v.planningUnitValues[planningUnit], int(v.Precision()))

	difference := newPlanningUnitValue - oldPlanningUnitValue
	v.value += difference

	v.value = math.RoundFloat(v.value, int(v.Precision()))
}

func (v *PerPlanningUnitDecisionVariable) PlanningUnitValue(planningUnit planningunit.Id) float64 {
	return v.planningUnitValues[planningUnit]
}

func (v *PerPlanningUnitDecisionVariable) ValuesPerPlanningUnit() PlanningUnitValueMap {
	return v.planningUnitValues
}

func (v *PerPlanningUnitDecisionVariable) NotifyObservers() {
	for _, observer := range v.Observers() {
		observer.ObserveDecisionVariable(v)
	}
}
