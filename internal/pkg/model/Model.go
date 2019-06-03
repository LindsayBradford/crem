// Copyright (c) 2018 Australian Rivers Institute.

package model

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

const LogLevel logging.Level = "Model"

type Model interface {
	name.Nameable
	name.Identifiable
	DecisionVariableContainer

	Initialise()
	TearDown()

	TryRandomChange()
	AcceptChange()
	RevertChange()

	DeepClone() Model

	ActiveManagementActions() []action.ManagementAction
	PlanningUnits() solution.PlanningUnitIds
}

// ContainedLogger defines an interface embedding a Model
type Container interface {
	Model() Model
	SetModel(model Model)
}

// ContainedModel is a struct offering a default implementation of ContainedLogger
type ContainedModel struct {
	model Model
}

func (c *ContainedModel) Model() Model {
	return c.model
}

func (c *ContainedModel) SetModel(model Model) {
	c.model = model
}

type DecisionVariableContainer interface {
	DecisionVariables() *variable.DecisionVariableMap

	DecisionVariable(name string) variable.DecisionVariable
	DecisionVariableChange(decisionVariableName string) float64
}

var NullModel = new(nullModel)

func NewNullModel() *nullModel {
	newModel := new(nullModel).WithName("NullModel")
	return newModel
}

type nullModel struct {
	name.NameContainer
	name.IdentifiableContainer
}

func (nm *nullModel) WithName(name string) *nullModel {
	nm.SetName(name)
	return nm
}

func (nm *nullModel) Initialise()                                      {}
func (nm *nullModel) TearDown()                                        {}
func (nm *nullModel) TryRandomChange()                                 {}
func (nm *nullModel) AcceptChange()                                    {}
func (nm *nullModel) RevertChange()                                    {}
func (nm *nullModel) DecisionVariables() *variable.DecisionVariableMap { return nil }
func (nm *nullModel) DecisionVariable(name string) variable.DecisionVariable {
	newVariable := variable.NewSimpleDecisionVariable(name)
	return &newVariable
}
func (nm *nullModel) DecisionVariableChange(decisionVariableName string) float64 { return 0 }
func (nm *nullModel) SetDecisionVariable(name string, value float64)             {}
func (nm *nullModel) ActiveManagementActions() []action.ManagementAction         { return nil }
func (nm *nullModel) PlanningUnits() solution.PlanningUnitIds                    { return nil }
func (nm *nullModel) DeepClone() Model                                           { return nm }
