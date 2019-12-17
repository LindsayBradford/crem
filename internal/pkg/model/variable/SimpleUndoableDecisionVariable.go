// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"github.com/LindsayBradford/crem/pkg/command"
	"github.com/LindsayBradford/crem/pkg/math"
)

var _ UndoableDecisionVariable = new(SimpleUndoableDecisionVariable)

func NewUndoableDecisionVariable(name string) *SimpleUndoableDecisionVariable {
	variable := new(SimpleUndoableDecisionVariable)

	variable.SetName(name)
	variable.SetPrecision(defaultPrecision)
	variable.SetUnitOfMeasure(NotApplicable)

	variable.command = new(UndoableValueCommand).ForVariable(variable)

	return variable
}

type SimpleUndoableDecisionVariable struct {
	SimpleDecisionVariable

	ContainedDecisionVariableObservers
	command *UndoableValueCommand
}

func (v *SimpleUndoableDecisionVariable) UndoableValue() float64 {
	return v.command.Value()
}

func (v *SimpleUndoableDecisionVariable) SetUndoableValue(value float64) {
	v.command.WithValue(value)
}

func (v *SimpleUndoableDecisionVariable) SetUndoableChange(change float64) {
	v.command.WithChange(change)
}

func (v *SimpleUndoableDecisionVariable) DifferenceInValues() float64 {
	return v.command.Change()
}

func (v *SimpleUndoableDecisionVariable) ApplyDoneValue() {
	v.command.Do()
}

func (v *SimpleUndoableDecisionVariable) ApplyUndoneValue() {
	v.command.Undo()
}

type UndoableValueCommand struct {
	command.BaseCommand

	undoneValue float64
	doneValue   float64
}

func (c *UndoableValueCommand) ForVariable(variable *SimpleUndoableDecisionVariable) *UndoableValueCommand {
	c.WithTarget(variable)
	return c
}

func (c *UndoableValueCommand) WithValue(value float64) *UndoableValueCommand {
	c.SetValue(value)
	return c
}

func (c *UndoableValueCommand) SetValue(value float64) {
	c.undoneValue = c.Variable().Value()
	roundedValue := math.RoundFloat(value, int(c.Variable().Precision()))
	c.doneValue = roundedValue
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

func (c *UndoableValueCommand) Variable() *SimpleUndoableDecisionVariable {
	return c.Target().(*SimpleUndoableDecisionVariable)
}

func (c *UndoableValueCommand) Value() float64 {
	return c.doneValue
}

func (c *UndoableValueCommand) Change() float64 {
	return c.doneValue - c.undoneValue
}
