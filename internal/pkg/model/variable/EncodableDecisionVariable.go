// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/LindsayBradford/crem/pkg/strings"
)

var currencyConverter = strings.NewConverter().Localised().WithFloatingPointPrecision(2).PaddingZeros()
var defaultConverter = strings.NewConverter().Localised().WithFloatingPointPrecision(3).PaddingZeros()

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
	PlanningUnit planningunit.Id
	Value        float64
}

type PlanningUnitValues []PlanningUnitValue

func (v PlanningUnitValues) Len() int {
	return len(v)
}

func (v PlanningUnitValues) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v PlanningUnitValues) Less(i, j int) bool {
	return v[i].PlanningUnit < v[j].PlanningUnit
}

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
		roundedValue := math.RoundFloat(planningUnitValue, int(variable.Precision()))
		if roundedValue == 0 {
			continue
		}
		newValue := PlanningUnitValue{
			PlanningUnit: planningUnitId,
			Value:        roundedValue,
		}
		values = append(values, newValue)
	}

	sort.Sort(values)

	return values
}

func (v *EncodeableDecisionVariable) MarshalJSON() ([]byte, error) {
	// TODO: Code-stink is high.  Scrub this down.
	buffer := bytes.NewBufferString("{")

	var key string
	var value string

	key = "Name"
	value = v.Name
	buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\",", key, value))

	key = "Measure"
	value = v.Measure.String()
	buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\",", key, value))

	key = "Value"

	switch v.Measure {
	case Dollars:
		value = currencyConverter.Convert(v.Value)
	default:
		value = defaultConverter.Convert(v.Value)
	}

	buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"", key, value))

	if len(v.ValuePerPlanningUnit) > 0 {
		buffer.WriteString(",")
		key = "ValuePerPlanningUnit"
		buffer.WriteString(fmt.Sprintf("\"%s\":[", key))

		length := len(v.ValuePerPlanningUnit)
		count := 0
		for _, planningUnitValue := range v.ValuePerPlanningUnit {

			key = planningUnitValue.PlanningUnit.String()

			var formattedPlanningUnitValue string
			switch v.Measure {
			case Dollars:
				formattedPlanningUnitValue = currencyConverter.Convert(planningUnitValue.Value)
			default:
				formattedPlanningUnitValue = defaultConverter.Convert(planningUnitValue.Value)
			}

			buffer.WriteString(fmt.Sprintf("{\"PlanningUnit\":\"%s\", \"Value\":\"%s\"}", key, formattedPlanningUnitValue))

			count++
			if count < length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("]")
	}

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
