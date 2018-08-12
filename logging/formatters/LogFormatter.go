// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package formatters defines logging formatters that take LogAttributes and convert them into logging-ready strings.
package formatters

import . "github.com/LindsayBradford/crm/logging/shared"

// LogFormatter describes an interface for the formatting of LogAttributes into some logging-ready string.
// Instances of LogHandler are expected to delegate any formatting of the supplied attributes to a LogFormatter.
type LogFormatter interface {
	// Initialises any necessary solution the formatter requires prior to being used.
	Initialise()

	// Format converts the supplied attributes into a representative 'logging ready' string.
	Format(attributes LogAttributes) string
}

// nullFormatMessage is the message supplied if the NullFormatter is left as the Formatter of a LogHandler.
const nullFormatMessage = "No formatter specified. Using the NullFormatter."

// NullFormatter implements a 'null object' placeholder formatter that is supplied by default if one is not specified.
// It returns a static message as per nullFormatMessage as a reminder that a proper formatter must be supplied
// for the logging to do anything meaningful.
type NullFormatter struct{}

func (formatter *NullFormatter) Initialise() {}

func (formatter *NullFormatter) Format(attributes LogAttributes) string {
	return nullFormatMessage
}

// The default label for a LogAttributes entry that is used for storing free-form messages.
const MessageNameLabel = "Message"

const MessageErrorLabel = "Error"

// RawMessageFormatter implements the Formatter interface by ignoring all logAttributes attributes supplied except
// the 'message' (MessageNameLabel) attribute, and returns a "formatted" string exactly as was supplied
// in that attribute.
type RawMessageFormatter struct {
	messageLabel string
	errorLabel   string
}

func (formatter *RawMessageFormatter) Initialise() {
	formatter.messageLabel = MessageNameLabel
	formatter.errorLabel = MessageErrorLabel
}

func (formatter *RawMessageFormatter) Format(attributes LogAttributes) string {
	for _, attribute := range attributes {
		if attribute.Name == formatter.messageLabel || attribute.Name == formatter.errorLabel {
			return attribute.Value.(string)
		}
	}
	return ""
}
