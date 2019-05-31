// Copyright (c) 2019 Australian Rivers Institute.

package dominance

type Candidate interface {
	IsComparable(otherCandidate Candidate) bool

	Dominates(otherCandidate Candidate) bool
	IsDominatedBy(otherCandidate Candidate) bool

	NoDominancePresent(otherCandidate Candidate) bool
}
