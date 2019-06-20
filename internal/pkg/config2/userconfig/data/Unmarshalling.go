// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"fmt"

	"github.com/LindsayBradford/crem/pkg/strings"
)

type unmarshalContext struct {
	configKey          string
	validValues        []string
	textToValidate     string
	assignmentFunction func()
}

func processUnmarshalContext(context unmarshalContext) error {
	if valueIsInList(context.textToValidate, context.validValues...) {
		context.assignmentFunction()
		return nil
	}
	return generateErrorFromContext(context)
}

func valueIsInList(value string, list ...string) bool {
	for _, listEntry := range list {
		if value == listEntry {
			return true
		}
	}
	return false
}

func generateErrorFromContext(context unmarshalContext) error {
	const errorTemplate = "invalid value \"%v\" specified for key \"%s\"; should be one of: %s"
	return fmt.Errorf(errorTemplate, context.textToValidate, context.configKey, listToString(context.validValues...))
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
