// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

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

func (c *HillSlopeRevegetationCommand) WithChange(changeValue float64) *HillSlopeRevegetationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	return c
}

func (c *HillSlopeRevegetationCommand) variable() *ParticulateNitrogenProduction {
	return c.Target().(*ParticulateNitrogenProduction)
}

func (c *HillSlopeRevegetationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	return command.Done
}

func (c *HillSlopeRevegetationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	return command.UnDone
}
