// Copyright (c) 2019 Australian Rivers Institute.

package nitrogenproduction

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction2"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/command"
)

type RiverBankRestorationCommand struct {
	variable.ChangePerPlanningUnitDecisionVariableCommand

	doneRiparianVegetationProportion   float64
	undoneRiparianVegetationProportion float64

	undoneRiverbankContribution float64
	doneRiverbankContribution   float64
}

func (c *RiverBankRestorationCommand) ForVariable(variable variable.PlanningUnitDecisionVariable) *RiverBankRestorationCommand {
	c.WithTarget(variable)
	return c
}

func (c *RiverBankRestorationCommand) InPlanningUnit(planningUnit planningunit.Id) *RiverBankRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.InPlanningUnit(planningUnit)
	return c
}

//func (c *RiverBankRestorationCommand) WithVegetationBuffer(vegetationBuffer float64) *RiverBankRestorationCommand {
//	planningUnitAttributes := c.variable().planningUnitAttributes[c.PlanningUnit()]
//	c.undoneRiparianVegetationProportion = planningUnitAttributes.Value(RiverbankVegetationProportion).(float64)
//	c.doneRiparianVegetationProportion = vegetationBuffer
//	return c
//}

func (c *RiverBankRestorationCommand) WithChange(changeValue float64) *RiverBankRestorationCommand {
	c.ChangePerPlanningUnitDecisionVariableCommand.WithChange(changeValue)
	//
	//c.undoneRiverbankContribution = c.riverbankSedimentContribution()
	//c.doneRiverbankContribution = c.undoneRiverbankContribution + changeValue
	//
	return c
}

func (c *RiverBankRestorationCommand) variable() *sedimentproduction2.SedimentProduction2 {
	return c.Target().(*sedimentproduction2.SedimentProduction2)
}

func (c *RiverBankRestorationCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.DoUnguarded()
	//c.setRiparianVegetation(c.doneRiparianVegetationProportion)
	//c.setRiverbankSedimentContribution(c.doneRiverbankContribution)
	return command.Done
}

func (c *RiverBankRestorationCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.ChangePerPlanningUnitDecisionVariableCommand.UndoUnguarded()
	//c.setRiparianVegetation(c.undoneRiparianVegetationProportion)
	//c.setRiverbankSedimentContribution(c.undoneRiverbankContribution)
	return command.UnDone
}

//func (c *RiverBankRestorationCommand) setRiparianVegetation(proportion float64) {
//	c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(RiverbankVegetationProportion, proportion)
//}
//
//func (c *RiverBankRestorationCommand) setRiverbankSedimentContribution(sedimentContribution float64) {
//	c.variable().planningUnitAttributes[c.PlanningUnit()].Replace(RiverbankSedimentContribution, sedimentContribution)
//}
//
//func (c *RiverBankRestorationCommand) riverbankSedimentContribution() float64 {
//	planningUnitAttributes := c.variable().PlanningUnitAttributes()[c.PlanningUnit()]
//	return planningUnitAttributes.Value(RiverbankSedimentContribution).(float64)
//}
