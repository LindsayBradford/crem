// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

import (
	"github.com/LindsayBradford/crem/pkg/command"
	"github.com/LindsayBradford/crem/pkg/math"
)

func NewUndoableDecisionVariable(name string) *UndoableDecisionVariable {
	variable := new(UndoableDecisionVariable)

	variable.SetName(name)
	variable.SetPrecision(defaultPrecision)
	variable.SetUnitOfMeasure(NotApplicable)

	variable.command = new(UndoableValueCommand).ForVariable(variable)

	return variable
}

type UndoableDecisionVariable struct {
	SimpleDecisionVariable

	ContainedDecisionVariableObservers
	command *UndoableValueCommand
}

func (v *UndoableDecisionVariable) InductiveValue() float64 {
	return v.command.Value()
}

func (v *UndoableDecisionVariable) SetInductiveChange(change float64) {
	v.command.WithChange(change)
}

func (v *UndoableDecisionVariable) DifferenceInValues() float64 {
	return v.command.Change()
}

func (v *UndoableDecisionVariable) AcceptInductiveValue() {
	v.command.Do()
}

func (v *UndoableDecisionVariable) RejectInductiveValue() {
	v.command.Undo()
}

type UndoableValueCommand struct {
	command.BaseCommand

	undoneValue float64
	doneValue   float64
}

func (c *UndoableValueCommand) ForVariable(variable *UndoableDecisionVariable) *UndoableValueCommand {
	c.WithTarget(variable)
	return c
}

func (c *UndoableValueCommand) WithChange(changeValue float64) *UndoableValueCommand {
	c.SetChange(changeValue)
	return c
}

func (c *UndoableValueCommand) SetChange(changeValue float64) {
	c.undoneValue = c.Variable().Value()
	roundedChangeValue := math.RoundFloat(changeValue, int(c.Variable().Precision()))
	c.doneValue = c.undoneValue + roundedChangeValue
}

func (c *UndoableValueCommand) Do() command.CommandStatus {
	if c.BaseCommand.Do() == command.NoChange {
		return command.NoChange
	}
	c.DoUnguarded()
	return command.Done
}

func (c *UndoableValueCommand) DoUnguarded() {
	c.Variable().SetValue(c.doneValue)
}

func (c *UndoableValueCommand) Undo() command.CommandStatus {
	if c.BaseCommand.Undo() == command.NoChange {
		return command.NoChange
	}
	c.UndoUnguarded()
	return command.UnDone
}

func (c *UndoableValueCommand) UndoUnguarded() {
	c.Variable().SetValue(c.undoneValue)
}

func (c *UndoableValueCommand) Variable() *UndoableDecisionVariable {
	return c.Target().(*UndoableDecisionVariable)
}

func (c *UndoableValueCommand) Value() float64 {
	return c.doneValue
}

func (c *UndoableValueCommand) Change() float64 {
	return c.doneValue - c.undoneValue
}
