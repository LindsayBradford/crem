// Copyright (c) 2019 Australian Rivers Institute.

package formatters

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/attributes"
)

type jsonStringer struct{}

func (s *jsonStringer) String() string {
	return "stringerValue"
}

func ExampleJsonFormatter_Format() {
	exampleAttributes := attributes.Attributes{
		{Name: "One", Value: "valueOne"},
		{Name: "Two", Value: 42},
		{Name: "Three", Value: uint64(0)},
		{Name: "Four", Value: 42.42},
		{Name: "Five", Value: new(jsonStringer)},
		{Name: "Six", Value: true},
		{Name: "Seven", Value: 7777.777777},
	}

	exampleFormatter := new(JsonFormatter)

	exampleJson := exampleFormatter.Format(exampleAttributes)
	fmt.Print(exampleJson)

	// Output: {"One": "valueOne", "Two": 42, "Three": 0, "Four": 42.42, "Five": "stringerValue", "Six": true, "Seven": 7,777.777777}
}
