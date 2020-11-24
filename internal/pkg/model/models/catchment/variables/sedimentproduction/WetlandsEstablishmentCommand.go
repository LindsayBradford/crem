// Copyright (c) 2019 Australian Rivers Institute.

package sedimentproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
)

type WetlandsEstablishmentCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand

	undoneRemovalEfficiency float64
	doneRemovalEfficiency   float64
}

func (c *WetlandsEstablishmentCommand) ForVariable(variable variable.PlanningUnitDecisionVariable) *WetlandsEstablishmentCommand {
	c.WithTarget(variable)
	return c
}

func (c *WetlandsEstablishmentCommand) InPlanningUnit(planningUnit planningunit.Id) *WetlandsEstablishmentCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.InPlanningUnit(planningUnit)
	return c
}

func (c *WetlandsEstablishmentCommand) WithRemovalEfficiency(efficiency float64) *WetlandsEstablishmentCommand {
	c.undoneRemovalEfficiency = c.removalEfficiency()
	c.doneRemovalEfficiency = efficiency
	return c
}

func (c *WetlandsEstablishmentCommand) WithChange(changeValue float64) *WetlandsEstablishmentCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	return c
}

func (c *WetlandsEstablishmentCommand) variable() *SedimentProduction {
	return c.Target().(*SedimentProduction)
}

func (c *WetlandsEstablishmentCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	c.setRemovalEfficiency(c.doneRemovalEfficiency)
	return command.Done
}

func (c *WetlandsEstablishmentCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	c.setRemovalEfficiency(c.undoneRemovalEfficiency)
	return command.UnDone
}

func (c *WetlandsEstablishmentCommand) setRemovalEfficiency(sedimentContribution float64) {
	c.variable().planningUnitAttributes[c.PlanningUnit()] =
		c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(WetlandRemovalEfficiency, sedimentContribution)
}

func (c *WetlandsEstablishmentCommand) removalEfficiency() float64 {
	planningUnitAttributes := c.variable().planningUnitAttributes[c.PlanningUnit()]
	return planningUnitAttributes.Value(WetlandRemovalEfficiency).(float64)
}

func (c *WetlandsEstablishmentCommand) DoneRemovalEfficiency() float64 {
	return c.doneRemovalEfficiency
}

func (c *WetlandsEstablishmentCommand) UndoneRemovalEfficiency() float64 {
	return c.undoneRemovalEfficiency
}
