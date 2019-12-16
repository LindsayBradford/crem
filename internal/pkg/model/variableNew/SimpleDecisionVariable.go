// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

const defaultPrecision = 3

func NewSimpleDecisionVariable(name string) SimpleDecisionVariable {
	variable := SimpleDecisionVariable{name: name, value: 0}
	variable.SetPrecision(defaultPrecision)
	variable.SetUnitOfMeasure(NotApplicable)
	return variable
}

type SimpleDecisionVariable struct {
	name  string
	value float64
	ContainedUnitOfMeasure
	ContainedPrecision
}

func (v *SimpleDecisionVariable) Name() string           { return v.name }
func (v *SimpleDecisionVariable) SetName(name string)    { v.name = name }
func (v *SimpleDecisionVariable) Value() float64         { return v.value }
func (v *SimpleDecisionVariable) SetValue(value float64) { v.value = value }
