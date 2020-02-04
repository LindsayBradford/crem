// Copyright (c) 2020 Australian Rivers Institute.

package variable

import (
	"testing"

	. "github.com/onsi/gomega"
)

const keyOne = "keyOne"
const keyTwo = "keyTwo"
const keyThree = "keyThree"

func TestDecisionVariableMap_SortedKeys_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	mapUnderTest := buildMapUnderTest()
	sortedKeys := mapUnderTest.SortedKeys()

	g.Expect(len(mapUnderTest)).To(BeNumerically(equalTo, len(sortedKeys)))

	g.Expect(sortedKeys[0]).To(Equal(keyOne))
	g.Expect(sortedKeys[1]).To(Equal(keyThree))
	g.Expect(sortedKeys[2]).To(Equal(keyTwo))
}

func TestDecisionVariableMap_SortedKeyIndex_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	mapUnderTest := buildMapUnderTest()

	keyIndexOne := mapUnderTest.SortedKeyIndex(keyOne)
	g.Expect(keyIndexOne).To(BeNumerically(equalTo, 0))

	keyIndexTwo := mapUnderTest.SortedKeyIndex(keyTwo)
	g.Expect(keyIndexTwo).To(BeNumerically(equalTo, 2))

	keyIndexThree := mapUnderTest.SortedKeyIndex(keyThree)
	g.Expect(keyIndexThree).To(BeNumerically(equalTo, 1))

	noKeyIndex := mapUnderTest.SortedKeyIndex("noKeyPresent")
	g.Expect(noKeyIndex).To(BeNumerically(equalTo, KeyNotFound))
}

func buildMapUnderTest() DecisionVariableMap {
	builtMap := make(DecisionVariableMap, 3)

	builtMap[keyOne] = NewSimpleDecisionVariable(keyOne)
	builtMap[keyTwo] = NewSimpleDecisionVariable(keyTwo)
	builtMap[keyThree] = NewSimpleDecisionVariable(keyThree)

	return builtMap
}
