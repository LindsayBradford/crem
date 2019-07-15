// Copyright (c) 2019 Australian Rivers Institute.

package action

import (
	"testing"

	. "github.com/onsi/gomega"
)

const testType ManagementActionType = "test"
const testVariableName ModelVariableName = "testVariable"

func NewTestManagementAction() *SimpleManagementAction {
	action := new(SimpleManagementAction).
		WithPlanningUnit(1).
		WithType(testType).
		WithVariable(testVariableName, 0.5)

	return action
}

func TestSimpleManagementAction_Subscribe(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testSpyOne := new(spyObserver)
	testSpyTwo := new(spyObserver)

	actionUnderTest := NewTestManagementAction()

	// when
	actionUnderTest.Subscribe(testSpyOne, testSpyTwo)

	// then
	g.Expect(len(actionUnderTest.observers)).To(Equal(2))
}

func TestSimpleManagementAction_ToggleActivation(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testSpyOne := new(spyObserver)
	testSpyTwo := new(spyObserver)

	actionUnderTest := NewTestManagementAction()

	// when
	actionUnderTest.Subscribe(testSpyOne, testSpyTwo)

	// then
	g.Expect(testSpyOne.LastObserved()).To(BeNil())
	g.Expect(testSpyTwo.LastObserved()).To(BeNil())

	actionUnderTest.ToggleActivation()

	g.Expect(actionUnderTest.IsActive()).To(BeTrue())
	g.Expect(testSpyOne.LastObserved()).To(Equal(actionUnderTest))
	g.Expect(testSpyTwo.LastObserved()).To(Equal(actionUnderTest))

	// when
	testSpyOne.Reset()
	testSpyTwo.Reset()

	g.Expect(testSpyOne.LastObserved()).To(BeNil())
	g.Expect(testSpyTwo.LastObserved()).To(BeNil())

	actionUnderTest.ToggleActivation()

	// then
	g.Expect(actionUnderTest.IsActive()).To(BeFalse())
	g.Expect(testSpyOne.LastObserved()).To(Equal(actionUnderTest))
	g.Expect(testSpyTwo.LastObserved()).To(Equal(actionUnderTest))
}

func TestSimpleManagementAction_InitialisingActivation(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testSpyOne := new(spyObserver)
	testSpyTwo := new(spyObserver)

	actionUnderTest := NewTestManagementAction()

	// when
	actionUnderTest.Subscribe(testSpyOne, testSpyTwo)

	// then
	g.Expect(testSpyOne.LastObserved()).To(BeNil())
	g.Expect(testSpyTwo.LastObserved()).To(BeNil())

	actionUnderTest.InitialisingActivation()

	g.Expect(actionUnderTest.IsActive()).To(BeTrue())
	g.Expect(testSpyOne.LastObserved()).To(Equal(actionUnderTest))
	g.Expect(testSpyTwo.LastObserved()).To(Equal(actionUnderTest))

	// when
	testSpyOne.Reset()
	testSpyTwo.Reset()

	g.Expect(testSpyOne.LastObserved()).To(BeNil())
	g.Expect(testSpyTwo.LastObserved()).To(BeNil())

	expectedPanicCall := func() {
		actionUnderTest.InitialisingActivation()
	}

	// then
	g.Expect(expectedPanicCall).To(Panic())
	g.Expect(testSpyOne.LastObserved()).To(BeNil())
	g.Expect(testSpyTwo.LastObserved()).To(BeNil())
}

func TestSimpleManagementAction_ToggleActivationUnobserved(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testSpyOne := new(spyObserver)

	actionUnderTest := NewTestManagementAction()

	// when
	actionUnderTest.Subscribe(testSpyOne)

	// then
	g.Expect(testSpyOne.LastObserved()).To(BeNil())

	actionUnderTest.ToggleActivationUnobserved()

	g.Expect(actionUnderTest.IsActive()).To(BeTrue())
	g.Expect(testSpyOne.LastObserved()).To(BeNil())
}
