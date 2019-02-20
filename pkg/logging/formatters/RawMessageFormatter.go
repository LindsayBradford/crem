// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package formatters defines observer formatters that take Attributes and convert them into observer-ready strings.
package formatters

import (
	"github.com/LindsayBradford/crem/pkg/attributes"
)

// The default label for a Attributes entry that is used for storing free-form messages.
const MessageNameLabel = "Message"

const MessageErrorLabel = "Error"
const MessageWarnLabel = "Warn"

// RawMessageFormatter implements the Formatter interface by ignoring all logAttributes attributes supplied except
// the 'Message, 'Error' and 'Warn' attribute, returning a "formatted" string exactly as per the attribute's value.
type RawMessageFormatter struct{}

func (formatter *RawMessageFormatter) Format(attributes attributes.Attributes) string {
	for _, attribute := range attributes {
		if isSupportedName(attribute.Name) {
			return attribute.Value.(string)
		}
	}
	return ""
}

func isSupportedName(name string) bool {
	switch name {
	case MessageNameLabel, MessageErrorLabel, MessageWarnLabel:
		return true
	}
	return false
}
