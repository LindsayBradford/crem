// Copyright (c) 2018 Australian Rivers Institute.

package model

type Model interface {
	Name() string
	SetName(name string)
	TryRandomChange()

	AcceptChange()
	RevertChange()

	DecisionVariable(name string) (DecisionVariable, error)
	DecisionVariableChange(decisionVariableName string) (float64, error)

	Clone() Model
}

var NullModel = new(nullModel)

func NewNullModel() *nullModel {
	newModel := new(nullModel).WithName("NullModel")
	return newModel
}

type nullModel struct {
	name string
}

func (nm *nullModel) Name() string { return nm.name }
func (nm *nullModel) WithName(name string) *nullModel {
	nm.SetName(name)
	return nm
}
func (nm *nullModel) SetName(name string) { nm.name = name }
func (nm *nullModel) TryRandomChange()    {}
func (nm *nullModel) AcceptChange()       {}
func (nm *nullModel) RevertChange()       {}
func (nm *nullModel) DecisionVariable(name string) (DecisionVariable, error) {
	newVariable := DecisionVariableImpl{
		name:  name,
		value: 0,
	}
	return &newVariable, nil
}
func (nm *nullModel) DecisionVariableChange(decisionVariableName string) (float64, error) {
	return 0, nil
}
func (nm *nullModel) SetDecisionVariable(name string, value float64) error { return nil }
func (nm *nullModel) Clone() Model {
	clone := *nm
	return &clone
}
