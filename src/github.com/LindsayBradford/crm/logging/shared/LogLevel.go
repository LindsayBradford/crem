// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

// LogLevel is a string identifier for a 'typical' set of log levels we might want to report against, ranging from
// DEBUG at the least critical/most noisy of log levels, to ERROR as the most critical/least noisy.
type LogLevel string

const (
	DEBUG LogLevel = "Debug"
	INFO LogLevel = "Info"
	WARN LogLevel = "Warn"
	ERROR LogLevel = "Error"
)