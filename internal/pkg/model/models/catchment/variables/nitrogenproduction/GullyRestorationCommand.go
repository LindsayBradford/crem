// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
)

type GullyRestorationCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand
}

func (c *GullyRestorationCommand) ForVariable(variable variable.PlanningUnitDecisionVariable) *GullyRestorationCommand {
	c.WithTarget(variable)
	return c
}

func (c *GullyRestorationCommand) InPlanningUnit(planningUnit planningunit.Id) *GullyRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.InPlanningUnit(planningUnit)
	return c
}

func (c *GullyRestorationCommand) WithChange(changeValue float64) *GullyRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	return c
}

func (c *GullyRestorationCommand) variable() *ParticulateNitrogenProduction {
	return c.Target().(*ParticulateNitrogenProduction)
}

func (c *GullyRestorationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	return command.Done
}

func (c *GullyRestorationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	return command.UnDone
}
