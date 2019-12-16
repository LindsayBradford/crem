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

func (dvi *SimpleDecisionVariable) Name() string           { return dvi.name }
func (dvi *SimpleDecisionVariable) SetName(name string)    { dvi.name = name }
func (dvi *SimpleDecisionVariable) Value() float64         { return dvi.value }
func (dvi *SimpleDecisionVariable) SetValue(value float64) { dvi.value = value }
