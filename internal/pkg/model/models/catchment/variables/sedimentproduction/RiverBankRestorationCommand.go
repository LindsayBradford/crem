// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/LindsayBradford/crem/pkg/command"
)

type RiverBankRestorationCommand struct {
	variableNew.ChangePerPlanningUnitDecisionVariableCommand

	doneRiparianVegetationProportion   float64
	undoneRiparianVegetationProportion float64
}

func (c *RiverBankRestorationCommand) ForVariable(variable variableNew.PlanningUnitDecisionVariable) *RiverBankRestorationCommand {
	c.WithTarget(variable)
	return c
}

func (c *RiverBankRestorationCommand) InPlanningUnit(planningUnit planningunit.Id) *RiverBankRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.InPlanningUnit(planningUnit)
	return c
}

func (c *RiverBankRestorationCommand) WithVegetationBuffer(vegetationBuffer float64) *RiverBankRestorationCommand {
	c.undoneRiparianVegetationProportion = c.variable().riparianVegetationProportionPerPlanningUnit[c.PlanningUnit()]
	c.doneRiparianVegetationProportion = vegetationBuffer
	return c
}

func (c *RiverBankRestorationCommand) WithChange(changeValue float64) *RiverBankRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	return c
}

func (c *RiverBankRestorationCommand) variable() *SedimentProduction {
	return c.Target().(*SedimentProduction)
}

func (c *RiverBankRestorationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	c.setRiparianVegetation(c.doneRiparianVegetationProportion)
	return command.Done
}

func (c *RiverBankRestorationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	c.setRiparianVegetation(c.undoneRiparianVegetationProportion)
	return command.UnDone
}

func (c *RiverBankRestorationCommand) setRiparianVegetation(proportion float64) {
	c.variable().riparianVegetationProportionPerPlanningUnit[c.PlanningUnit()] = proportion
}
