// Copyright (c) 2019 Australian Rivers Institute.

package variables

import "github.com/LindsayBradford/crem/internal/pkg/annealing/model"

const (
	_ = iota
)

type DecisionVariables struct {
	model.VolatileDecisionVariables
}

func (dv *DecisionVariables) Initialise() *DecisionVariables {
	dv.VolatileDecisionVariables = model.NewVolatileDecisionVariables()
	return dv
}

type ContainedDecisionVariables struct {
	decisionVariables DecisionVariables
}

func (c *ContainedDecisionVariables) DecisionVariables() *DecisionVariables {
	return &c.decisionVariables
}

func (c *ContainedDecisionVariables) DecisionVariable(name string) model.DecisionVariable {
	return c.decisionVariables.Variable(name)
}

func (c *ContainedDecisionVariables) DecisionVariableChange(variableName string) float64 {
	decisionVariable := c.decisionVariables.Variable(variableName)
	return decisionVariable.ChangeInValue()
}
