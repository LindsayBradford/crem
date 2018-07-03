// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// handlers package defines handlers responsible for formatting log entries (delegated to LogFormatter) and delivering
// the formatted entries to whatever log destinations are needed.
package handlers

import (
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
)

// LogHandler defines an interface for the handling of logging. It sets out methods for logging at the various supported
// LogLevels of either a free-form string (traditional), or LogAttributes (for machine-friendly logging). It delegates
// formatting to a LogFormatter, and resolution of log destination streams to LogLevelDestinations.
type LogHandler interface {
	Debug(message string)
	DebugWithAttributes(logAttributes LogAttributes)

	Info(message string)
	InfoWithAttributes(logAttributes LogAttributes)

	Warn(message string)
	WarnWithAttributes(logAttributes LogAttributes)

	Error(message string)
	ErrorWithAttributes(logAttributes LogAttributes)
	ErrorWithError(err error)

	LogAtLevel(logLevel LogLevel, message string)
	LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes)

	Initialise()

	BeingDiscarded(logLevel LogLevel) bool

	SetDestinations(*LogLevelDestinations)
	Destinations() *LogLevelDestinations

	SetFormatter(formatter LogFormatter)
	Formatter() LogFormatter
}

// LogHandlerBase is a base struct that implements default behaviour that matches the LogHandler interface
type LogHandlerBase struct {
	destinations *LogLevelDestinations
	formatter    LogFormatter
}

// SetDestinations allows a pre-defined LogLevelDestinations instance to be assigned, and subsequently used for
// log destination stream resolution.
func (this *LogHandlerBase) SetDestinations(destinations *LogLevelDestinations) {
	this.destinations = destinations
}

func (this *LogHandlerBase) Destinations() *LogLevelDestinations {
	return this.destinations
}

// SetFormatter tells the LogHandlerBase to use the supplied formatter for preparing a given log entry for writing
// to its final LogLevelDestination
func (this *LogHandlerBase) SetFormatter(formatter LogFormatter) {
	this.formatter = formatter
}

func (this *LogHandlerBase) Formatter() LogFormatter {
	return this.formatter
}

func (this *LogHandlerBase)  BeingDiscarded(logLevel LogLevel) bool {
	return this.destinations.Destinations[logLevel] == DISCARD
}

type NullLogHandler struct {}

func (this *NullLogHandler) Debug(message string) {}
func (this *NullLogHandler) DebugWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) Info(message string) {}
func (this *NullLogHandler) InfoWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) Warn(message string) {}
func (this *NullLogHandler) WarnWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) Error(message string) {}
func (this *NullLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) ErrorWithError(err error) {}
func (this *NullLogHandler) LogAtLevel(logLevel LogLevel, message string) {}
func (this *NullLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {}
func (this *NullLogHandler) Initialise() {}
func (this *NullLogHandler) SetDestinations(*LogLevelDestinations) {}
func (this *NullLogHandler) Destinations() *LogLevelDestinations {return nil}
func (this *NullLogHandler) SetFormatter(formatter LogFormatter) {}
func (this *NullLogHandler) Formatter() LogFormatter{return &NullFormatter{}}
func (this *NullLogHandler) BeingDiscarded(logLevel LogLevel) bool { return true }