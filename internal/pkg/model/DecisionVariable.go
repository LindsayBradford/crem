// Copyright (c) 2018 Australian Rivers Institute.

package model

type DecisionVariables map[string]DecisionVariable
type DecisionVariableImplementations map[string]DecisionVariableImpl

func NewDecisionVariableImplementations() DecisionVariableImplementations {
	return make(DecisionVariableImplementations, 1)
}

type DecisionVariable interface {
	Name() string
	Value() float64
}

const NullDecisionVariableName = "NullDecisionVariable"

var NullDecisionVariable = new(nullDecisionVariable)

type nullDecisionVariable struct{}

func (ndv *nullDecisionVariable) Name() string           { return NullDecisionVariableName }
func (ndv *nullDecisionVariable) SetName(name string)    {}
func (ndv *nullDecisionVariable) Value() float64         { return 0 }
func (ndv *nullDecisionVariable) SetValue(value float64) {}

const ObjectiveValue = "ObjectiveValue"

var ObjectiveValueDecisionVariable = &DecisionVariableImpl{
	name:  ObjectiveValue,
	value: 0,
}

type DecisionVariableImpl struct {
	name  string
	value float64
}

func (dvi DecisionVariableImpl) Name() string           { return dvi.name }
func (dvi DecisionVariableImpl) SetName(name string)    { dvi.name = name }
func (dvi DecisionVariableImpl) Value() float64         { return dvi.value }
func (dvi DecisionVariableImpl) SetValue(value float64) { dvi.value = value }
