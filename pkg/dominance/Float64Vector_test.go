// Copyright (c) 2019 Australian Rivers Institute.

package dominance

import (
	// "github.com/LindsayBradford/crem/internal/pkg/rand"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	equalTo = "=="

	dominate    = "dominate"
	notDominate = "not dominate"

	beDominatedBy    = "be dominated by"
	notBeDominatedBy = "not be dominated by"
)

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
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = -1e-3

	// then
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 0
	indexedOtherCandidateVector[0] = 1e-3

	// then
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 1e-3

	// then
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedOtherCandidateVector[1] = 2e-3

	// then
	t.Logf("Expected vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[1] = 2e-3

	// then
	t.Logf("Expected vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedOtherCandidateVector[2] = 3e-3

	// then
	t.Logf("Expected vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[2] = 3e-3

	// then
	t.Logf("Expected vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[2] = 1e-3

	// then
	t.Logf("Expected vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 3e-3

	// then
	t.Logf("Expected vectorUnderTest %v ts %s otherCandidateVector: %v", vectorUnderTest, notDominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 0
	indexedVectorUnderTest[1] = 0
	indexedVectorUnderTest[2] = 0

	// then
	t.Logf("Expected vectorUnderTest %v ts %s otherCandidateVector: %v", vectorUnderTest, dominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeTrue())

	// when
	indexedVectorUnderTest[0] = 1e-3
	indexedVectorUnderTest[1] = 1e-3
	indexedVectorUnderTest[2] = 1e-3

	indexedOtherCandidateVector[0] = 2e-3
	indexedOtherCandidateVector[1] = 2e-3
	indexedOtherCandidateVector[2] = 2e-3

	// then
	t.Logf("Expected vectorUnderTest %v ts %s otherCandidateVector: %v", vectorUnderTest, dominate, otherCandidateVector)
	g.Expect(vectorUnderTest.Dominates(otherCandidateVector)).To(BeTrue())

}

func TestFloat64Vector_Dominates_ValidlyIdentifiesNoDominancePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 3
	vectorUnderTest := NewFloat64(expectedSize)

	otherCandidateVector := NewFloat64(expectedSize)

	// then
	t.Logf("Expecting no dominamce present between vectorUnderTest: %v and otherCandidateVector: %v", vectorUnderTest, otherCandidateVector)
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
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notBeDominatedBy, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeFalse())

	// when
	indexedOtherCandidateVector[0] = 1e-3

	// then
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notBeDominatedBy, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 1e-3

	// then
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, notBeDominatedBy, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeFalse())

	// when
	indexedVectorUnderTest[0] = 2e-3
	indexedVectorUnderTest[1] = 1e-3
	indexedVectorUnderTest[2] = 1e-3

	// then
	t.Logf("Expecting vectorUnderTest %v to %s otherCandidateVector: %v", vectorUnderTest, beDominatedBy, otherCandidateVector)
	g.Expect(vectorUnderTest.IsDominatedBy(otherCandidateVector)).To(BeTrue())
}
