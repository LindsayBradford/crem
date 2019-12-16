// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"testing"

	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestCommand_DoUndo_NoPanic(t *testing.T) {
	g := NewGomegaWithT(t)

	variableUnderTest := NewPerPlanningUnitDecisionVariable()
	testCommand := new(ChangePerPlanningUnitDecisionVariableCommand).ForVariable(variableUnderTest)

	doRunner := func() {
		testCommand.Do()
	}

	g.Expect(doRunner).ToNot(Panic())

	undoRunner := func() {
		testCommand.Undo()
	}

	g.Expect(undoRunner).ToNot(Panic())
}

func TestCommand_DoUndo_ChangesAsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	const oldValue = 5
	const newValue = 10
	const planningUnit = 42

	variableUnderTest := NewPerPlanningUnitDecisionVariable()
	variableUnderTest.SetPlanningUnitValue(planningUnit, oldValue)

	g.Expect(variableUnderTest.Value()).To(BeNumerically(equalTo, oldValue))
	g.Expect(variableUnderTest.PlanningUnitValue(planningUnit)).To(BeNumerically(equalTo, oldValue))

	// when

	commandUnderTest := new(ChangePerPlanningUnitDecisionVariableCommand).
		ForVariable(variableUnderTest).
		InPlanningUnit(planningUnit).
		WithChange(newValue)

	// then

	g.Expect(variableUnderTest.Value()).To(BeNumerically(equalTo, oldValue))
	g.Expect(variableUnderTest.PlanningUnitValue(planningUnit)).To(BeNumerically(equalTo, oldValue))

	commandUnderTest.Do()

	g.Expect(variableUnderTest.Value()).To(BeNumerically(equalTo, oldValue+newValue))
	g.Expect(variableUnderTest.PlanningUnitValue(planningUnit)).To(BeNumerically(equalTo, oldValue+newValue))

	commandUnderTest.Undo()

	g.Expect(variableUnderTest.Value()).To(BeNumerically(equalTo, oldValue))
	g.Expect(variableUnderTest.PlanningUnitValue(planningUnit)).To(BeNumerically(equalTo, oldValue))
}
