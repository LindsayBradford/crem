// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package errors

import (
	"errors"

	. "github.com/onsi/gomega"
)
import "testing"

func TestCompositeError_add(t *testing.T) {
	g := NewGomegaWithT(t)

	errorUnderTest := NewComposite("testingComposite")

	expectedSubError0 := errors.New("subError0")

	g.Expect(
		errorUnderTest.Size()).To(BeZero(),
		"A new composite error should have zero size")

	errorUnderTest.Add(expectedSubError0)

	g.Expect(
		errorUnderTest.Size()).To(BeIdenticalTo(1),
		"A composite error size should grow by one after add")

	g.Expect(
		errorUnderTest.SubError(0)).To(BeIdenticalTo(expectedSubError0),
		"Composite error should store first sub-error at index 0")
}
