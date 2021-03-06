// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	equalTo = "=="

	expectedDefaultEncoding        = "0:0"
	expectedBoundaryValuesEncoding = "8000000000000001:1"
)

func TestBooleanArchive_New(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 130
	expectedArchiveSize := 3

	// when
	archiveUnderTest := New(expectedSize)

	// then
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, expectedSize))
	g.Expect(archiveUnderTest.ArchiveLen()).To(BeNumerically(equalTo, expectedArchiveSize))

	for index := 0; index < expectedSize; index++ {
		g.Expect(archiveUnderTest.Value(index)).To(BeFalse())
	}
}

func TestBooleanArchive_SetValue_OverArchiveRage(t *testing.T) {
	g := NewGomegaWithT(t)
	random := rand.NewTimeSeeded()

	// given
	expectedSize := 130
	archiveUnderTest := New(expectedSize)

	numberToSetTrue := 7
	expectedTrueIndexes := make([]int, numberToSetTrue)
	for current := 0; current < numberToSetTrue; current++ {
		indexToSetTrue := random.Intn(expectedSize)
		for previous := 0; previous < current; previous++ {
			duplicateIndexFound := true
			for duplicateIndexFound {
				if expectedTrueIndexes[previous] == indexToSetTrue {
					indexToSetTrue = random.Intn(expectedSize)
				} else {
					duplicateIndexFound = false
				}
			}
		}
		expectedTrueIndexes[current] = indexToSetTrue
	}
	t.Logf("Archive indexes that should be set to true: %v", expectedTrueIndexes)

	// when
	for current := 0; current < numberToSetTrue; current++ {
		archiveUnderTest.SetValue(expectedTrueIndexes[current], true)
	}

	actualTrueIndexes := make([]int, 0)
	falseCount := 0
	for index := 0; index < expectedSize; index++ {
		if archiveUnderTest.Value(index) == false {
			falseCount++
		} else {
			actualTrueIndexes = append(actualTrueIndexes, index)
		}
	}
	t.Logf("Archive indexes are set to true: %v", actualTrueIndexes)

	// then
	g.Expect(expectedTrueIndexes).To(ConsistOf(actualTrueIndexes))
	g.Expect(falseCount).To(BeNumerically(equalTo, expectedSize-numberToSetTrue))
}

func TestBooleanArchive_SetValue_Toggling(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedSize := 5
	archiveUnderTest := New(expectedSize)
	indexToTest := 2

	// given
	valuesToAssign := []bool{false, false, true, true, false, false, true, true}

	for _, testValue := range valuesToAssign {
		// when
		archiveUnderTest.SetValue(indexToTest, testValue)

		// then
		if testValue {
			g.Expect(archiveUnderTest.Value(indexToTest)).To(BeTrue())
		} else {
			g.Expect(archiveUnderTest.Value(indexToTest)).To(BeFalse())
		}
	}
}

func TestBooleanArchive_SetValue_OutsideValidRange(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 5
	archiveUnderTest := New(expectedSize)

	// when
	indexToTest := 4
	archiveUnderTest.SetValue(indexToTest, true)

	// then
	g.Expect(archiveUnderTest.Value(indexToTest)).To(BeTrue())

	// when
	indexToTest = 5
	outOfBoundsSet := func() {
		archiveUnderTest.SetValue(indexToTest, true)
	}

	// then
	g.Expect(outOfBoundsSet).To(Panic())
}

func TestBooleanArchive_Value_OutsideValidRange(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 5
	archiveUnderTest := New(expectedSize)

	// when
	indexToTest := 4
	expectedValue := archiveUnderTest.Value(indexToTest)

	// then
	g.Expect(expectedValue).To(BeFalse())

	// when
	indexToTest = 5
	outOfBoundsValue := func() {
		archiveUnderTest.Value(indexToTest)
	}

	// then
	g.Expect(outOfBoundsValue).To(Panic())
}

func TestBooleanArchive_IsEquivalentTo_ValidResponses(t *testing.T) {
	g := NewGomegaWithT(t)
	random := rand.NewTimeSeeded()

	expectedSize := 200
	firstArchiveUnderTest := New(expectedSize)
	secondArchiveUnderTest := New(expectedSize)

	// given

	numberToSetTrue := 10
	expectedTrueIndexes := make([]int, numberToSetTrue)
	for current := 0; current < numberToSetTrue; current++ {
		indexToSetTrue := random.Intn(expectedSize)
		for previous := 0; previous < current; previous++ {
			duplicateIndexFound := true
			for duplicateIndexFound {
				if expectedTrueIndexes[previous] == indexToSetTrue {
					indexToSetTrue = random.Intn(expectedSize)
				} else {
					duplicateIndexFound = false
				}
			}
		}
		expectedTrueIndexes[current] = indexToSetTrue
	}
	t.Logf("Archive indexes that should be set to true: %v", expectedTrueIndexes)

	for current := 0; current < numberToSetTrue; current++ {
		// when
		firstArchiveUnderTest.SetValue(expectedTrueIndexes[current], true)

		// then
		g.Expect(firstArchiveUnderTest.IsEquivalentTo(secondArchiveUnderTest)).To(BeFalse())

		// when
		secondArchiveUnderTest.SetValue(expectedTrueIndexes[current], true)

		// then
		g.Expect(firstArchiveUnderTest.IsEquivalentTo(secondArchiveUnderTest)).To(BeTrue())
	}
}

func TestBooleanArchive_IsEquivalentTo_InvalidTest(t *testing.T) {
	g := NewGomegaWithT(t)

	baseSize := 200
	baseArchiveUnderTest := New(baseSize)
	sameSizedArchiveUnderTest := New(baseSize)

	g.Expect(baseArchiveUnderTest.IsEquivalentTo(sameSizedArchiveUnderTest)).To(BeTrue())

	differentSize := baseSize + 1
	differentlySizedArchiveUnderTest := New(differentSize)

	g.Expect(baseArchiveUnderTest.IsEquivalentTo(differentlySizedArchiveUnderTest)).To(BeFalse())
}

func TestBooleanArchive_EncodingDefault_Success(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	expectedSize := 100
	archiveUnderTest := New(expectedSize)

	// when
	encoding := archiveUnderTest.Encoding()
	t.Log(encoding)

	// then

	g.Expect(encoding).To(Equal(expectedDefaultEncoding))
}

func TestBooleanArchive_Encoding_DelimiterBoundariesCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	expectedSize := 100
	archiveUnderTest := New(expectedSize)

	// when

	archiveUnderTest.SetValue(0, true)
	encodingUnderTest := archiveUnderTest.Encoding()

	// then

	t.Log(encodingUnderTest)
	g.Expect(encodingUnderTest).To(Equal("1:0"))

	// when

	archiveUnderTest.SetValue(64, true)
	encodingUnderTest = archiveUnderTest.Encoding()

	// then

	t.Log(encodingUnderTest)
	g.Expect(encodingUnderTest).To(Equal("1:1"))

	// when

	archiveUnderTest.SetValue(63, true)
	threeValueEncoding := archiveUnderTest.Encoding()

	// then

	t.Log(threeValueEncoding)
	g.Expect(threeValueEncoding).To(Equal(expectedBoundaryValuesEncoding))
}

func TestBooleanArchive_Decode_Success(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	expectedSize := 100
	archiveUnderTest := New(expectedSize)

	// when

	archiveUnderTest.Decode(expectedBoundaryValuesEncoding)
	encodingUnderTest := archiveUnderTest.Encoding()

	// then

	t.Log(encodingUnderTest)
	g.Expect(encodingUnderTest).To(Equal(expectedBoundaryValuesEncoding))

	// when

	archiveUnderTest.Decode(expectedDefaultEncoding)
	encodingUnderTest = archiveUnderTest.Encoding()

	// then

	t.Log(encodingUnderTest)
	g.Expect(encodingUnderTest).To(Equal(expectedDefaultEncoding))

	// when

	outOfBoundsEncoding := "8000000000000001:1"
	archiveUnderTest.Decode(outOfBoundsEncoding)
	encodingUnderTest = archiveUnderTest.Encoding()

	// then

	t.Log(encodingUnderTest)
	g.Expect(encodingUnderTest).To(Equal(expectedBoundaryValuesEncoding))
}

func TestBooleanArchive_Decode_OutOfBoundsBecomeBounded(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	expectedSize := 100
	archiveUnderTest := New(expectedSize)

	// when

	outOfBoundsEncoding := "8000000000000001:8000000000000001" // bit 127 is true here, but we only use range 0-99
	archiveUnderTest.Decode(outOfBoundsEncoding)
	encodingUnderTest := archiveUnderTest.Encoding()

	// then  -- decoding expected to zero out flagged but unused bits (size of 100 means bits 100 thru 127 are unused).

	t.Log(encodingUnderTest)
	g.Expect(encodingUnderTest).To(Equal(expectedBoundaryValuesEncoding))
}
