// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"io"
	"io/ioutil"
	"os"
)

// LogDestination captures the output stream that a log should be written to.
type LogDestination io.Writer

// Three default LogDestination entries are provided by the package. DISCARD will cause the log entries to be written
// nowhere. STDOUT to standard console output, and STDERR for standard error console output.
var (
	DISCARD LogDestination = ioutil.Discard
	STDOUT  LogDestination = os.Stdout
	STDERR  LogDestination = os.Stderr
)

// LogLevelDestinations is a mapping of LogAtLevel values to LogDestination values.
type LogLevelDestinations struct {
	Destinations map[LogLevel]LogDestination
}

// Initialise creates and returns a LogLevelDestinations instance with a default LogLevelDestinations map.
// Specifically, DEBUG is discarded, INFO and WARN are delivered to STDOUT, and ERROR to STDERR.
func (lld *LogLevelDestinations) Initialise() *LogLevelDestinations {
	lld.Destinations = map[LogLevel]LogDestination{
		DEBUG: DISCARD,
		INFO:  STDOUT,
		WARN:  STDOUT,
		ERROR: STDERR,
	}

	return lld
}

// WithOverride is a fluent method for overriding the existing LogLevelDestinations map entry for logLevel to instead
// point to the new destination supplied.
func (lld *LogLevelDestinations) WithOverride(logLevel LogLevel, destination LogDestination) *LogLevelDestinations {
	lld.Override(logLevel, destination)
	return lld
}

// Override remaps the given LogLevelDestination's logLevel mapping to the new destination supplied.
func (lld *LogLevelDestinations) Override(logLevel LogLevel, destination LogDestination) {
	lld.Destinations[logLevel] = destination
}

// Override remaps the given LogLevelDestination's logLevel mapping to the new destination supplied.
func (lld *LogLevelDestinations) SupportsLogLevel(logLevel LogLevel) bool {
	_, present := lld.Destinations[logLevel]
	return present
}
