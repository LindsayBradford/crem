// Copyright (c) 2018 Australian Rivers Institute.

// loggers package defines loggers responsible for formatters log entries (delegated to Formatter) and delivering
// the formatted entries to whatever log destinations are needed.
package logging

import (
	"errors"

	"github.com/LindsayBradford/crem/pkg/name"
)

// Logger defines an interface for the handling of observer. It sets out methods for observer at the various supported
// LogLevels of either a free-form string (traditional), or Attributes (for machine-friendly observer). It delegates
// formatters to a Formatter, and resolution of log destination streams to Destinations.
type Logger interface {
	name.Nameable

	Debug(message interface{})
	Info(message interface{})
	Warn(message interface{})
	Error(message interface{})

	LogAtLevel(logLevel Level, message interface{})
	LogAtLevelWithAttributes(logLevel Level, attributes Attributes)

	Initialise()

	BeingDiscarded(logLevel Level) bool

	SetDestinations(*Destinations)
	Destinations() *Destinations

	SetFormatter(formatter Formatter)
	Formatter() Formatter

	SupportsLogLevel(logLevel Level) bool
	Override(logLevel Level, destination Destination)
}

// Container defines an interface for users wishing to embed a Logger.
type Container interface {
	SetLogHandler(logger Logger) error
	LogHandler() Logger
}

// ContainedLogger is a struct offering a default implementation of Container
type ContainedLogger struct {
	logHandler Logger
}

func (cl *ContainedLogger) SetLogHandler(logHandler Logger) error {
	if logHandler == nil {
		return errors.New("invalid attempt to set log handler to nil value")
	}
	cl.logHandler = logHandler
	return nil
}

func (cl *ContainedLogger) LogHandler() Logger {
	return cl.logHandler
}
