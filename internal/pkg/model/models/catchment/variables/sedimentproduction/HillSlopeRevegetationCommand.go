// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/LindsayBradford/crem/pkg/command"
)

type HillSlopeRevegetationCommand struct {
	variableNew.ChangePerPlanningUnitDecisionVariableCommand

	doneHillSlopeVegetationProportion   float64
	undoneHillSlopeVegetationProportion float64
}

func (c *HillSlopeRevegetationCommand) ForVariable(variable variableNew.PlanningUnitDecisionVariable) *HillSlopeRevegetationCommand {
	c.WithTarget(variable)
	return c
}

func (c *HillSlopeRevegetationCommand) InPlanningUnit(planningUnit planningunit.Id) *HillSlopeRevegetationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.InPlanningUnit(planningUnit)
	return c
}

func (c *HillSlopeRevegetationCommand) WithVegetationBuffer(vegetationBuffer float64) *HillSlopeRevegetationCommand {
	c.undoneHillSlopeVegetationProportion = c.variable().hillSlopeVegetationProportionPerPlanningUnit[c.PlanningUnit()]
	c.doneHillSlopeVegetationProportion = vegetationBuffer
	return c
}

func (c *HillSlopeRevegetationCommand) WithChange(changeValue float64) *HillSlopeRevegetationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
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
	return command.Done
}

func (c *HillSlopeRevegetationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	c.setHillSlopeVegetation(c.undoneHillSlopeVegetationProportion)
	return command.UnDone
}

func (c *HillSlopeRevegetationCommand) setHillSlopeVegetation(proportion float64) {
	c.variable().hillSlopeVegetationProportionPerPlanningUnit[c.PlanningUnit()] = proportion
}
