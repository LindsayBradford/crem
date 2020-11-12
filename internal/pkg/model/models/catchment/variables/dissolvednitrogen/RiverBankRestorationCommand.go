// Copyright (c) 2019 Australian Rivers Institute.

package dissolvednitrogen

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
)

type RiverBankRestorationCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand

	undoneRiparianVegetationProportion float64
	doneRiparianVegetationProportion   float64

	undoneRemovalEfficiency float64
	doneRemovalEfficiency   float64

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
	planningUnitAttributes := c.variable().subCatchmentAttributes[c.PlanningUnit()]
	c.undoneRiparianVegetationProportion = planningUnitAttributes.Value(ProportionOfRiparianVegetation).(float64)
	c.doneRiparianVegetationProportion = proportion
	return c
}

func (c *RiverBankRestorationCommand) WithRemovalEfficiency(efficiency float64) *RiverBankRestorationCommand {
	planningUnitAttributes := c.variable().subCatchmentAttributes[c.PlanningUnit()]
	c.undoneRemovalEfficiency = planningUnitAttributes.Value(RiparianDissolvedNitrogenRemovalEfficiency).(float64)
	c.doneRemovalEfficiency = efficiency
	return c
}

func (c *RiverBankRestorationCommand) WithNitrogenContribution(contribution float64) *RiverBankRestorationCommand {
	c.undoneRiparianContribution = c.riparianNitrogenContribution()
	c.doneRiparianContribution = contribution
	return c
}

func (c *RiverBankRestorationCommand) WithChange(changeValue float64) *RiverBankRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	return c
}

func (c *RiverBankRestorationCommand) variable() *DissolvedNitrogenProduction {
	return c.Target().(*DissolvedNitrogenProduction)
}

func (c *RiverBankRestorationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	c.setRiparianVegetationProportion(c.doneRiparianVegetationProportion)
	c.setRemovalEfficiency(c.doneRemovalEfficiency)
	c.setRiparianNitrogenContribution(c.doneRiparianContribution)
	return command.Done
}

func (c *RiverBankRestorationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	c.setRiparianVegetationProportion(c.undoneRiparianVegetationProportion)
	c.setRemovalEfficiency(c.undoneRemovalEfficiency)
	c.setRiparianNitrogenContribution(c.undoneRiparianContribution)
	return command.UnDone
}

func (c *RiverBankRestorationCommand) setRiparianVegetationProportion(proportion float64) {
	c.variable().subCatchmentAttributes[c.PlanningUnit()] =
		c.variable().subCatchmentAttributes[c.PlanningUnit()].Replace(ProportionOfRiparianVegetation, proportion)
}

func (c *RiverBankRestorationCommand) setRemovalEfficiency(proportion float64) {
	c.variable().subCatchmentAttributes[c.PlanningUnit()] =
		c.variable().subCatchmentAttributes[c.PlanningUnit()].Replace(RiparianDissolvedNitrogenRemovalEfficiency, proportion)
}

func (c *RiverBankRestorationCommand) riparianNitrogenContribution() float64 {
	planningUnitAttributes := c.variable().subCatchmentAttributes[c.PlanningUnit()]
	return planningUnitAttributes.Value(RiparianNitrogenContribution).(float64)
}

func (c *RiverBankRestorationCommand) setRiparianNitrogenContribution(contribution float64) {
	c.variable().subCatchmentAttributes[c.PlanningUnit()] =
		c.variable().subCatchmentAttributes[c.PlanningUnit()].Replace(RiparianNitrogenContribution, contribution)
}
