// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
)

const (
	_ = iota
)

type DecisionVariables struct {
	variable.InductiveDecisionVariables
}

func (dv *DecisionVariables) Initialise() *DecisionVariables {
	dv.InductiveDecisionVariables = variable.NewInductiveDecisionVariables()
	return dv
}

type ContainedDecisionVariables struct {
	decisionVariables DecisionVariables
}

func (c *ContainedDecisionVariables) DecisionVariables() *DecisionVariables {
	return &c.decisionVariables
}

func (c *ContainedDecisionVariables) DecisionVariable(name string) variable.DecisionVariable {
	return c.decisionVariables.Variable(name)
}

func (c *ContainedDecisionVariables) DecisionVariableChange(variableName string) float64 {
	return c.decisionVariables.DifferenceInValues(variableName)
}
