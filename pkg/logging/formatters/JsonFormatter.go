// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package formatters

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/strings"
)

const (
	escapedQuote = "\""
	nullString   = "null"
	comma        = ", "
	colon        = ": "
	openBracket  = "{"
	closeBracket = "}"
)

// JsonFormatter formats a Attributes array into an equivalent JSON encoding.
type JsonFormatter struct{}

func (formatter *JsonFormatter) Format(attributes attributes.Attributes) string {
	var builder strings.FluentBuilder

	builder.Add(openBracket)

	needsComma := false
	for _, attribute := range attributes {
		if !needsComma {
			needsComma = true
		} else {
			builder.Add(comma)
		}

		builder.
			Add(escapedQuote, attribute.Name, escapedQuote).
			Add(colon, jsonValueToString(attribute.Value))
	}

	builder.Add(closeBracket)

	return builder.String()
}

func jsonValueToString(value interface{}) string {
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
