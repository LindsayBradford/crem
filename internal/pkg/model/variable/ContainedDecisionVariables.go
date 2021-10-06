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

func (c *ContainedDecisionVariables) NameMappedVariables() *DecisionVariableMap {
	variableMap := make(DecisionVariableMap, 0)
	for _, variable := range c.UndoableDecisionVariables {
		variableMap[variable.Name()] = variable
	}
	return &variableMap
}

func (c *ContainedDecisionVariables) CreationOrderedVariables() UndoableDecisionVariables {
	return c.UndoableDecisionVariables
}

func (c *ContainedDecisionVariables) DecisionVariableNames() []string {
	names := make([]string, 0)
	for _, variable := range c.UndoableDecisionVariables {
		names = append(names, variable.Name())
	}
	return names
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
