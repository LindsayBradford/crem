// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestSolution_MatchErrors_NilForIdenticalInitialised(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution("test")
	otherSolution := NewSolution("test")

	matchErrors := solutionUnderTest.MatchErrors(otherSolution)

	g.Expect(matchErrors).To(BeNil())
}

func TestSolution_MatchErrors_MismatchedIds(t *testing.T) {
	g := NewGomegaWithT(t)

	solutionUnderTest := NewSolution("test")
	mismatchedSolution := NewSolution("mismatch")

	matchErrors := solutionUnderTest.MatchErrors(mismatchedSolution)

	if matchErrors != nil {
		t.Log(matchErrors)
	}

	g.Expect(matchErrors).To(Not(BeNil()))
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

	solutionUnderTest := NewSolution("test")
	mismatchedSolution := NewSolution("test")

	solutionUnderTest.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 2)
	solutionUnderTest.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "matchingValues",
		Value: 0.0,
	}
	solutionUnderTest.DecisionVariables[1] = variableNew.EncodeableDecisionVariable{
		Name:  "mismatchingValues",
		Value: 1.0,
	}

	mismatchedSolution.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 2)
	mismatchedSolution.DecisionVariables[0] = variableNew.EncodeableDecisionVariable{
		Name:  "matchingValues",
		Value: 0,
	}
	mismatchedSolution.DecisionVariables[1] = variableNew.EncodeableDecisionVariable{
		Name:  "mismatchingValues",
		Value: 42.0,
	}

	const expectedErrors = 1

	matchErrors := solutionUnderTest.MatchErrors(mismatchedSolution)

	if matchErrors != nil {
		t.Log(matchErrors)
	}
	g.Expect(matchErrors).To(Not(BeNil()))
	g.Expect(matchErrors.Size()).To(BeNumerically(equalTo, expectedErrors))
}
