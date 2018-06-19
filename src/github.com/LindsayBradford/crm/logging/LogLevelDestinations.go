// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	"io"
	"io/ioutil"
	"os"
)

type LogDestination io.Writer

var (
	DISCARD LogDestination  = ioutil.Discard
	STDOUT LogDestination = os.Stdout
	STDERR LogDestination = os.Stderr
)

type LogLevelDestinations struct {
	destinations map[LogLevel] LogDestination
}

func (this *LogLevelDestinations) Initialise()  *LogLevelDestinations {
		this = new(LogLevelDestinations)

		this.destinations = map[LogLevel] LogDestination {
			DEBUG: DISCARD,
			INFO: STDOUT,
			WARN: STDOUT,
			ERROR: STDERR,
		}

		return this
}

func (this *LogLevelDestinations) WithOverride(logLevel LogLevel, destination LogDestination)  *LogLevelDestinations {
	this.Override(logLevel, destination)
	return this
}

func (this *LogLevelDestinations) Override(logLevel LogLevel, destination LogDestination) {
	this.destinations[logLevel] = destination
}