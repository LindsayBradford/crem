// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	"errors"
	"fmt"

	"github.com/LindsayBradford/crm/config"
	crmerrors "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
)

const defaultLoggerIndex = 0

// LogHandlerBuilder is a fluent constructor of a LogHandler, allowing the caller to specify various formatters and
// log-level destinations to best suit their needs.
type LogHandlerBuilder struct {
	logHandler  LogHandler
	buildErrors *crmerrors.CompositeError
}

func (builder *LogHandlerBuilder) FromConfig(loggingConfig []config.LoggerConfig) ([]LogHandler, error) {
	listError := new(crmerrors.CompositeError)

	handlerList := make([]LogHandler, 1)
	handlerList[defaultLoggerIndex] = builder.buildDefaultLogHandler(listError)

	for _, currConfig := range loggingConfig {
		newLogger, newLogError := builder.deriveLogHandler(currConfig, listError)
		if newLogError != nil {
			listError.Add(newLogError)
		}
		handlerList = append(handlerList, newLogger)
	}

	if listError.Size() > 0 {
		return nil, listError
	}
	return handlerList, nil
}

func (builder *LogHandlerBuilder) buildDefaultLogHandler(listError *crmerrors.CompositeError) LogHandler {
	defaultLogger, defaultLogError := builder.ForDefaultLogHandler().Build()
	if defaultLogError != nil {
		// TODO: Prime candidate for error wrapping?
		listError.Add(errors.New("failed creating default log handler"))
	}
	return defaultLogger
}

func (builder *LogHandlerBuilder) deriveLogHandler(currConfig config.LoggerConfig, listError *crmerrors.CompositeError) (LogHandler, error) {
	builder.deriveBaseLogHandler(currConfig)
	builder.deriveConfiguredLogLevelDestinations(currConfig, listError)
	return builder.Build()
}

func (builder *LogHandlerBuilder) deriveBaseLogHandler(currConfig config.LoggerConfig) {
	switch currConfig.Type {
	case config.NativeLibrary, config.UnspecifiedLoggerType:
		builder.
			ForNativeLibraryLogHandler().
			WithFormatter(deriveLogFormatter(currConfig)).
			WithName(currConfig.Name)
	case config.BareBones:
		builder.
			ForBareBonesLogHandler().
			WithFormatter(deriveLogFormatter(currConfig)).
			WithName(currConfig.Name)
	default:
		panic("Should not reach here")
	}
}

func (builder *LogHandlerBuilder) deriveConfiguredLogLevelDestinations(currConfig config.LoggerConfig, listError *crmerrors.CompositeError) {
	for configLogLevel, configDestination := range currConfig.LogLevelDestinations {
		logLevel, destination := deriveLogLevelAndDestination(configLogLevel, configDestination, listError)
		builder.WithLogLevelDestination(logLevel, destination)
	}
}

func deriveLogLevelAndDestination(configLogLevel string, configDestination string, listError *crmerrors.CompositeError) (LogLevel, LogDestination) {
	logLevel := deriveLogLevel(configLogLevel)
	destination := deriveDestination(configDestination, configLogLevel, listError)
	return logLevel, destination
}

func deriveLogLevel(configLogLevel string) LogLevel {
	var derivedLogLevel LogLevel
	switch configLogLevel {
	case "Debugging":
		derivedLogLevel = DEBUG
	case "Information":
		derivedLogLevel = INFO
	case "Warnings":
		derivedLogLevel = WARN
	case "Errors":
		derivedLogLevel = ERROR
	default:
		derivedLogLevel = LogLevel(configLogLevel)
	}
	return derivedLogLevel
}

func deriveDestination(configDestination string, configLogLevel string, listError *crmerrors.CompositeError) LogDestination {
	var derivedDestination LogDestination
	switch configDestination {
	case "StandardOutput":
		derivedDestination = STDOUT
	case "StandardError":
		derivedDestination = STDERR
	case "Discarded":
		derivedDestination = DISCARD
	default:
		listError.Add(
			fmt.Errorf("attempted to map log level [%s] to unrecognised destination [%s]",
				configLogLevel, configDestination))
	}
	return derivedDestination
}

func deriveLogFormatter(loggerConfig config.LoggerConfig) LogFormatter {
	switch loggerConfig.Formatter {
	case config.RawMessage, config.UnspecifiedFormatterType:
		return new(RawMessageFormatter)
	case config.Json:
		return new(JsonFormatter)
	case config.NameValuePair:
		return new(NameValuePairFormatter)
	default:
		panic("Should not reach here")
	}
}

func (builder *LogHandlerBuilder) ForDefaultLogHandler() *LogHandlerBuilder {
	return builder.
		ForNativeLibraryLogHandler().
		WithName("DefaultLogHandler").
		WithFormatter(new(RawMessageFormatter))
}

// ForNativeLibraryLogHandler instructs LogHandlerBuilder to use the native built-in go library wrapper as its
// LogHandler
func (builder *LogHandlerBuilder) ForNativeLibraryLogHandler() *LogHandlerBuilder {
	builder.buildErrors = crmerrors.NewComposite("Failed to build valid LogHandler")

	newHandler := new(NativeLibraryLogHandler)

	defaultDestinations := new(LogLevelDestinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(NullFormatter))
	newHandler.Initialise()

	builder.logHandler = newHandler
	return builder
}

// ForNativeLibraryLogHandler instructs LogHandlerBuilder to use the native built-in go library wrapper as its
// LogHandler
func (builder *LogHandlerBuilder) ForBareBonesLogHandler() *LogHandlerBuilder {
	builder.buildErrors = crmerrors.NewComposite("Failed to build valid LogHandler")

	newHandler := new(BareBonesLogHandler)

	defaultDestinations := new(LogLevelDestinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(NullFormatter))
	newHandler.Initialise()

	builder.logHandler = newHandler
	return builder
}

// WithName instructs LogHandlerBuilder to label the LogHandler being built with the specified human-friendly name.
func (builder *LogHandlerBuilder) WithName(name string) *LogHandlerBuilder {
	handlerBeingBuilt := builder.logHandler
	handlerBeingBuilt.SetName(name)
	return builder
}

// WithFormatter instructs LogHandlerBuilder to ensure that the LogHandler constructed will use formatter for its log
// entry formatting. If not called, the default NullFormatter will be used.
func (builder *LogHandlerBuilder) WithFormatter(formatter LogFormatter) *LogHandlerBuilder {
	formatter.Initialise()

	handlerBeingBuilt := builder.logHandler
	handlerBeingBuilt.SetFormatter(formatter)

	return builder
}

// WithLogLevelDestination instructs LogHandlerBuilder to override the existing LogLevelDestinations with a new
// destination for the given logLevel.
func (builder *LogHandlerBuilder) WithLogLevelDestination(logLevel LogLevel, destination LogDestination) *LogHandlerBuilder {
	handlerBeingBuilt := builder.logHandler

	handlerDestinations := handlerBeingBuilt.Destinations()
	handlerDestinations.Override(logLevel, destination)
	if nativeLibraryHandler, ok := handlerBeingBuilt.(*NativeLibraryLogHandler); ok {
		nativeLibraryHandler.addLogLevel(logLevel)
	}

	return builder
}

// Build instructs LogHandlerBuilder to finalise building its LogHandler, and return it to he caller.
func (builder *LogHandlerBuilder) Build() (LogHandler, error) {
	handlerBeingBuilt := builder.logHandler
	if builder.buildErrors.Size() == 0 {
		return handlerBeingBuilt, nil
	} else {
		return handlerBeingBuilt, builder.buildErrors
	}
}
