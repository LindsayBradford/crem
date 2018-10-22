// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	"io"
	"time"

	. "github.com/LindsayBradford/crem/logging/formatters"
	. "github.com/LindsayBradford/crem/logging/shared"
	"github.com/LindsayBradford/crem/strings"
)

// value is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type int

type BareBonesLogHandler struct {
	LogHandlerBase
}

func (handler *BareBonesLogHandler) Initialise() {}

func (handler *BareBonesLogHandler) WithFormatter(formatter LogFormatter) *BareBonesLogHandler {
	formatter.Initialise()
	handler.formatter = formatter
	return handler
}

func (handler *BareBonesLogHandler) Debug(message interface{}) {
	handler.LogAtLevel(DEBUG, message)
}

func (handler *BareBonesLogHandler) Info(message interface{}) {
	handler.LogAtLevel(INFO, message)
}

func (handler *BareBonesLogHandler) Warn(message interface{}) {
	handler.LogAtLevel(WARN, message)
}

func (handler *BareBonesLogHandler) Error(message interface{}) {
	handler.LogAtLevel(ERROR, message)
}

func (handler *BareBonesLogHandler) LogAtLevel(logLevel LogLevel, message interface{}) {
	messageAttributes := toLogAttributes(message)
	messageAttributes = prependLogLevel(logLevel, messageAttributes)
	messageAttributes = prependTimestamp(messageAttributes)
	handler.writeString(logLevel, handler.formatter.Format(messageAttributes))
}

func (handler *BareBonesLogHandler) writeString(logLevel LogLevel, text string) {
	var builder strings.FluentBuilder
	builder.Add(text, "\n")
	io.WriteString(handler.deriveDestination(logLevel), builder.String())
}

func (handler *BareBonesLogHandler) deriveDestination(logLevel LogLevel) LogDestination {
	return handler.destinations.Destinations[logLevel]
}

func prependTimestamp(oldSlice []NameValuePair) []NameValuePair {
	timeAsString := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	newPair := NameValuePair{Name: "Time", Value: timeAsString}
	return append([]NameValuePair{newPair}, oldSlice...)
}

func prependLogLevel(logLevel LogLevel, oldSlice []NameValuePair) []NameValuePair {
	newPair := NameValuePair{Name: "LogLevel", Value: string(logLevel)}
	return append([]NameValuePair{newPair}, oldSlice...)
}
