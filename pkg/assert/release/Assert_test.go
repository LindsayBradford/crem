// Copyright (c) 2019 Australian Rivers Institute.

package assert

import (
	"testing"

	. "github.com/onsi/gomega"
)

const (
	conditionIsHolding = true
	conditionIsBroken  = false
)

func TestAssertion_Holds_NoPanic(t *testing.T) {
	g := NewGomegaWithT(t)

	assertRunner := func() {
		That(conditionIsHolding).Holds()
	}

	g.Expect(assertRunner).ToNot(Panic())
}

func TestAssertion_DoesntHold_NoPanic(t *testing.T) {
	g := NewGomegaWithT(t)

	assertRunner := func() {
		That(conditionIsBroken).Holds()
	}

	g.Expect(assertRunner).ToNot(Panic())
}
