// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package formatters

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/logging"
)

type nvpStringer struct{}

func (s *nvpStringer) String() string {
	return "stringerValue"
}

func ExampleNameValuePairFormatter_Format() {
	exampleAttributes := logging.Attributes{
		{Name: "One", Value: "valueOne"},
		{Name: "Two", Value: 42},
		{Name: "Three", Value: uint64(0)},
		{Name: "Four", Value: 42.42},
		{Name: "Five", Value: new(nvpStringer)},
		{Name: "Six", Value: true},
	}

	exampleFormatter := new(NameValuePairFormatter)
	exampleFormatter.Initialise()

	exampleNvp := exampleFormatter.Format(exampleAttributes)
	fmt.Print(exampleNvp)

	// Output: One="valueOne", Two=42, Three=0, Four=42.42, Five="stringerValue", Six=true
}
