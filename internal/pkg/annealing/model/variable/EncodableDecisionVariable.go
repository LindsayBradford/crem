// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/LindsayBradford/crem/pkg/math"

type EncodeableDecisionVariables []EncodeableDecisionVariable

func (v EncodeableDecisionVariables) Len() int {
	return len(v)
}

func (v EncodeableDecisionVariables) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v EncodeableDecisionVariables) Less(i, j int) bool {
	return v[i].Name < v[j].Name
}

type PlanningUnitValue struct {
	PlanningUnit string
	Value        float64
}

type PlanningUnitValues []PlanningUnitValue

type EncodeableDecisionVariable struct {
	Name                 string
	Value                float64
	Measure              UnitOfMeasure      `json:"UnitOfMeasure"`
	ValuePerPlanningUnit PlanningUnitValues `json:",omitempty"`
}

func MakeEncodeable(variable DecisionVariable) EncodeableDecisionVariable {
	return EncodeableDecisionVariable{
		Name:                 variable.Name(),
		Value:                math.RoundFloat(variable.Value(), int(variable.Precision())),
		Measure:              variable.UnitOfMeasure(),
		ValuePerPlanningUnit: encodeValuesPerPlanningUnit(variable),
	}
}

func encodeValuesPerPlanningUnit(variable DecisionVariable) PlanningUnitValues {
	variablePerPlanningUnit, isVariablePerPlanningUnit := variable.(PlanningUnitDecisionVariable)
	if !isVariablePerPlanningUnit {
		return nil
	}

	rawValues := variablePerPlanningUnit.ValuesPerPlanningUnit()

	values := make(PlanningUnitValues, 0)
	for planningUnitId, planningUnitValue := range rawValues {
		newValue := PlanningUnitValue{
			PlanningUnit: planningUnitId,
			Value:        math.RoundFloat(planningUnitValue, int(variable.Precision())),
		}
		values = append(values, newValue)
	}

	return values
}
