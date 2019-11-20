// Copyright (c) 2019 Australian Rivers Institute.

package command

import (
	"testing"

	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestBaseCommand_DoUndo_NoPanic(t *testing.T) {
	g := NewGomegaWithT(t)

	testCommand := new(BaseCommand).
		WithTarget(new(Counter))

	doRunner := func() {
		testCommand.Do()
	}

	g.Expect(doRunner).ToNot(Panic())

	undoRunner := func() {
		testCommand.Undo()
	}

	g.Expect(undoRunner).ToNot(Panic())
}

type Counter struct {
	value int
}

const (
	increment = "increment"
	decrement = "decrement"
)

type DummyCommand struct {
	BaseCommand
}

func (dc *DummyCommand) WithTarget(target interface{}) *DummyCommand {
	dc.target = target
	return dc
}

func (dc *DummyCommand) WithAttribute(name string, value interface{}) *DummyCommand {
	dc.ContainedAttributes.AddAttribute(name, value)
	return dc
}

func (dc *DummyCommand) counter() *Counter {
	return dc.target.(*Counter)
}

func (dc *DummyCommand) incrementValue() int {
	if dc.HasAttribute(increment) {
		if increment, isInteger := dc.Attribute(increment).(int); isInteger {
			return increment
		}
	}
	return 0
}

func (dc *DummyCommand) decrementValue() int {
	if dc.HasAttribute(decrement) {
		if decrement, isInteger := dc.Attribute(decrement).(int); isInteger {
			return decrement
		}
	}
	return 0
}

func (dc *DummyCommand) Do() {
	dc.counter().value += dc.incrementValue()
	dc.counter().value -= dc.decrementValue()
}

func (dc *DummyCommand) Undo() {
	dc.counter().value -= dc.incrementValue()
	dc.counter().value += dc.decrementValue()
}

func TestDummyCommand_DoUndo_CounterCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	counterUnderTest := new(Counter)
	const expectedInitialValue = 0
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))

	const expectedValue = 4
	const offset = 2

	commandUnderTest := new(DummyCommand).
		WithTarget(counterUnderTest).
		WithAttribute(increment, expectedValue+offset).
		WithAttribute(decrement, offset)

	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))

	commandUnderTest.Do()
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedValue))

	commandUnderTest.Undo()
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))
}

func TestMultipleDummyCommandSequence_CounterCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	const loopSize = 5
	commandSequence := make([]*DummyCommand, loopSize)

	counterUnderTest := new(Counter)
	const expectedInitialValue = 0
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))

	var buildIndex = 0
	for range commandSequence {
		commandSequence[buildIndex] = new(DummyCommand).
			WithTarget(counterUnderTest).
			WithAttribute(increment, buildIndex+1).
			WithAttribute(decrement, buildIndex)
		buildIndex++
	}

	for doIndex, command := range commandSequence {
		command.Do()
		g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, doIndex+1))
	}

	for undoIndex := loopSize - 1; undoIndex >= 0; undoIndex-- {
		commandSequence[undoIndex].Undo()
		g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, undoIndex))
	}
}
