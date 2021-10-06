// Copyright (c) 2018 Australian Rivers Institute.

package model

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

const LogLevel logging.Level = "Model"

type InitialisationType int

const (
	AsIs InitialisationType = iota
	Random
	Unchanged
)

type Model interface {
	name.Nameable
	name.Identifiable
	DecisionVariableContainer

	Initialise(initialisationType InitialisationType)
	Randomize()
	TearDown()

	DoRandomChange()
	UndoChange()

	TryRandomChange()
	ChangeIsValid() (bool, *errors.CompositeError)
	AcceptChange()
	RevertChange()

	IsEquivalentTo(Model) bool
	SynchroniseTo(Model)

	DeepClone() Model

	ManagementActions() []action.ManagementAction
	ActiveManagementActions() []action.ManagementAction
	SetManagementAction(index int, value bool)
	SetManagementActionUnobserved(index int, value bool)

	PlanningUnits() planningunit.Ids
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
	NameMappedVariables() *variable.DecisionVariableMap

	DecisionVariable(name string) variable.DecisionVariable
	OffersDecisionVariable(name string) bool
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

func (nm *nullModel) Initialise(initialisationType InitialisationType) {}
func (nm *nullModel) Randomize()                                       {}
func (nm *nullModel) TearDown()                                        {}

func (nm *nullModel) DoRandomChange() {}
func (nm *nullModel) UndoChange()     {}

func (nm *nullModel) TryRandomChange()                              {}
func (nm *nullModel) ChangeIsValid() (bool, *errors.CompositeError) { return true, nil }
func (nm *nullModel) AcceptChange()                                 {}
func (nm *nullModel) RevertChange()                                 {}

func (nm *nullModel) NameMappedVariables() *variable.DecisionVariableMap { return nil }
func (nm *nullModel) DecisionVariable(name string) variable.DecisionVariable {
	newVariable := variable.NewSimpleDecisionVariable(name)
	return newVariable
}

func (nm *nullModel) OffersDecisionVariable(name string) bool                    { return true }
func (nm *nullModel) DecisionVariableChange(decisionVariableName string) float64 { return 0 }
func (nm *nullModel) SetDecisionVariable(name string, value float64)             {}

func (nm *nullModel) ManagementActions() []action.ManagementAction        { return nil }
func (nm *nullModel) ActiveManagementActions() []action.ManagementAction  { return nil }
func (nm *nullModel) SetManagementAction(index int, value bool)           {}
func (nm *nullModel) SetManagementActionUnobserved(index int, value bool) {}

func (nm *nullModel) PlanningUnits() planningunit.Ids { return nil }

func (nm *nullModel) IsEquivalentTo(Model) bool { return false }

func (nm *nullModel) DeepClone() Model          { return nm }
func (nm *nullModel) SynchroniseTo(model Model) {}
