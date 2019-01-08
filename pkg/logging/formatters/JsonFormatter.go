// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package formatters

import (
	"fmt"
	"strconv"

	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// JsonFormatter formats a Attributes array into an equivalent JSON encoding.
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

		builder.Add("\"", attribute.Name, "\": ", jsonValueToString(attribute.Value))
	}

	builder.Add("}")
	return builder.String()
}

func jsonValueToString(value interface{}) string {
	switch value.(type) {
	case bool:
		return strconv.FormatBool(value.(bool))
	case int:
		return strconv.Itoa(value.(int))
	case uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case float64:
		return strconv.FormatFloat(value.(float64), 'g', -1, 64)
	case string:
		return "\"" + value.(string) + "\""
	case fmt.Stringer:
		typeConvertedValue := value.(fmt.Stringer)
		return "\"" + typeConvertedValue.String() + "\""
	}
	return "null"
}
