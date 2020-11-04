// Copyright (c) 2019 Australian Rivers Institute.

package dissolvednitrogen

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
)

type HillSlopeRevegetationCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand

	undoneHillSlopeContribution float64
	doneHillSlopeContribution   float64
}

func (c *HillSlopeRevegetationCommand) ForVariable(variable variable.PlanningUnitDecisionVariable) *HillSlopeRevegetationCommand {
	c.WithTarget(variable)
	return c
}

func (c *HillSlopeRevegetationCommand) InPlanningUnit(planningUnit planningunit.Id) *HillSlopeRevegetationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.InPlanningUnit(planningUnit)
	return c
}

func (c *HillSlopeRevegetationCommand) WithFilteredNitrogenContribution(contribution float64) *HillSlopeRevegetationCommand {
	c.undoneHillSlopeContribution = c.hillSlopeNitrogenContribution()
	c.doneHillSlopeContribution = contribution
	return c
}

func (c *HillSlopeRevegetationCommand) WithChange(changeValue float64) *HillSlopeRevegetationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	return c
}

func (c *HillSlopeRevegetationCommand) variable() *DissolvedNitrogenProduction {
	return c.Target().(*DissolvedNitrogenProduction)
}

func (c *HillSlopeRevegetationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	c.setHillSlopeNitrogenContribution(c.doneHillSlopeContribution)
	return command.Done
}

func (c *HillSlopeRevegetationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	c.setHillSlopeNitrogenContribution(c.undoneHillSlopeContribution)
	return command.UnDone
}

func (c *HillSlopeRevegetationCommand) setHillSlopeNitrogenContribution(nitrogenContribution float64) {
	c.variable().subCatchmentAttributes[c.PlanningUnit()] =
		c.variable().subCatchmentAttributes[c.PlanningUnit()].Replace(HillSlopeNitrogenContribution, nitrogenContribution)
}

func (c *HillSlopeRevegetationCommand) hillSlopeNitrogenContribution() float64 {
	planningUnitAttributes := c.variable().subCatchmentAttributes[c.PlanningUnit()]
	return planningUnitAttributes.Value(HillSlopeNitrogenContribution).(float64)
}

func (c *HillSlopeRevegetationCommand) DoneHillSlopeContribution() float64 {
	return c.doneHillSlopeContribution
}

func (c *HillSlopeRevegetationCommand) UndoneHillSlopeContribution() float64 {
	return c.undoneHillSlopeContribution
}
