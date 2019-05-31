// Copyright (c) 2019 Australian Rivers Institute.

package dominance

import (
	// "github.com/LindsayBradford/crem/internal/pkg/rand"
	"testing"

	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestFloat64Vector_New(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 4

	// when
	vectorUnderTest := *NewFloat64(expectedSize)

	// then
	g.Expect(len(vectorUnderTest)).To(BeNumerically(equalTo, expectedSize))

	for index := 0; index < expectedSize; index++ {
		g.Expect(vectorUnderTest[index]).To(BeNumerically(equalTo, 0))
	}
}

func TestFloat64Vector_IsComparable_SameSizeSameTypeVectorsComparable(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 4
	vectorUnderTest := NewFloat64(expectedSize)
	comparableVector := NewFloat64(expectedSize)

	// when
	actualComparableValue := vectorUnderTest.IsComparable(comparableVector)

	// then
	g.Expect(actualComparableValue).To(BeTrue())

}

func TestFloat64Vector_IsComparable_DifferentSizedVectorsNotComparable(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testVectorSize := 4
	vectorUnderTest := NewFloat64(testVectorSize)
	invalidVectorSize := 3
	expectedInvalidVector := NewFloat64(invalidVectorSize)

	// when
	actualComparableValue := vectorUnderTest.IsComparable(expectedInvalidVector)

	// then
	g.Expect(actualComparableValue).To(BeFalse())
}

type dummyCandidate []int

func (v *dummyCandidate) IsComparable(otherCandidate Candidate) bool       { return false }
func (v *dummyCandidate) Dominates(otherCandidate Candidate) bool          { return false }
func (v *dummyCandidate) IsDominatedBy(otherCandidate Candidate) bool      { return false }
func (v *dummyCandidate) NoDominancePresent(otherCandidate Candidate) bool { return true }

func TestFloat64Vector_IsComparable_TypeDifferentVectorsNotComparable(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 4
	vectorUnderTest := NewFloat64(expectedSize)
	dummyVector := make(dummyCandidate, expectedSize)

	// when
	actualComparableValue := vectorUnderTest.IsComparable(&dummyVector)

	// then
	g.Expect(actualComparableValue).To(BeFalse())
}

func TestFloat64Vector_Dominates_ValidlyIdentifiesDomination(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 3
	vectorUnderTest := NewFloat64(expectedSize)
	indexedVectorUnderTest := *vectorUnderTest

	otherCandidateVector := NewFloat64(expectedSize)
	indexedOtherCandidateVector := *otherCandidateVector

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedOtherCandidateVector[0] = 1e-10

	// then
	t.Logf("Expected true: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeTrue())

	// when
	indexedVectorUnderTest[0] = 1e-10

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedOtherCandidateVector[1] = 2e-10

	// then
	t.Logf("Expected true: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeTrue())

	// when
	indexedVectorUnderTest[1] = 2e-10

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedOtherCandidateVector[2] = 3e-10

	// then
	t.Logf("Expected true: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeTrue())

	// when
	indexedVectorUnderTest[2] = 3e-10

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[2] = 1e-10

	// then
	t.Logf("Expected true: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeTrue())

	// when
	indexedVectorUnderTest[0] = 3e-10

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())
}

func TestFloat64Vector_Dominates_ValidlyIdentifiesNoDominancePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 3
	vectorUnderTest := NewFloat64(expectedSize)
	// indexedVectorUnderTest := *vectorUnderTest

	otherCandidateVector := NewFloat64(expectedSize)
	// indexedOtherCandidateVector := *otherCandidateVector

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.NoDominancePresent(otherCandidateVector)).To(BeTrue())
}

func TestFloat64Vector_Dominates_ValidlyIdentifiesIsDominatedBy(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 3
	vectorUnderTest := NewFloat64(expectedSize)
	indexedVectorUnderTest := *vectorUnderTest

	otherCandidateVector := NewFloat64(expectedSize)
	indexedOtherCandidateVector := *otherCandidateVector

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeFalse())

	// when
	indexedOtherCandidateVector[0] = 1e-10

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 1e-10

	// then
	t.Logf("Expected false: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 2e-10

	// then
	t.Logf("Expected true: vectorUnderTest: %v, otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeTrue())

}
