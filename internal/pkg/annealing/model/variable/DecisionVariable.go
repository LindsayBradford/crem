// Copyright (c) 2019 Australian Rivers Institute.

package variable

type DecisionVariables map[string]DecisionVariable
type DecisionVariableImplementations map[string]DecisionVariableImpl

func NewDecisionVariableImplementations() DecisionVariableImplementations {
	return make(DecisionVariableImplementations, 1)
}

type DecisionVariable interface {
	Name() string
	Value() float64
}

func NewDecisionVariableImpl(name string) DecisionVariableImpl {
	return DecisionVariableImpl{name: name, value: 0}
}

type DecisionVariableImpl struct {
	name  string
	value float64
}

func (dvi DecisionVariableImpl) Name() string           { return dvi.name }
func (dvi DecisionVariableImpl) SetName(name string)    { dvi.name = name }
func (dvi DecisionVariableImpl) Value() float64         { return dvi.value }
func (dvi DecisionVariableImpl) SetValue(value float64) { dvi.value = value }
