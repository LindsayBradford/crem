// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
)

type HillSlopeRevegetationCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand

	undoneHillSlopeVegetationProportion float64
	doneHillSlopeVegetationProportion   float64

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

func (c *HillSlopeRevegetationCommand) WithSedimentContribution(contribution float64) *HillSlopeRevegetationCommand {
	c.undoneHillSlopeContribution = c.hillSlopeSedimentContribution()
	c.doneHillSlopeContribution = contribution
	return c
}

func (c *HillSlopeRevegetationCommand) WithChange(changeValue float64) *HillSlopeRevegetationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	return c
}

func (c *HillSlopeRevegetationCommand) variable() *SedimentProduction {
	return c.Target().(*SedimentProduction)
}

func (c *HillSlopeRevegetationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	c.setHillSlopeVegetation(c.doneHillSlopeVegetationProportion)
	c.setHillSlopeSedimentContribution(c.doneHillSlopeContribution)
	return command.Done
}

func (c *HillSlopeRevegetationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	c.setHillSlopeVegetation(c.undoneHillSlopeVegetationProportion)
	c.setHillSlopeSedimentContribution(c.undoneHillSlopeContribution)
	return command.UnDone
}

func (c *HillSlopeRevegetationCommand) setHillSlopeVegetation(proportion float64) {
	c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(HillSlopeVegetationProportion, proportion)
}

func (c *HillSlopeRevegetationCommand) setHillSlopeSedimentContribution(sedimentContribution float64) {
	c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(HillSlopeSedimentContribution, sedimentContribution)
}

func (c *HillSlopeRevegetationCommand) hillSlopeSedimentContribution() float64 {
	planningUnitAttributes := c.variable().planningUnitAttributes[c.PlanningUnit()]
	return planningUnitAttributes.Value(HillSlopeSedimentContribution).(float64)
}

func (c *HillSlopeRevegetationCommand) DoneHillSlopeContribution() float64 {
	return c.doneHillSlopeContribution
}

func (c *HillSlopeRevegetationCommand) UndoneHillSlopeContribution() float64 {
	return c.undoneHillSlopeContribution
}
