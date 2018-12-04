// Copyright (c) 2018 Australian Rivers Institute.

package model

const ObjectiveValue = "ObjectiveValue"

type Model interface {
	TryRandomChange()

	AcceptChange()
	RevertChange()

	DecisionVariable(name string) (DecisionVariable, error)
	Change(decisionVariable DecisionVariable) (float64, error)
}

var NullModel = new(nullModel)

type nullModel struct{}

func (nm *nullModel) Name() string     { return "NullModel" }
func (nm *nullModel) TryRandomChange() {}
func (nm *nullModel) AcceptChange()    {}
func (nm *nullModel) RevertChange()    {}
func (nm *nullModel) DecisionVariable(name string) (DecisionVariable, error) {
	return NullDecisionVariable, nil
}
func (nm *nullModel) Change(decisionVariable DecisionVariable) (float64, error) { return 0, nil }
func (nm *nullModel) SetDecisionVariable(name string, value float64) error      { return nil }
