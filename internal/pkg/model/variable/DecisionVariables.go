// Copyright (c) 2019 Australian Rivers Institute.

package variable

const (
	_ = iota
)

type ContainedDecisionVariables struct {
	InductiveDecisionVariables
}

func (c *ContainedDecisionVariables) Initialise() {
	c.InductiveDecisionVariables = NewInductiveDecisionVariables()
}

func (c *ContainedDecisionVariables) DecisionVariables() *DecisionVariableMap {
	inductiveVariables := c.InductiveDecisionVariables
	vanillaVariables := make(DecisionVariableMap, 0)
	for _, inductiveVariable := range inductiveVariables {
		vanillaVariables[inductiveVariable.Name()] = inductiveVariable
	}
	return &vanillaVariables
}

func (c *ContainedDecisionVariables) DecisionVariable(name string) DecisionVariable {
	return c.InductiveDecisionVariables.Variable(name)
}

func (c *ContainedDecisionVariables) OffersDecisionVariable(name string) bool {
	return c.InductiveDecisionVariables.Variable(name) != nil
}

func (c *ContainedDecisionVariables) DecisionVariableChange(variableName string) float64 {
	return c.InductiveDecisionVariables.DifferenceInValues(variableName)
}
