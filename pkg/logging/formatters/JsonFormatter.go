// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package formatters

import (
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// JsonFormatter formats a Attributes array into an equivalent JSON encoding.
// TODO: Supply example encoding.
type JsonFormatter struct{}

func (formatter *JsonFormatter) Initialise() {}

func (formatter *JsonFormatter) Format(attributes logging.Attributes) string {
	var builder strings.FluentBuilder

	builder.Add("{")
	needsComma := false

	for _, attribute := range attributes {
		if !needsComma {
			needsComma = true
		} else {
			builder.Add(", ")
		}
		builder.Add("\"", attribute.Name, "\": \"", attribute.Value.(string), "\"")
	}

	builder.Add("}")
	return builder.String()
}