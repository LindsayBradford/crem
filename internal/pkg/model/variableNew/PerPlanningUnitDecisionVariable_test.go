// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	. "github.com/onsi/gomega"
)

const approx = "~"

func TestVariable_Initial(t *testing.T) {
	g := NewGomegaWithT(t)

	variableUnderTest := NewPerPlanningUnitDecisionVariable()
	g.Expect(variableUnderTest.Value()).To(BeNumerically("==", 0))
}

func TestVariable_SetOnePlanningUnit(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	const planningUnitUnderTest = 42
	const expectedValue = 5

	variableUnderTest := NewPerPlanningUnitDecisionVariable()

	// when

	variableUnderTest.SetPlanningUnitValue(planningUnitUnderTest, expectedValue)

	// then

	g.Expect(variableUnderTest.Value()).To(BeNumerically(equalTo, expectedValue))
	g.Expect(variableUnderTest.PlanningUnitValue(planningUnitUnderTest)).To(BeNumerically(equalTo, expectedValue))
}

func TestVariable_SetSeveralPlanningUnits_WorksAsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	rnd := rand.NewTimeSeeded()
	const maxLoops = 10
	loops := 1 + rnd.Intn(maxLoops)

	variableUnderTest := NewPerPlanningUnitDecisionVariable()

	var total float64
	for i := 0; i < loops; i++ {
		value := rnd.Float64Unitary()
		planningUnit := planningunit.Id(i)
		variableUnderTest.SetPlanningUnitValue(planningUnit, value)
		total += value
	}

	g.Expect(variableUnderTest.Value()).To(BeNumerically(approx, total))
}

func TestVariable_SetPlanningUnitRepeatedly_WorksAsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	const planningUnitUnderTest = 42

	rnd := rand.NewTimeSeeded()
	const maxLoops = 10
	loops := 1 + rnd.Intn(maxLoops)

	variableUnderTest := NewPerPlanningUnitDecisionVariable()

	var lastValue float64
	for i := 0; i < loops; i++ {
		value := rnd.Float64Unitary()
		variableUnderTest.SetPlanningUnitValue(planningUnitUnderTest, value)
		lastValue = value
	}

	g.Expect(variableUnderTest.PlanningUnitValue(planningUnitUnderTest)).To(BeNumerically(approx, lastValue))
	g.Expect(variableUnderTest.Value()).To(BeNumerically(approx, lastValue))
}
