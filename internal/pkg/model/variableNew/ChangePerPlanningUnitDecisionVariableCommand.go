// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/pkg/command"
)

type ChangePerPlanningUnitDecisionVariableCommand struct {
	command.BaseCommand

	undoneValue  float64
	doneValue    float64
	planningUnit planningunit.Id
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) ForVariable(variable *PerPlanningUnitDecisionVariable) *ChangePerPlanningUnitDecisionVariableCommand {
	c.WithTarget(variable)
	return c
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) WithValue(newValue float64) *ChangePerPlanningUnitDecisionVariableCommand {
	c.undoneValue = c.variable().value
	c.doneValue = newValue
	return c
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) InPlanningUnit(planningUnit planningunit.Id) *ChangePerPlanningUnitDecisionVariableCommand {
	c.planningUnit = planningUnit
	return c
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) Do() {
	c.variable().SetPlanningUnitValue(c.planningUnit, c.doneValue)
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) Undo() {
	c.variable().SetPlanningUnitValue(c.planningUnit, c.undoneValue)
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) variable() *PerPlanningUnitDecisionVariable {
	return c.Target().(*PerPlanningUnitDecisionVariable)
}
