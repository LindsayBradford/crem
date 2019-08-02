// Copyright (c) 2019 Australian Rivers Institute.

package action

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"testing"

	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestManagementActions_Initialise(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	actionsUnderTest := new(ModelManagementActions)
	// when
	actionsUnderTest.Initialise()

	// then
	g.Expect(actionsUnderTest.lastApplied).To(BeNil())
}

func TestManagementActions_Add(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ModelManagementActions)
	actionsUnderTest.Initialise()

	// when
	dummyAction1 := buildDummyAction(1)
	actionsUnderTest.Add(dummyAction1)

	// then
	g.Expect(len(actionsUnderTest.actions)).To(BeNumerically("==", 1))
	g.Expect(actionsUnderTest.actions).To(ContainElement(dummyAction1))

	// when
	dummyAction2 := buildDummyAction(2)
	actionsUnderTest.Add(dummyAction2)

	// then
	g.Expect(len(actionsUnderTest.actions)).To(BeNumerically("==", 2))
	g.Expect(actionsUnderTest.actions).To(ContainElement(dummyAction1))
	g.Expect(actionsUnderTest.actions).To(ContainElement(dummyAction2))

	g.Expect(actionsUnderTest.lastApplied).To(BeNil())
}

const ManagementActionsTestType ManagementActionType = "ManagementActionsTestType"

func buildDummyAction(planningUnit planningunit.Id) ManagementAction {
	newAction := new(SimpleManagementAction).
		WithPlanningUnit(planningUnit).
		WithType(ManagementActionsTestType).
		WithVariable("dummyVar", 1)

	return newAction
}

func TestManagementActions_RandomlyToggleOneActivation_NonePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ModelManagementActions)
	actionsUnderTest.Initialise()

	// when
	actionsUnderTest.RandomlyToggleOneActivation()

	// then

	g.Expect(actionsUnderTest.lastApplied).To(Equal(NullManagementAction))
}

func TestManagementActions_RandomlyToggleOneActivation_OnePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ModelManagementActions)
	actionsUnderTest.Initialise()

	dummyAction1 := buildDummyAction(1)
	actionsUnderTest.Add(dummyAction1)

	// when
	actionsUnderTest.RandomlyToggleOneActivation()

	// then
	g.Expect(actionsUnderTest.lastApplied).To(Equal(dummyAction1))
	g.Expect(dummyAction1.IsActive()).To(BeTrue())
}

func TestManagementActions_RandomlyInitialise(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ModelManagementActions)
	actionsUnderTest.Initialise()

	actionSpy := new(spyObserver)

	dummyActions := []ManagementAction{
		buildDummyAction(1),
		buildDummyAction(2),
		buildDummyAction(3),
		buildDummyAction(4),
		buildDummyAction(4),
	}

	for _, action := range dummyActions {
		action.Subscribe(actionSpy)
		actionsUnderTest.Add(action)
	}

	// when
	actionsUnderTest.RandomlyInitialise()
	expectedActivations := actionSpy.ObservationsCounted()

	// then
	var actualActivations uint
	for _, action := range dummyActions {
		if action.IsActive() {
			actualActivations++
		}
	}
	t.Log(fmt.Sprintf("%d dummy actions randomly activated", actualActivations))

	g.Expect(actualActivations).To(BeNumerically(equalTo, expectedActivations))
}

func TestManagementActions_UndoLastActivationToggleUnobserved(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	actionsUnderTest := new(ModelManagementActions)
	actionsUnderTest.Initialise()

	dummyAction := buildDummyAction(1)

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
	actionsUnderTest.ToggleLastActivationUnobserved()

	// then
	g.Expect(dummyAction.IsActive()).To(BeFalse())
	g.Expect(actionSpy.LastObserved()).To(BeNil())
}
