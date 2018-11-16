// Copyright (c) 2018 Australian Rivers Institute.

package logging

// LogAtLevel is a string identifier for a 'typical' set of log levels we might want to report against, ranging from
// DEBUG at the least critical/most noisy of log levels, to ERROR as the most critical/least noisy.
type Level string

const (
	DEBUG Level = "Debug"
	INFO  Level = "Info"
	WARN  Level = "Warn"
	ERROR Level = "Error"
)
