// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
)

const (
	_ = iota
)

type DecisionVariables struct {
}

type ContainedDecisionVariables struct {
	variable.InductiveDecisionVariables
}

func (c *ContainedDecisionVariables) Initialise() {
	c.InductiveDecisionVariables = variable.NewInductiveDecisionVariables()
}

func (c *ContainedDecisionVariables) DecisionVariables() *variable.DecisionVariables {
	inductiveVariables := c.InductiveDecisionVariables
	vanillaVariables := make(variable.DecisionVariables, 0)
	for _, inductiveVariable := range inductiveVariables {
		vanillaVariables[inductiveVariable.Name()] = inductiveVariable
	}
	return &vanillaVariables
}

func (c *ContainedDecisionVariables) DecisionVariable(name string) variable.DecisionVariable {
	return c.InductiveDecisionVariables.Variable(name)
}

func (c *ContainedDecisionVariables) DecisionVariableChange(variableName string) float64 {
	return c.InductiveDecisionVariables.DifferenceInValues(variableName)
}
