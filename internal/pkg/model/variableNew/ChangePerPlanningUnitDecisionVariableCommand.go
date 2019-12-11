// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/pkg/command"
	"github.com/LindsayBradford/crem/pkg/math"
)

type ChangeCommand interface {
	command.Command
	Value() float64
	SetChange(change float64)
	Change() float64
}

type ChangePerPlanningUnitDecisionVariableCommand struct {
	command.BaseCommand

	undoneValue  float64
	doneValue    float64
	planningUnit planningunit.Id
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) ForVariable(variable PlanningUnitDecisionVariable) *ChangePerPlanningUnitDecisionVariableCommand {
	c.WithTarget(variable)
	return c
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) WithChange(changeValue float64) *ChangePerPlanningUnitDecisionVariableCommand {
	c.SetChange(changeValue)
	return c
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) SetChange(changeValue float64) {
	c.undoneValue = c.Variable().ValuesPerPlanningUnit()[c.planningUnit]
	roundedChangeValue := math.RoundFloat(changeValue, int(c.Variable().Precision()))
	c.doneValue = c.undoneValue + roundedChangeValue
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) InPlanningUnit(planningUnit planningunit.Id) *ChangePerPlanningUnitDecisionVariableCommand {
	c.planningUnit = planningUnit
	return c
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) PlanningUnit() planningunit.Id {
	return c.planningUnit
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.DoUnguarded()
	return command.Done
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) DoUnguarded() {
	c.Variable().SetPlanningUnitValue(c.planningUnit, c.doneValue)
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.UndoUnguarded()
	return command.UnDone
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) UndoUnguarded() {
	c.Variable().SetPlanningUnitValue(c.planningUnit, c.undoneValue)
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) Variable() PlanningUnitDecisionVariable {
	return c.Target().(PlanningUnitDecisionVariable)
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) Value() float64 {
	return c.doneValue
}

func (c *ChangePerPlanningUnitDecisionVariableCommand) Change() float64 {
	return c.doneValue - c.undoneValue
}
