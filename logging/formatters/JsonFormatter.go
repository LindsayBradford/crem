// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package formatters

import (
	. "github.com/LindsayBradford/crm/logging/shared"
	"github.com/LindsayBradford/crm/strings"
)

// JsonFormatter formats a LogAttributes array into an equivalent JSON encoding.
// TODO: Supply example encoding.
type JsonFormatter struct{}

func (formatter *JsonFormatter) Initialise() {}

func (formatter *JsonFormatter) Format(attributes LogAttributes) string {
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
