// Copyright (c) 2019 Australian Rivers Institute.

package command

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestEmptyCompositeCommand_DoUndo_DoesNotPanic(t *testing.T) {
	g := NewGomegaWithT(t)

	testCommand := NewCompositeCommand()

	doRunner := func() {
		testCommand.Do()
	}

	g.Expect(doRunner).ToNot(Panic())

	undoRunner := func() {
		testCommand.Undo()
	}

	g.Expect(undoRunner).ToNot(Panic())
}

func TestCompositeCommand_DoUndo_CounterCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedInitialValue = 2
	counterUnderTest := new(Counter)
	counterUnderTest.value = expectedInitialValue
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))

	const value1 = 3
	command1 := new(dummyCommand).
		WithTarget(counterUnderTest).
		WithAttribute(increment, value1)

	const value2 = 5
	command2 := new(dummyCommand).
		WithTarget(counterUnderTest).
		WithAttribute(decrement, value2)

	const value3 = 8
	command3 := new(dummyCommand).
		WithTarget(counterUnderTest).
		WithAttribute(increment, value3)

	commandUnderTest := NewCompositeCommand().ComposedOf(command1, command2, command3)

	const expectedValue = expectedInitialValue + value1 - value2 + value3

	commandUnderTest.Do()
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedValue))

	commandUnderTest.Undo()
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))
}
