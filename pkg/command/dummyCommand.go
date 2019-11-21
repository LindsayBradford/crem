// Copyright (c) 2019 Australian Rivers Institute.

package command

type Counter struct {
	value int
}

const (
	increment = "increment"
	decrement = "decrement"
)

type dummyCommand struct {
	BaseCommand
}

func (dc *dummyCommand) WithTarget(target interface{}) *dummyCommand {
	dc.target = target
	return dc
}

func (dc *dummyCommand) WithAttribute(name string, value interface{}) *dummyCommand {
	dc.ContainedAttributes.AddAttribute(name, value)
	return dc
}

func (dc *dummyCommand) counter() *Counter {
	return dc.target.(*Counter)
}

func (dc *dummyCommand) incrementValue() int {
	if dc.HasAttribute(increment) {
		if increment, isInteger := dc.Attribute(increment).(int); isInteger {
			return increment
		}
	}
	return 0
}

func (dc *dummyCommand) decrementValue() int {
	if dc.HasAttribute(decrement) {
		if decrement, isInteger := dc.Attribute(decrement).(int); isInteger {
			return decrement
		}
	}
	return 0
}

func (dc *dummyCommand) Do() {
	dc.counter().value += dc.incrementValue()
	dc.counter().value -= dc.decrementValue()
}

func (dc *dummyCommand) Undo() {
	dc.counter().value -= dc.incrementValue()
	dc.counter().value += dc.decrementValue()
}
