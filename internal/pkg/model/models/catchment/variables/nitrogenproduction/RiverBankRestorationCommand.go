// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
	"github.com/LindsayBradford/crem/pkg/math"
)

type RiverBankRestorationCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand

	doneRiparianVegetationProportion   float64
	undoneRiparianVegetationProportion float64

	undoneRiparianContribution float64
	doneRiparianContribution   float64
}

func (c *RiverBankRestorationCommand) ForVariable(variable variable.PlanningUnitDecisionVariable) *RiverBankRestorationCommand {
	c.WithTarget(variable)
	return c
}

func (c *RiverBankRestorationCommand) InPlanningUnit(planningUnit planningunit.Id) *RiverBankRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.InPlanningUnit(planningUnit)
	return c
}

func (c *RiverBankRestorationCommand) WithVegetationProportion(proportion float64) *RiverBankRestorationCommand {
	planningUnitAttributes := c.variable().planningUnitAttributes[c.PlanningUnit()]
	c.undoneRiparianVegetationProportion = planningUnitAttributes.Value(RiverbankVegetationProportion).(float64)
	change := proportion - c.undoneRiparianVegetationProportion
	roundedChange := math.RoundFloat(change, int(c.variable().Precision()))
	c.doneRiparianVegetationProportion = c.undoneRiparianVegetationProportion + roundedChange
	return c
}

func (c *RiverBankRestorationCommand) WithNitrogenContribution(contribution float64) *RiverBankRestorationCommand {
	c.undoneRiparianContribution = c.riparianNitrogenContribution()
	change := contribution - c.undoneRiparianContribution
	roundedChange := math.RoundFloat(change, int(c.variable().Precision()))
	c.doneRiparianContribution = c.undoneRiparianContribution + roundedChange
	return c
}

func (c *RiverBankRestorationCommand) variable() *ParticulateNitrogenProduction {
	return c.Target().(*ParticulateNitrogenProduction)
}

func (c *RiverBankRestorationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.setRiparianVegetationProportion(c.doneRiparianVegetationProportion)
	c.setRiparianContribution(c.doneRiparianContribution)
	return command.Done
}

func (c *RiverBankRestorationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.setRiparianVegetationProportion(c.undoneRiparianVegetationProportion)
	c.setRiparianContribution(c.undoneRiparianContribution)
	return command.UnDone
}

func (c *RiverBankRestorationCommand) setRiparianVegetationProportion(proportion float64) {
	c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(RiverbankVegetationProportion, proportion)
}

func (c *RiverBankRestorationCommand) riparianNitrogenContribution() float64 {
	planningUnitAttributes := c.variable().planningUnitAttributes[c.PlanningUnit()]
	return planningUnitAttributes.Value(RiparianNitrogenContribution).(float64)
}

func (c *RiverBankRestorationCommand) setRiparianContribution(contribution float64) {
	c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(RiparianNitrogenContribution, contribution)
}
