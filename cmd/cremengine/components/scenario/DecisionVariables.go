// Copyright (c) 2019 Australian Rivers Institute.

package scenario

import "github.com/LindsayBradford/crem/internal/pkg/model"

const (
	_                   = iota
	SedimentLoad string = "SedimentLoad"
)

type DecisionVariables struct {
	model.VolatileDecisionVariables
}

func (dv *DecisionVariables) Initialise() *DecisionVariables {
	dv.VolatileDecisionVariables = model.NewVolatileDecisionVariables()
	dv.buildDecisionVariables()
	return dv
}

func (dv *DecisionVariables) buildDecisionVariables() {
	dv.NewForName(SedimentLoad)
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
