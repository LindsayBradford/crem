// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

// loggers package defines loggers responsible for formatters log entries (delegated to Formatter) and delivering
// the formatted entries to whatever log destinations are needed.
package logging

// Logger defines an interface for the handling of observer. It sets out methods for observer at the various supported
// LogLevels of either a free-form string (traditional), or Attributes (for machine-friendly observer). It delegates
// formatters to a Formatter, and resolution of log destination streams to Destinations.
type Logger interface {
	Name() string
	SetName(name string)

	Debug(message interface{})

	Info(message interface{})

	Warn(message interface{})

	Error(message interface{})

	LogAtLevel(logLevel Level, message interface{})

	Initialise()

	BeingDiscarded(logLevel Level) bool

	SetDestinations(*Destinations)
	Destinations() *Destinations

	SetFormatter(formatter Formatter)
	Formatter() Formatter

	SupportsLogLevel(logLevel Level) bool
	Override(logLevel Level, destination Destination)
}
