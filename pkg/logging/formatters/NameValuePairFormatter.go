// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package formatters

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/strings"
)

const (
	equals = "="
)

// NameValuePairFormatter formats a Attributes array into a string of comma-separated name-value pairs.
type NameValuePairFormatter struct{}

func (formatter *NameValuePairFormatter) Format(attributes attributes.Attributes) string {
	var builder strings.FluentBuilder

	needsComma := false

	for _, attribute := range attributes {
		if !needsComma {
			needsComma = true
		} else {
			builder.Add(comma)
		}
		builder.Add(attribute.Name, equals, nvpValueToString(attribute.Value))
	}
	return builder.String()
}

func nvpValueToString(value interface{}) string {
	if r := recover(); r != nil {
		return nullString
	}

	switch value.(type) {
	case string, fmt.Stringer:
		return escapedQuote + strings.Convert(value) + escapedQuote
	default:
		return strings.Convert(value)
	}

	return nullString
}
