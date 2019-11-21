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

func TestDummyCommand_DoUndo_CounterCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	counterUnderTest := new(Counter)
	const expectedInitialValue = 0
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))

	const expectedValue = 4
	const offset = 2

	commandUnderTest := new(dummyCommand).
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
	commandSequence := make([]*dummyCommand, loopSize)

	counterUnderTest := new(Counter)
	const expectedInitialValue = 0
	g.Expect(counterUnderTest.value).To(BeNumerically(equalTo, expectedInitialValue))

	var buildIndex = 0
	for range commandSequence {
		commandSequence[buildIndex] = new(dummyCommand).
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
