// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	"fmt"
	. "github.com/LindsayBradford/crm/logging/shared"
	. "github.com/LindsayBradford/crm/logging/formatters"
	"io"
	"github.com/LindsayBradford/crm/strings"
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
	logAttributes := LogAttributes{ NameValuePair{ MESSAGE_LABEL, message }}
	logAttributes = prependLogLevel(DEBUG, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(DEBUG, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) DebugWithAttributes(logAttributes LogAttributes) {
	logAttributes = prependLogLevel(DEBUG, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(DEBUG, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) Info(message string) {
	logAttributes := LogAttributes{ NameValuePair{ MESSAGE_LABEL, message }}
	logAttributes = prependLogLevel(INFO, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(INFO, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) InfoWithAttributes(logAttributes LogAttributes) {
	logAttributes = prependLogLevel(INFO, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(INFO, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) Warn(message string) {
	logAttributes := LogAttributes{ NameValuePair{ MESSAGE_LABEL, message }}
	logAttributes = prependLogLevel(WARN, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(WARN, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) WarnWithAttributes(logAttributes LogAttributes) {
	logAttributes = prependLogLevel(WARN, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(WARN, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) Error(message string) {
	logAttributes := LogAttributes{ NameValuePair{ MESSAGE_LABEL, message }}
	logAttributes = prependLogLevel(ERROR, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(ERROR, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {
	logAttributes = prependLogLevel(ERROR, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(ERROR, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) ErrorWithError(err error) {
	logAttributes := LogAttributes{ NameValuePair{ "Error", fmt.Sprintf(err.Error())}}
	logAttributes = prependLogLevel(WARN, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	this.writeString(ERROR, this.formatter.Format(logAttributes))
}

func (this *BareBonesLogHandler) writeString(logLevel LogLevel, text string) {
	var builder strings.FluentBuilder
	builder.Add(text, "\n")

	// TODO:  THis probably needs a semaphore once we go concurrent.
	io.WriteString(this.destinations.Destinations[logLevel], builder.String())
}

func prependTimestamp(oldSlice []NameValuePair) []NameValuePair {
	timeAsString := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	return append([]NameValuePair{{"Time", timeAsString}}, oldSlice...)
}

func prependLogLevel(logLevel LogLevel, oldSlice []NameValuePair) []NameValuePair {
	return append([]NameValuePair{{"LogLevel", string(logLevel)}}, oldSlice...)
}

func prepend (newValue NameValuePair, oldSlice []NameValuePair) []NameValuePair {
	return append([]NameValuePair{newValue}, oldSlice...)
}