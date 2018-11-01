// Copyright (c) 2018 Australian Rivers Institute.

package loggers

import (
	"github.com/LindsayBradford/crem/logging"
	"github.com/LindsayBradford/crem/logging/formatters"
)

// LoggerBase is a base struct that implements default behaviour that matches the Logger interface
type LoggerBase struct {
	name         string
	destinations *logging.Destinations
	formatter    logging.Formatter
}

// SetName allows a human-friendly name to be assigned to the loghandler to make it easier to configure
func (lb *LoggerBase) SetName(name string) {
	lb.name = name
}

func (lb *LoggerBase) Name() string {
	return lb.name
}

// SetDestinations allows a pre-defined Destinations instance to be assigned, and subsequently used for
// log destination stream resolution.
func (lb *LoggerBase) SetDestinations(destinations *logging.Destinations) {
	lb.destinations = destinations
}

func (lb *LoggerBase) Destinations() *logging.Destinations {
	return lb.destinations
}

// SetFormatter tells the LoggerBase to use the supplied formatter for preparing a given log entry for writing
// to its final LogLevelDestination
func (lb *LoggerBase) SetFormatter(formatter logging.Formatter) {
	lb.formatter = formatter
}

func (lb *LoggerBase) Formatter() logging.Formatter {
	return lb.formatter
}

func (lb *LoggerBase) BeingDiscarded(logLevel logging.Level) bool {
	return lb.destinations.Destinations[logLevel] == logging.DISCARD
}

func (lb *LoggerBase) SupportsLogLevel(logLevel logging.Level) bool {
	_, present := lb.destinations.Destinations[logLevel]
	return present
}

func (lb *LoggerBase) Override(logLevel logging.Level, destination logging.Destination) {
	lb.destinations.Override(logLevel, destination)
}

func toLogAttributes(message interface{}) logging.Attributes {
	switch message.(type) {
	case string:
		return logging.Attributes{logging.NameValuePair{Name: formatters.MessageNameLabel, Value: message}}
	case error:
		messageAsError, _ := message.(error)
		return logging.Attributes{logging.NameValuePair{Name: formatters.MessageErrorLabel, Value: messageAsError.Error()}}
	case logging.Attributes:
		messageAsAttributes, _ := message.(logging.Attributes)
		return messageAsAttributes
	default:
		panic("should not get here")
	}
}
