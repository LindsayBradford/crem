// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package formatters defines observer formatters that take LogAttributes and convert them into observer-ready strings.
package formatters

import . "github.com/LindsayBradford/crem/logging/shared"

// LogFormatter describes an interface for the formatting of LogAttributes into some observer-ready string.
// Instances of LogHandler are expected to delegate any formatting of the supplied attributes to a LogFormatter.
type LogFormatter interface {
	// Initialises any necessary explorer the formatter requires prior to being used.
	Initialise()

	// Format converts the supplied attributes into a representative 'observer ready' string.
	Format(attributes LogAttributes) string
}

// nullFormatMessage is the message supplied if the NullFormatter is left as the Formatter of a LogHandler.
const nullFormatMessage = "No formatter specified. Using the NullFormatter."

// NullFormatter implements a 'null object' placeholder formatter that is supplied by default if one is not specified.
// It returns a static message as per nullFormatMessage as a reminder that a proper formatter must be supplied
// for the observer to do anything meaningful.
type NullFormatter struct{}

func (formatter *NullFormatter) Initialise() {}

func (formatter *NullFormatter) Format(attributes LogAttributes) string {
	return nullFormatMessage
}

// The default label for a LogAttributes entry that is used for storing free-form messages.
const MessageNameLabel = "Message"

const MessageErrorLabel = "Error"
const MessageWarnLabel = "Warn"

// RawMessageFormatter implements the Formatter interface by ignoring all logAttributes attributes supplied except
// the 'Message, 'Error' and 'Warn' attribute, returning a "formatted" string exactly as per the attribute's value.
type RawMessageFormatter struct{}

func (formatter *RawMessageFormatter) Initialise() {}

func (formatter *RawMessageFormatter) Format(attributes LogAttributes) string {
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
