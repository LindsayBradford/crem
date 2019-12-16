// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/LindsayBradford/crem/internal/pkg/model/variableNew"

const (
	_ = iota
)

type ContainedDecisionVariables struct {
	InductiveDecisionVariables
}

func (c *ContainedDecisionVariables) Initialise() {
	c.InductiveDecisionVariables = NewInductiveDecisionVariables()
}

func (c *ContainedDecisionVariables) DecisionVariables() *variableNew.DecisionVariableMap {
	inductiveVariables := c.InductiveDecisionVariables
	vanillaVariables := make(variableNew.DecisionVariableMap, 0)
	for _, inductiveVariable := range inductiveVariables {
		vanillaVariables[inductiveVariable.Name()] = inductiveVariable
	}
	return &vanillaVariables
}

func (c *ContainedDecisionVariables) DecisionVariable(name string) variableNew.DecisionVariable {
	return c.InductiveDecisionVariables.Variable(name)
}

func (c *ContainedDecisionVariables) OffersDecisionVariable(name string) bool {
	return c.InductiveDecisionVariables.Variable(name) != nil
}

func (c *ContainedDecisionVariables) DecisionVariableChange(variableName string) float64 {
	return c.InductiveDecisionVariables.DifferenceInValues(variableName)
}
