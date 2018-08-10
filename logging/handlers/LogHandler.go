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
	Name() string
	SetName(name string)

	Debug(message interface{})

	Info(message interface{})

	Warn(message interface{})

	Error(message interface{})

	LogAtLevel(logLevel LogLevel, message interface{})

	Initialise()

	BeingDiscarded(logLevel LogLevel) bool

	SetDestinations(*LogLevelDestinations)
	Destinations() *LogLevelDestinations

	SetFormatter(formatter LogFormatter)
	Formatter() LogFormatter

	SupportsLogLevel(logLevel LogLevel) bool
	Override(logLevel LogLevel, destination LogDestination)
}

// LogHandlerBase is a base struct that implements default behaviour that matches the LogHandler interface
type LogHandlerBase struct {
	name         string
	destinations *LogLevelDestinations
	formatter    LogFormatter
}

// SetName allows a human-friendly name to be assigned to the loghandler to make it easier to configure
func (handlerBase *LogHandlerBase) SetName(name string) {
	handlerBase.name = name
}

func (handlerBase *LogHandlerBase) Name() string {
	return handlerBase.name
}

// SetDestinations allows a pre-defined LogLevelDestinations instance to be assigned, and subsequently used for
// log destination stream resolution.
func (handlerBase *LogHandlerBase) SetDestinations(destinations *LogLevelDestinations) {
	handlerBase.destinations = destinations
}

func (handlerBase *LogHandlerBase) Destinations() *LogLevelDestinations {
	return handlerBase.destinations
}

// SetFormatter tells the LogHandlerBase to use the supplied formatter for preparing a given log entry for writing
// to its final LogLevelDestination
func (handlerBase *LogHandlerBase) SetFormatter(formatter LogFormatter) {
	handlerBase.formatter = formatter
}

func (handlerBase *LogHandlerBase) Formatter() LogFormatter {
	return handlerBase.formatter
}

func (handlerBase *LogHandlerBase) BeingDiscarded(logLevel LogLevel) bool {
	return handlerBase.destinations.Destinations[logLevel] == DISCARD
}

func (handlerBase *LogHandlerBase) SupportsLogLevel(logLevel LogLevel) bool {
	_, present := handlerBase.destinations.Destinations[logLevel]
	return present
}

func (handlerBase *LogHandlerBase) Override(logLevel LogLevel, destination LogDestination) {
	handlerBase.destinations.Override(logLevel, destination)
}

func toLogAttributes(message interface{}) LogAttributes {
	switch message.(type) {
	case string:
		return LogAttributes{NameValuePair{Name: MessageNameLabel, Value: message}}
	case error:
		messageAsError, _ := message.(error)
		return LogAttributes{NameValuePair{Name: "Error", Value: messageAsError.Error()}}
	case LogAttributes:
		messageAsAttributes, _ := message.(LogAttributes)
		return messageAsAttributes
	default:
		panic("should not get here")
	}
}
