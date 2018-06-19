// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

type LogLevel string

const (
	DEBUG LogLevel = "Debug"
	INFO LogLevel = "Info"
	WARN LogLevel = "WarnWithAttributes"
	ERROR LogLevel = "Error"
)

type NameValuePair struct {
	Name  string
	Value interface{}
}

type LogAttributes []NameValuePair

type LogAttributeFormatter interface {
	Initialise()
	Format(attributes LogAttributes) string
}

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

	SetDestinations(*LogLevelDestinations)
	Destinations() *LogLevelDestinations

	SetFormatter(formatter LogAttributeFormatter)
	Formatter() LogAttributeFormatter
}

type LogHandlerBase struct {
	destinations *LogLevelDestinations
	formatter LogAttributeFormatter
}

func (this *LogHandlerBase) SetDestinations(destinations *LogLevelDestinations) {
	this.destinations = destinations
}

func (this *LogHandlerBase) Destinations() *LogLevelDestinations {
	return this.destinations
}

func (this *LogHandlerBase) SetFormatter(formatter LogAttributeFormatter) {
	this.formatter = formatter
}

func (this *LogHandlerBase) Formatter() LogAttributeFormatter {
	return this.formatter
}
