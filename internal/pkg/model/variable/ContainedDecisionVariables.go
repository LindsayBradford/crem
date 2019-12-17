// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package variable

const (
	_ = iota
)

type ContainedDecisionVariables struct {
	UndoableDecisionVariables
}

func (c *ContainedDecisionVariables) Initialise() {
	c.UndoableDecisionVariables = NewUndoableDecisionVariables()
}

func (c *ContainedDecisionVariables) DecisionVariables() *DecisionVariableMap {
	inductiveVariables := c.UndoableDecisionVariables
	vanillaVariables := make(DecisionVariableMap, 0)
	for _, inductiveVariable := range inductiveVariables {
		vanillaVariables[inductiveVariable.Name()] = inductiveVariable
	}
	return &vanillaVariables
}

func (c *ContainedDecisionVariables) DecisionVariable(name string) DecisionVariable {
	return c.UndoableDecisionVariables.Variable(name)
}

func (c *ContainedDecisionVariables) OffersDecisionVariable(name string) bool {
	return c.UndoableDecisionVariables.Variable(name) != nil
}

func (c *ContainedDecisionVariables) DecisionVariableChange(variableName string) float64 {
	return c.UndoableDecisionVariables.DifferenceInValues(variableName)
}
