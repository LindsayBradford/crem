// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
)

var _ variable.ChangeCommand = new(GullyRestorationCommand)

type GullyRestorationCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand

	undoneGullyContribution float64
	doneGullyContribution   float64
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

	c.undoneGullyContribution = c.gullySedimentContribution()
	c.doneGullyContribution = c.undoneGullyContribution + changeValue

	return c
}

func (c *GullyRestorationCommand) variable() *SedimentProduction {
	return c.Target().(*SedimentProduction)
}

func (c *GullyRestorationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	c.setGullySedimentContribution(c.doneGullyContribution)
	return command.Done
}

func (c *GullyRestorationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	c.setGullySedimentContribution(c.undoneGullyContribution)
	return command.UnDone
}

func (c *GullyRestorationCommand) setGullySedimentContribution(sedimentContribution float64) {
	c.variable().planningUnitAttributes[c.PlanningUnit()] =
		c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(GullySedimentContribution, sedimentContribution)
}

func (c *GullyRestorationCommand) gullySedimentContribution() float64 {
	planningUnitAttributes := c.variable().planningUnitAttributes[c.PlanningUnit()]
	return planningUnitAttributes.Value(GullySedimentContribution).(float64)
}

func (c *GullyRestorationCommand) DoneGullyContribution() float64 {
	return c.doneGullyContribution
}

func (c *GullyRestorationCommand) UndoneGullyContribution() float64 {
	return c.undoneGullyContribution
}
