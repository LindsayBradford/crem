// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package formatters

import (
	"github.com/LindsayBradford/crem/logging"
	"github.com/LindsayBradford/crem/strings"
)

// NameValuePairFormatter formats a Attributes array into a string of comma-separated name-value pairs.
// TODO: Supply example encoding.
type NameValuePairFormatter struct{}

func (formatter *NameValuePairFormatter) Initialise() {}

func (formatter *NameValuePairFormatter) Format(attributes logging.Attributes) string {
	var builder strings.FluentBuilder

	needsComma := false

	for _, attribute := range attributes {
		if !needsComma {
			needsComma = true
		} else {
			builder.Add(", ")
		}
		builder.Add(attribute.Name, "=\"", attribute.Value.(string), "\"")
	}
	return builder.String()
}
