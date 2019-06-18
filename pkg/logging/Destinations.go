// Copyright (c) 2018 Australian Rivers Institute.

package logging

import (
	"io"
	"io/ioutil"
	"os"
)

// Destination captures the output stream that a log should be written to.
type Destination io.Writer

// Three default Destination entries are provided by the package. DISCARD will cause the log entries to be written
// nowhere. STDOUT to standard console output, and STDERR for standard error console output.
var (
	DISCARD Destination = ioutil.Discard
	STDOUT  Destination = os.Stdout
	STDERR  Destination = os.Stderr
)

// Destinations is a mapping of LogAtLevel values to Destination values.
type Destinations struct {
	Destinations map[Level]Destination
}

// CreateEmpty creates and returns a Destinations instance with a default Destinations map.
// Specifically, DEBUG is discarded, INFO and WARN are delivered to STDOUT, and ERROR to STDERR.
func (d *Destinations) Initialise() *Destinations {
	d.Destinations = map[Level]Destination{
		DEBUG: DISCARD,
		INFO:  STDOUT,
		WARN:  STDOUT,
		ERROR: STDERR,
	}

	return d
}

// WithOverride is a fluent method for overriding the existing Destinations map entry for logLevel to instead
// point to the new destination supplied.
func (d *Destinations) WithOverride(logLevel Level, destination Destination) *Destinations {
	d.Override(logLevel, destination)
	return d
}

// Override remaps the given LogLevelDestination's logLevel mapping to the new destination supplied.
func (d *Destinations) Override(logLevel Level, destination Destination) {
	d.Destinations[logLevel] = destination
}

// Override remaps the given LogLevelDestination's logLevel mapping to the new destination supplied.
func (d *Destinations) SupportsLogLevel(logLevel Level) bool {
	_, present := d.Destinations[logLevel]
	return present
}
