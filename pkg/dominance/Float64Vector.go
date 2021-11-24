// Copyright (c) 2019 Australian Rivers Institute.

package dominance

type Float64Vector []float64

var _ Candidate = &Float64Vector{}

func NewFloat64(size int) *Float64Vector {
	vector := make(Float64Vector, size)
	return &vector
}

func (v *Float64Vector) IsComparable(otherCandidate Candidate) bool {
	if !isFloat64Vector(otherCandidate) {
		return false
	}

	otherCandidateAsVector := asFloat64Vector(otherCandidate)
	if !vectorLengthsMatch(v, otherCandidateAsVector) {
		return false
	}

	return true
}

func isFloat64Vector(otherCandidate Candidate) bool {
	switch otherCandidate.(type) {
	case *Float64Vector:
		return true
	default:
		return false
	}
}

func asFloat64Vector(otherCandidate Candidate) *Float64Vector {
	return otherCandidate.(*Float64Vector)
}

func vectorLengthsMatch(firstVector *Float64Vector, secondVector *Float64Vector) bool {
	return len(*firstVector) == len(*secondVector)
}

func (v *Float64Vector) Dominates(otherCandidate Candidate) bool {
	otherCandidateAsVector := *asFloat64Vector(otherCandidate)
	thisCandidateAsVector := *v

	// Efficient Learning Machines, 2015, Chapter 10.
	// X dominates Y ⇔ ∀ i:ℕ ∈ 0..|X|-1 ⦁ X[i] ≤ Y[i] ∧ ∃ j:ℕ ∈ 0..|X|-1 ⦁ X[i] < Y[i]

	// ∀ i:ℕ ∈ 0..|X|-1 ⦁ X[i] ≤ Y[i]
	for index := range thisCandidateAsVector {
		if thisCandidateAsVector[index] > otherCandidateAsVector[index] {
			return false
		}
	}

	// ∃ j:ℕ ∈ 0..|X|-1 ⦁ X[i] < Y[i]
	for index := range thisCandidateAsVector {
		if thisCandidateAsVector[index] < otherCandidateAsVector[index] {
			return true
		}
	}
	return false
}

func (v *Float64Vector) allLessThanValuesIn(otherCandidate Candidate) bool {
	otherCandidateAsVector := *asFloat64Vector(otherCandidate)
	thisCandidateAsVector := *v

	for index := range thisCandidateAsVector {
		if thisCandidateAsVector[index] >= otherCandidateAsVector[index] {
			return false
		}
	}
	return true
}

func (v *Float64Vector) anyLessThanValuesIn(otherCandidate Candidate) bool {
	otherCandidateAsVector := *asFloat64Vector(otherCandidate)
	thisCandidateAsVector := *v

	for index := range thisCandidateAsVector {
		if thisCandidateAsVector[index] < otherCandidateAsVector[index] {
			return true
		}
	}
	return false
}

func (v *Float64Vector) allEqualOrLessThanValuesIn(otherCandidate Candidate) bool {
	otherCandidateAsVector := *asFloat64Vector(otherCandidate)
	thisCandidateAsVector := *v

	for index := range otherCandidateAsVector {
		if !(thisCandidateAsVector[index] <= otherCandidateAsVector[index]) {
			return false
		}
	}
	return true
}

func (v *Float64Vector) IsDominatedBy(otherCandidate Candidate) bool {
	return otherCandidate.Dominates(v)
}

func (v *Float64Vector) DominancePresent(otherCandidate Candidate) bool {
	return v.Dominates(otherCandidate) || otherCandidate.Dominates(v)
}

func (v *Float64Vector) NoDominancePresent(otherCandidate Candidate) bool {
	return !(v.DominancePresent(otherCandidate))
}
