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

type EncodeableDecisionVariable struct {
	Name    string
	Value   float64
	Measure UnitOfMeasure `json:"UnitOfMeasure"`
}

func MakeEncodeable(variable DecisionVariable) EncodeableDecisionVariable {
	return EncodeableDecisionVariable{
		Name:    variable.Name(),
		Value:   math.RoundFloat(variable.Value(), int(variable.Precision())),
		Measure: variable.UnitOfMeasure(),
	}
}
