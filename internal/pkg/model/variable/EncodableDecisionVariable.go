// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/LindsayBradford/crem/pkg/strings"
	"sort"
	strings2 "strings"
)

var currencyConverter = strings.NewConverter().Localised().WithFloatingPointPrecision(2).PaddingZeros()
var defaultConverter = strings.NewConverter().Localised().WithFloatingPointPrecision(3).PaddingZeros()

type EncodeableDecisionVariables []EncodeableDecisionVariable

const (
	nameKey                 = "Name"
	measureKey              = "Measure"
	valueKey                = "Value"
	valuePerPlanningUnitKey = "ValuePerPlanningUnit"

	comma      = ","
	openBrace  = "{"
	closeBrace = "}"
)

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
	planningUnitValues := v.deriveFormattedPerPlanningUnitValues()

	perAttributeJson := new(strings.FluentBuilder).
		Add(openBrace).
		Add(formatKeyValuePair(nameKey, v.Name)).Add(comma).
		Add(formatKeyValuePair(measureKey, v.Measure.String())).Add(comma).
		Add(formatKeyValuePair(valueKey, v.formatMeasureValue(v.Value))).
		AddIf(v.hasValuesPerPlanningUnit(), comma, formatKeyArrayPair(valuePerPlanningUnitKey, planningUnitValues)).
		Add(closeBrace).
		String()

	return []byte(perAttributeJson), nil
}

func (v *EncodeableDecisionVariable) deriveFormattedPerPlanningUnitValues() []string {
	perPlanningUnitValues := make([]string, 0)
	for _, planningUnitValue := range v.ValuePerPlanningUnit {
		formattedValue := v.formatPlanningUnitValue(planningUnitValue)
		perPlanningUnitValues = append(perPlanningUnitValues, formattedValue)
	}
	return perPlanningUnitValues
}

func (v *EncodeableDecisionVariable) formatPlanningUnitValue(planningUnitValue PlanningUnitValue) string {
	key := planningUnitValue.PlanningUnit.String()
	formattedValue := v.formatMeasureValue(planningUnitValue.Value)
	return fmt.Sprintf("{\"PlanningUnit\":\"%s\", \"Value\":\"%s\"}", key, formattedValue)
}

func (v *EncodeableDecisionVariable) hasValuesPerPlanningUnit() bool {
	return len(v.ValuePerPlanningUnit) > 0
}

func (v *EncodeableDecisionVariable) formatMeasureValue(value float64) string {
	switch v.Measure {
	case Dollars:
		return currencyConverter.Convert(value)
	default:
		return defaultConverter.Convert(value)
	}
	assert.That(false).WithFailureMessage("Should not reach here").Holds()
	return ""
}

func formatKeyValuePair(key string, value string) string {
	return fmt.Sprintf("\"%s\":\"%s\"", key, value)
}

func formatKeyArrayPair(key string, values []string) string {
	commaSeparatedValues := strings2.Join(values[:], comma)
	return fmt.Sprintf("\"%s\": [%s]", key, commaSeparatedValues)
}
