// Copyright (c) 2019 Australian Rivers Institute.

package action

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestManagementActions_Initialise(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	actionsUnderTest := new(ManagementActions)

	// when
	actionsUnderTest.Initialise()

	// then
	g.Expect(actionsUnderTest.lastApplied).To(BeNil())
}

func TestManagementActions_Add(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ManagementActions)
	actionsUnderTest.Initialise()

	// when
	dummyAction1 := buildDummyAction("dummy1")
	actionsUnderTest.Add(dummyAction1)

	// then
	g.Expect(len(actionsUnderTest.actions)).To(BeNumerically("==", 1))
	g.Expect(actionsUnderTest.actions).To(ContainElement(dummyAction1))

	// when
	dummyAction2 := buildDummyAction("dummy2")
	actionsUnderTest.Add(dummyAction2)

	// then
	g.Expect(len(actionsUnderTest.actions)).To(BeNumerically("==", 2))
	g.Expect(actionsUnderTest.actions).To(ContainElement(dummyAction1))
	g.Expect(actionsUnderTest.actions).To(ContainElement(dummyAction2))

	g.Expect(actionsUnderTest.lastApplied).To(BeNil())
}

const ManagementActionsTestType ManagementActionType = "ManagementActionsTestType"

func buildDummyAction(planningUnit string) ManagementAction {
	newAction := new(SimpleManagementAction).
		WithPlanningUnit(planningUnit).
		WithType(ManagementActionsTestType).
		WithVariable("dummyVar", 1)

	return newAction
}

func TestManagementActions_RandomlyToggleOneActivation_NonePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ManagementActions)
	actionsUnderTest.Initialise()

	// when
	actionsUnderTest.RandomlyToggleOneActivation()

	// then

	g.Expect(actionsUnderTest.lastApplied).To(Equal(NullManagementAction))
}

func TestManagementActions_RandomlyToggleOneActivation_OnePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ManagementActions)
	actionsUnderTest.Initialise()

	dummyAction1 := buildDummyAction("dummy1")
	actionsUnderTest.Add(dummyAction1)

	// when
	actionsUnderTest.RandomlyToggleOneActivation()

	// then
	g.Expect(actionsUnderTest.lastApplied).To(Equal(dummyAction1))
	g.Expect(dummyAction1.IsActive()).To(BeTrue())
}

func TestManagementActions_RandomlyToggleAllActivations(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ManagementActions)
	actionsUnderTest.Initialise()

	actionSpy := new(spyObserver)

	dummyActions := []ManagementAction{
		buildDummyAction("dummy1"),
		buildDummyAction("dummy2"),
		buildDummyAction("dummy3"),
	}

	for _, action := range dummyActions {
		action.Subscribe(actionSpy)
		actionsUnderTest.Add(action)
	}

	// when
	actionsUnderTest.RandomlyToggleAllActivations()

	// then
	var actualActivations uint
	for _, action := range dummyActions {
		if action.IsActive() {
			actualActivations++
		}
	}
	t.Log(fmt.Sprintf("%d dummy actions randomly activated", actualActivations))

	g.Expect(actualActivations).To(BeNumerically("==", actionSpy.ObservationsCounted()))
}

func TestManagementActions_UndoLastActivationToggleUnobserved(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ManagementActions)
	actionsUnderTest.Initialise()

	dummyAction := buildDummyAction("dummy")

	actionSpy := new(spyObserver)

	dummyAction.Subscribe(actionSpy)
	actionsUnderTest.Add(dummyAction)

	// when
	actionsUnderTest.RandomlyToggleOneActivation()

	// then
	g.Expect(actionSpy.LastObserved()).To(Equal(dummyAction))
	g.Expect(dummyAction.IsActive()).To(BeTrue())

	// when
	actionSpy.Reset()
	actionsUnderTest.UndoLastActivationToggleUnobserved()

	// then
	g.Expect(dummyAction.IsActive()).To(BeFalse())
	g.Expect(actionSpy.LastObserved()).To(BeNil())
}
