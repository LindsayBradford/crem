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
		WithPlanningUnit("testPu").
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

func TestSimpleManagementAction_Activate(t *testing.T) {
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

	actionUnderTest.Activate()

	g.Expect(actionUnderTest.IsActive()).To(BeTrue())
	g.Expect(testSpyOne.LastObserved()).To(Equal(actionUnderTest))
	g.Expect(testSpyTwo.LastObserved()).To(Equal(actionUnderTest))
}

func TestSimpleManagementAction_Deactivate(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testSpyOne := new(spyObserver)

	actionUnderTest := NewTestManagementAction()

	// when
	actionUnderTest.Subscribe(testSpyOne)

	// then
	g.Expect(testSpyOne.LastObserved()).To(BeNil())

	actionUnderTest.Deactivate()

	g.Expect(testSpyOne.LastObserved()).To(BeNil())
}

func TestSimpleManagementAction_ToggleActivation(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testSpyOne := new(spyObserver)

	actionUnderTest := NewTestManagementAction()

	// when
	actionUnderTest.Subscribe(testSpyOne)

	// then
	g.Expect(testSpyOne.LastObserved()).To(BeNil())

	actionUnderTest.ToggleActivation()

	g.Expect(actionUnderTest.IsActive()).To(BeTrue())
	g.Expect(testSpyOne.LastObserved()).To(Equal(actionUnderTest))
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
