// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	"io"
	"time"

	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
	"github.com/LindsayBradford/crm/strings"
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

func (handler *BareBonesLogHandler) Debug(message string) {
	handler.LogAtLevel(DEBUG, message)
}

func (handler *BareBonesLogHandler) DebugWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(DEBUG, logAttributes)
}

func (handler *BareBonesLogHandler) Info(message string) {
	handler.LogAtLevel(INFO, message)
}

func (handler *BareBonesLogHandler) InfoWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(INFO, logAttributes)
}

func (handler *BareBonesLogHandler) Warn(message string) {
	handler.LogAtLevel(WARN, message)
}

func (handler *BareBonesLogHandler) WarnWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(WARN, logAttributes)
}

func (handler *BareBonesLogHandler) Error(message string) {
	handler.LogAtLevel(ERROR, message)
}

func (handler *BareBonesLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(ERROR, logAttributes)
}

func (handler *BareBonesLogHandler) ErrorWithError(err error) {
	logAttributes := LogAttributes{NameValuePair{Name: "Error", Value: err.Error()}}
	logAttributes = prependLogLevel(ERROR, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	handler.writeString(ERROR, handler.formatter.Format(logAttributes))
}

func (handler *BareBonesLogHandler) LogAtLevel(logLevel LogLevel, message string) {
	logAttributes := LogAttributes{NameValuePair{Name: MessageNameLabel, Value: message}}
	logAttributes = prependLogLevel(logLevel, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	handler.writeString(logLevel, handler.formatter.Format(logAttributes))
}

func (handler *BareBonesLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {
	logAttributes = prependLogLevel(logLevel, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	handler.writeString(logLevel, handler.formatter.Format(logAttributes))
}

func (handler *BareBonesLogHandler) writeString(logLevel LogLevel, text string) {
	var builder strings.FluentBuilder
	builder.Add(text, "\n")
	io.WriteString(handler.destinations.Destinations[logLevel], builder.String())
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
