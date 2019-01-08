// Copyright (c) 2019 Australian Rivers Institute.

package formatters

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/logging"
)

func ExampleNullFormatter_Format() {
	exampleAttributes := logging.Attributes{
		{Name: "NoMatter", Value: "Ignored anyway"},
	}

	exampleFormatter := new(NullFormatter)

	exampleJson := exampleFormatter.Format(exampleAttributes)
	fmt.Print(exampleJson)

	// Output:No formatter specified. Using the NullFormatter.
}
