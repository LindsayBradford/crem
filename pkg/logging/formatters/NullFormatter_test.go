// Copyright (c) 2019 Australian Rivers Institute.

package formatters

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/attributes"
)

func ExampleNullFormatter_Format() {
	exampleAttributes := attributes.Attributes{
		{Name: "NoMatter", Value: "Ignored anyway"},
	}

	exampleFormatter := new(NullFormatter)

	exampleJson := exampleFormatter.Format(exampleAttributes)
	fmt.Print(exampleJson)

	// Output:No formatter specified. Using the NullFormatter.
}
