// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package formatters

import (
	"fmt"
	"strconv"

	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/strings"
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
			builder.Add(", ")
		}
		builder.Add(attribute.Name, "=", nvpValueToString(attribute.Value))
	}
	return builder.String()
}

func nvpValueToString(value interface{}) string {
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
