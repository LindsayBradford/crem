// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/strings"
)

type UnmarshalContext struct {
	ConfigKey          string
	ValidValues        []string
	TextToValidate     string
	AssignmentFunction func()
}

func ProcessUnmarshalContext(context UnmarshalContext) error {
	if valueIsInList(context.TextToValidate, context.ValidValues...) {
		context.AssignmentFunction()
		return nil
	}
	return GenerateErrorFromContext(context)
}

func valueIsInList(value string, list ...string) bool {
	for _, listEntry := range list {
		if value == listEntry {
			return true
		}
	}
	return false
}

func GenerateErrorFromContext(context UnmarshalContext) error {
	const errorTemplate = "invalid Value \"%v\" specified for key \"%s\"; should be one of: %s"
	return fmt.Errorf(errorTemplate, context.TextToValidate, context.ConfigKey, listToString(context.ValidValues...))
}

func listToString(list ...string) string {
	builder := strings.FluentBuilder{}
	needsComma := false
	for _, entry := range list {
		if needsComma {
			builder.Add(", ")
		}

		builder.Add("\"", entry, "\"")
		needsComma = true
	}
	return builder.String()
}
