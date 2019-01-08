// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package formatters

import (
	"fmt"
	"testing"

	"github.com/LindsayBradford/crem/pkg/logging"
	. "github.com/onsi/gomega"
)

func ExampleRawMessageFormatter_Format() {
	expectedMessage := "here is an expected message"
	exampleAttributes := logging.Attributes{
		{Name: "Message", Value: expectedMessage},
		{Name: "NotAMessage", Value: "who cares?"},
	}

	exampleFormatter := new(RawMessageFormatter)
	exampleFormatter.Initialise()

	exampleMsg := exampleFormatter.Format(exampleAttributes)
	fmt.Print(exampleMsg)

	// Output:here is an expected message
}

func TestRawMessageFormatter_FormatError(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedMessage := "here is an error message"
	attribsUnderTest := logging.Attributes{
		{Name: string(logging.ERROR), Value: expectedMessage},
		{Name: "NotAMessage", Value: "who cares?"},
	}

	exampleFormatter := new(RawMessageFormatter)
	exampleFormatter.Initialise()

	actualMessage := exampleFormatter.Format(attribsUnderTest)
	g.Expect(actualMessage).To(Equal(expectedMessage))
}

func TestRawMessageFormatter_FormatWarn(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedMessage := "here is a warning message"
	attribsUnderTest := logging.Attributes{
		{Name: string(logging.WARN), Value: expectedMessage},
		{Name: "NotAMessage", Value: "who cares?"},
	}

	exampleFormatter := new(RawMessageFormatter)
	exampleFormatter.Initialise()

	actualMessage := exampleFormatter.Format(attribsUnderTest)
	g.Expect(actualMessage).To(Equal(expectedMessage))
}

func TestRawMessageFormatter_FormatEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedMessage := ""
	attribsUnderTest := logging.Attributes{
		{Name: "NotAMessage", Value: "who cares?"},
	}

	exampleFormatter := new(RawMessageFormatter)
	exampleFormatter.Initialise()

	actualMessage := exampleFormatter.Format(attribsUnderTest)
	g.Expect(actualMessage).To(Equal(expectedMessage))
}
