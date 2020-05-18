// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction2

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

func (c *HillSlopeRevegetationCommand) WithVegetationBuffer(vegetationBuffer float64) *HillSlopeRevegetationCommand {
	planningUnitAttributes := c.variable().planningUnitAttributes[c.PlanningUnit()]
	c.undoneHillSlopeVegetationProportion = planningUnitAttributes.Value(HillSlopeVegetationProportion).(float64)
	c.doneHillSlopeVegetationProportion = vegetationBuffer
	return c
}

func (c *HillSlopeRevegetationCommand) WithChange(changeValue float64) *HillSlopeRevegetationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)

	c.undoneHillSlopeContribution = c.hillSlopeSedimentContribution()
	c.doneHillSlopeContribution = c.undoneHillSlopeContribution + changeValue

	return c
}

func (c *HillSlopeRevegetationCommand) variable() *SedimentProduction2 {
	return c.Target().(*SedimentProduction2)
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
