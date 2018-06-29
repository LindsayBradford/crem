// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	"fmt"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
	"github.com/LindsayBradford/crm/strings"
	"io"
	"time"
)

// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type int

type BareBonesLogHandler struct {
	LogHandlerBase
}

func (this *BareBonesLogHandler) Initialise() {}

func (this *BareBonesLogHandler) WithFormatter(formatter LogFormatter) *BareBonesLogHandler {
	formatter.Initialise()
	this.formatter = formatter
	return this
}

func (this *BareBonesLogHandler) Debug(message string) {
	this.LogAtLevel(DEBUG, message)
}

func (this *BareBonesLogHandler) DebugWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(DEBUG, logAttributes)
}

func (this *BareBonesLogHandler) Info(message string) {
	this.LogAtLevel(INFO, message)
}

func (this *BareBonesLogHandler) InfoWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(INFO, logAttributes)
}

func (this *BareBonesLogHandler) Warn(message string) {
	this.LogAtLevel(WARN, message)
}

func (this *BareBonesLogHandler) WarnWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(WARN, logAttributes)
}

func (this *BareBonesLogHandler) Error(message string) {
	this.LogAtLevel(ERROR, message)
}

func (this *BareBonesLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(ERROR, logAttributes)
}

func (this *BareBonesLogHandler) ErrorWithError(err error) {
	logAttributes := LogAttributes{NameValuePair{"Error", fmt.Sprintf(err.Error())}}
	logAttributes = prependLogLevel(ERROR, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(ERROR, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) LogAtLevel(logLevel LogLevel, message string) {
	logAttributes := LogAttributes{NameValuePair{MESSAGE_LABEL, message}}
	logAttributes = prependLogLevel(logLevel, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(logLevel, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {
	logAttributes = prependLogLevel(logLevel, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(logLevel, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) writeString(logLevel LogLevel, text string) {
	var builder strings.FluentBuilder
	builder.Add(text, "\n")

	// TODO:  sync.Mutex per destination for when we are concurrent?
	io.WriteString(this.destinations.Destinations[logLevel], builder.String())
}

func prependTimestamp(oldSlice []NameValuePair) []NameValuePair {
	timeAsString := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	return append([]NameValuePair{{"Time", timeAsString}}, oldSlice...)
}

func prependLogLevel(logLevel LogLevel, oldSlice []NameValuePair) []NameValuePair {
	return append([]NameValuePair{{"LogLevel", string(logLevel)}}, oldSlice...)
}

func prepend(newValue NameValuePair, oldSlice []NameValuePair) []NameValuePair {
	return append([]NameValuePair{newValue}, oldSlice...)
}