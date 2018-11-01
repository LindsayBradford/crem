// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"fmt"
	"github.com/LindsayBradford/crem/annealing/observer"
	. "github.com/LindsayBradford/crem/errors"
	. "github.com/LindsayBradford/crem/logging/formatters"
	"github.com/LindsayBradford/crem/logging/handlers"
	"github.com/LindsayBradford/crem/logging/shared"
	"github.com/pkg/errors"
)

const defaultLoggerIndex = 0

var (
	baseBuilder = new(handlers.LogHandlerBuilder)
)

type LogHandlersBuilder struct {
	errors *CompositeError
	config []LoggerConfig
}

func (builder *LogHandlersBuilder) initialise() *LogHandlersBuilder {
	builder.errors = new(CompositeError)
	return builder
}

func (builder *LogHandlersBuilder) WithConfig(loggingConfig []LoggerConfig) *LogHandlersBuilder {
	builder.initialise()
	builder.config = loggingConfig
	return builder
}

func (builder *LogHandlersBuilder) Build() ([]handlers.LogHandler, error) {

	handlerList := make([]handlers.LogHandler, 0)
	handlerList = append(handlerList, builder.buildDefaultLogHandler()) // system-supplied default logger JIC config sux

	for _, currConfig := range builder.config {
		newLogHandler := builder.newHandlerFor(currConfig)
		handlerList = append(handlerList, newLogHandler)
	}

	handlerList = builder.makeFirstConfigSuppliedLoggerTheDefault(handlerList)

	if builder.errors.Size() > 0 {
		return handlerList, builder.errors
	} else {
		return handlerList, nil
	}
}

func (builder *LogHandlersBuilder) makeFirstConfigSuppliedLoggerTheDefault(handlerList []handlers.LogHandler) []handlers.LogHandler {
	if len(builder.config) > 0 && builder.errors.Size() == 0 {
		return handlerList[1:]
	}
	return handlerList
}

func (builder *LogHandlersBuilder) newHandlerFor(currConfig LoggerConfig) handlers.LogHandler {
	newLogger, newLogError := builder.deriveLogHandler(currConfig)
	if newLogError != nil {
		builder.errors.Add(newLogError)
	} else {
		ensureSupportForAnnealerLogLevel(newLogger)
	}
	return newLogger
}

func (builder *LogHandlersBuilder) buildDefaultLogHandler() handlers.LogHandler {
	defaultLogger, defaultLogError := baseBuilder.ForDefaultLogHandler().Build()
	if defaultLogError != nil {
		builder.errors.Add(
			errors.Wrap(defaultLogError, "failed creating default log handler"),
		)
	}
	ensureSupportForAnnealerLogLevel(defaultLogger)
	return defaultLogger
}

func ensureSupportForAnnealerLogLevel(handler handlers.LogHandler) {
	if !handler.SupportsLogLevel(observer.AnnealerLogLevel) {
		handler.Override(observer.AnnealerLogLevel, shared.STDOUT)
	}
}

func (builder *LogHandlersBuilder) deriveLogHandler(currConfig LoggerConfig) (handlers.LogHandler, error) {
	builder.deriveBaseLogHandler(currConfig)
	builder.deriveConfiguredLogLevelDestinations(currConfig)
	return baseBuilder.Build()
}

func (builder *LogHandlersBuilder) deriveBaseLogHandler(currConfig LoggerConfig) {
	switch currConfig.Type {
	case NativeLibrary, UnspecifiedLoggerType:
		baseBuilder.
			ForNativeLibraryLogHandler().
			WithFormatter(deriveLogFormatter(currConfig)).
			WithName(currConfig.Name)
	case BareBones:
		baseBuilder.
			ForBareBonesLogHandler().
			WithFormatter(deriveLogFormatter(currConfig)).
			WithName(currConfig.Name)
	default:
		panic("Should not reach here")
	}
}

func (builder *LogHandlersBuilder) deriveConfiguredLogLevelDestinations(currConfig LoggerConfig) {
	for configLogLevel, configDestination := range currConfig.LogLevelDestinations {
		logLevel, destination := builder.deriveLogLevelAndDestination(configLogLevel, configDestination)
		baseBuilder.WithLogLevelDestination(logLevel, destination)
	}
}

func (builder *LogHandlersBuilder) deriveLogLevelAndDestination(configLogLevel string, configDestination string) (shared.LogLevel, shared.LogDestination) {
	logLevel := builder.deriveLogLevel(configLogLevel)
	destination := builder.deriveDestination(configDestination, configLogLevel)
	return logLevel, destination
}

func (builder *LogHandlersBuilder) deriveLogLevel(configLogLevel string) shared.LogLevel {
	var derivedLogLevel shared.LogLevel
	switch configLogLevel {
	case "Debugging":
		derivedLogLevel = shared.DEBUG
	case "Information":
		derivedLogLevel = shared.INFO
	case "Warnings":
		derivedLogLevel = shared.WARN
	case "Errors":
		derivedLogLevel = shared.ERROR
	default:
		derivedLogLevel = shared.LogLevel(configLogLevel)
	}
	return derivedLogLevel
}

func (builder *LogHandlersBuilder) deriveDestination(configDestination string, configLogLevel string) shared.LogDestination {
	var derivedDestination shared.LogDestination
	switch configDestination {
	case "StandardOutput":
		derivedDestination = shared.STDOUT
	case "StandardError":
		derivedDestination = shared.STDERR
	case "Discarded":
		derivedDestination = shared.DISCARD
	default:
		builder.errors.Add(
			fmt.Errorf("attempted to map log level [%s] to unrecognised destination [%s]",
				configLogLevel, configDestination))
	}
	return derivedDestination
}

func deriveLogFormatter(loggerConfig LoggerConfig) LogFormatter {
	switch loggerConfig.Formatter {
	case RawMessage, UnspecifiedFormatterType:
		return new(RawMessageFormatter)
	case Json:
		return new(JsonFormatter)
	case NameValuePair:
		return new(NameValuePairFormatter)
	default:
		panic("Should not reach here")
	}
}
