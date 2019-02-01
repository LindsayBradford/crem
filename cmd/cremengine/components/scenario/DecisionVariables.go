// Copyright (c) 2019 Australian Rivers Institute.

package scenario

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

func (c *ContainedDecisionVariables) DecisionVariable(name string) model.DecisionVariable {
	return c.decisionVariables.Variable(name)
}

func (c *ContainedDecisionVariables) DecisionVariableChange(variableName string) float64 {
	decisionVariable := c.decisionVariables.Variable(variableName)
	difference := decisionVariable.TemporaryValue() - decisionVariable.Value()
	return difference
}
