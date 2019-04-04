// Copyright (c) 2019 Australian Rivers Institute.

package variable

type SimpleDecisionVariables map[string]*SimpleDecisionVariable

func NewSimpleDecisionVariables() SimpleDecisionVariables {
	return make(SimpleDecisionVariables, 1)
}

func NewSimpleDecisionVariable(name string) SimpleDecisionVariable {
	return SimpleDecisionVariable{name: name, value: 0}
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
