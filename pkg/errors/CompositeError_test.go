// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package errors

import (
	"errors"
	"fmt"

	. "github.com/onsi/gomega"
)
import "testing"

func TestCompositeError_add(t *testing.T) {
	g := NewGomegaWithT(t)

	errorUnderTest := New("testingComposite")

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

	g.Expect(
		errorUnderTest.Error()).To(Equal("subError0"))

	expectedSubError1 := errors.New("subError1")
	errorUnderTest.Add(expectedSubError1)

	g.Expect(
		errorUnderTest.Size()).To(BeIdenticalTo(2),
		"A composite error size should grow by one after add")

	g.Expect(
		errorUnderTest.SubError(1)).To(BeIdenticalTo(expectedSubError1),
		"Composite error should store second sub-error at index 1")

	actualCompositeErrorString := errorUnderTest.Error()

	g.Expect(actualCompositeErrorString).To(ContainSubstring(expectedSubError0.Error()),
		"Composite error should return error string of its first sub-error")

	g.Expect(actualCompositeErrorString).To(ContainSubstring(expectedSubError1.Error()),
		"Composite error should return error string of its second sub-error")
}

func ExampleCompositeError_Add() {
	newComposite := New("error prefix")
	newComposite.Add(errors.New("first error"))
	newComposite.Add(errors.New("second error"))
	fmt.Printf("%v", newComposite)

	// Output: error prefix, composed of: [
	// 	first error
	// 	second error
	// ]
}
