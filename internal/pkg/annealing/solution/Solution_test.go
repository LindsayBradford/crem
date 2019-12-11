// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	. "github.com/onsi/gomega"
)

const equalTo = "=="
const testSolutionId = "test"

func TestSolution_MatchErrors_NilForIdenticalInitialised(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution(testSolutionId)
	otherSolution := NewSolution(testSolutionId)

	matchErrors := solutionUnderTest.MatchErrors(otherSolution)

	g.Expect(matchErrors).To(BeNil())
}

func TestSolution_MatchErrors_MismatchedIds(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution(testSolutionId)
	mismatchedSolution := NewSolution(testSolutionId + "mismatch")

	matchErrors := solutionUnderTest.MatchErrors(mismatchedSolution)

	if matchErrors != nil {
		t.Log(matchErrors)
	}

	const expectedErrors = 1
	const expectedErrorMsg = "Solutions have mismatching Ids"

	g.Expect(matchErrors).To(Not(BeNil()))
	g.Expect(matchErrors.Size()).To(BeNumerically(equalTo, expectedErrors))
	g.Expect(matchErrors.SubError(0).Error()).To(ContainSubstring(expectedErrorMsg))
}

func TestSolution_MatchErrors_MissingVariables(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution("mine")
	mismatchedSolution := NewSolution("other") // expected error #1, mismatching ids

	solutionUnderTest.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 2)
	solutionUnderTest.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "only in mine", // expected error #2, variable only present here
		Value: 0,
	}
	solutionUnderTest.DecisionVariables[1] = variableNew.EncodeableDecisionVariable{
		Name:  "match",
		Value: 0,
	}

	mismatchedSolution.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 2)
	mismatchedSolution.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "match",
		Value: 0,
	}
	mismatchedSolution.DecisionVariables[1] = variableNew.EncodeableDecisionVariable{
		Name:  "only in other", // expected error #3, variable only present here
		Value: 0,
	}

	const expectedErrors = 3

	matchErrors := solutionUnderTest.MatchErrors(mismatchedSolution)

	if matchErrors != nil {
		t.Log(matchErrors)
	}
	g.Expect(matchErrors).To(Not(BeNil()))
	g.Expect(matchErrors.Size()).To(BeNumerically(equalTo, expectedErrors))
}

func TestSolution_MatchErrors_VariableValuesMismatch(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution(testSolutionId)
	mismatchedSolution := NewSolution(testSolutionId)

	solutionUnderTest.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 2)
	solutionUnderTest.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "matchingValues",
		Value: 0.0,
	}
	solutionUnderTest.DecisionVariables[1] = variableNew.EncodeableDecisionVariable{
		Name:  "mismatchingValues",
		Value: 1.0,
		ValuePerPlanningUnit: variableNew.PlanningUnitValues{
			variableNew.PlanningUnitValue{
				PlanningUnit: 42,
				Value:        1.0,
			},
		},
	}

	mismatchedSolution.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 2)
	mismatchedSolution.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "matchingValues",
		Value: 0,
	}
	mismatchedSolution.DecisionVariables[1] = variableNew.EncodeableDecisionVariable{
		Name:  "mismatchingValues",
		Value: 42.0,
		ValuePerPlanningUnit: variableNew.PlanningUnitValues{
			variableNew.PlanningUnitValue{
				PlanningUnit: 42,
				Value:        42.0,
			},
		},
	}

	const expectedErrors = 1
	const expectedErrorMsg = "variable [mismatchingValues] has mismatching values"

	matchErrors := solutionUnderTest.MatchErrors(mismatchedSolution)

	if matchErrors != nil {
		t.Log(matchErrors)
	}
	g.Expect(matchErrors).To(Not(BeNil()))
	g.Expect(matchErrors.Size()).To(BeNumerically(equalTo, expectedErrors))
	g.Expect(matchErrors.SubError(0).Error()).To(ContainSubstring(expectedErrorMsg))
}

func TestSolution_MatchErrors_VariableValueMatchesSumOfPlanningUnits_NoMatchErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution(testSolutionId)

	solutionUnderTest.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 1)
	solutionUnderTest.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "mismatchingValues",
		Value: 2.25,
		ValuePerPlanningUnit: variableNew.PlanningUnitValues{
			variableNew.PlanningUnitValue{
				PlanningUnit: 0,
				Value:        1.5,
			},
			variableNew.PlanningUnitValue{
				PlanningUnit: 1,
				Value:        0.75,
			},
		},
	}

	matchErrors := solutionUnderTest.MatchErrors(solutionUnderTest)

	if matchErrors != nil {
		t.Log(matchErrors)
	}

	g.Expect(matchErrors).To(BeNil())
}

func TestSolution_MatchErrors_VariableValueDoesntMatchSumOfPlanningUnits_MatchErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution(testSolutionId)

	solutionUnderTest.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 1)
	solutionUnderTest.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "mismatchingValues",
		Value: 3.0,
		ValuePerPlanningUnit: variableNew.PlanningUnitValues{
			variableNew.PlanningUnitValue{
				PlanningUnit: 0,
				Value:        1.5,
			},
			variableNew.PlanningUnitValue{
				PlanningUnit: 1,
				Value:        0.75,
			},
		},
	}

	const expectedErrors = 1
	const expectedErrorMsg = "but sum of planning units is"

	matchErrors := solutionUnderTest.MatchErrors(solutionUnderTest)

	if matchErrors != nil {
		t.Log(matchErrors)
	}

	g.Expect(matchErrors).To(Not(BeNil()))
	g.Expect(matchErrors.Size()).To(BeNumerically(equalTo, expectedErrors))
	g.Expect(matchErrors.SubError(0).Error()).To(ContainSubstring(expectedErrorMsg))
}
