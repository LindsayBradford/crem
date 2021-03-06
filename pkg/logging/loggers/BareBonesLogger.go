// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package loggers

import (
	"io"
	"time"

	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// value is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type int

type BareBonesLogger struct {
	LoggerBase
}

func (bbl *BareBonesLogger) Initialise() {}

func (bbl *BareBonesLogger) WithFormatter(formatter logging.Formatter) *BareBonesLogger {
	bbl.formatter = formatter
	return bbl
}

func (bbl *BareBonesLogger) Debug(message interface{}) {
	bbl.LogAtLevel(logging.DEBUG, message)
}

func (bbl *BareBonesLogger) Info(message interface{}) {
	bbl.LogAtLevel(logging.INFO, message)
}

func (bbl *BareBonesLogger) Warn(message interface{}) {
	bbl.LogAtLevel(logging.WARN, message)
}

func (bbl *BareBonesLogger) Error(message interface{}) {
	bbl.LogAtLevel(logging.ERROR, message)
}

func (bbl *BareBonesLogger) LogAtLevel(logLevel logging.Level, message interface{}) {
	messageAttributes := toLogAttributes(message)
	messageAttributes = prependLogLevel(logLevel, messageAttributes)
	messageAttributes = prependTimestamp(messageAttributes)
	bbl.writeString(logLevel, bbl.formatter.Format(messageAttributes))
}

func (bbl *BareBonesLogger) LogAtLevelWithAttributes(logLevel logging.Level, logAttributes attributes.Attributes) {
	logAttributes = prependLogLevel(logLevel, logAttributes)
	logAttributes = prependTimestamp(logAttributes)
	bbl.writeString(logLevel, bbl.formatter.Format(logAttributes))
}

func (bbl *BareBonesLogger) writeString(logLevel logging.Level, text string) {
	var builder strings.FluentBuilder
	builder.Add(text, "\n")
	io.WriteString(bbl.deriveDestination(logLevel), builder.String())
}

func (bbl *BareBonesLogger) deriveDestination(logLevel logging.Level) logging.Destination {
	return bbl.destinations.Destinations[logLevel]
}

func prependTimestamp(oldSlice []attributes.NameValuePair) []attributes.NameValuePair {
	timeAsString := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	newPair := attributes.NameValuePair{Name: "Time", Value: timeAsString}
	return append([]attributes.NameValuePair{newPair}, oldSlice...)
}

func prependLogLevel(logLevel logging.Level, oldSlice []attributes.NameValuePair) []attributes.NameValuePair {
	newPair := attributes.NameValuePair{Name: "Level", Value: string(logLevel)}
	return append([]attributes.NameValuePair{newPair}, oldSlice...)
}
