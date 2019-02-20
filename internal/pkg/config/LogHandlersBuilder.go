// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	. "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	. "github.com/LindsayBradford/crem/pkg/logging/formatters"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/pkg/errors"
)

const defaultLoggerIndex = 0

var (
	baseBuilder = new(loggers.Builder)
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

func (builder *LogHandlersBuilder) Build() ([]logging.Logger, error) {

	handlerList := make([]logging.Logger, 0)
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

func (builder *LogHandlersBuilder) makeFirstConfigSuppliedLoggerTheDefault(handlerList []logging.Logger) []logging.Logger {
	if len(builder.config) > 0 && builder.errors.Size() == 0 {
		return handlerList[1:]
	}
	return handlerList
}

func (builder *LogHandlersBuilder) newHandlerFor(currConfig LoggerConfig) logging.Logger {
	newLogger, newLogError := builder.deriveLogHandler(currConfig)
	if newLogError != nil {
		builder.errors.Add(newLogError)
	} else {
		ensureSupportForAnnealerLogLevel(newLogger)
		ensureSupportForModelLogLevel(newLogger)
	}
	return newLogger
}

func (builder *LogHandlersBuilder) buildDefaultLogHandler() logging.Logger {
	defaultLogger, defaultLogError := baseBuilder.ForDefaultLogHandler().Build()
	if defaultLogError != nil {
		builder.errors.Add(
			errors.Wrap(defaultLogError, "failed creating default log handler"),
		)
	}
	ensureSupportForAnnealerLogLevel(defaultLogger)
	ensureSupportForModelLogLevel(defaultLogger)
	return defaultLogger
}

func ensureSupportForAnnealerLogLevel(handler logging.Logger) {
	if !handler.SupportsLogLevel(observer.AnnealerLogLevel) {
		handler.Override(observer.AnnealerLogLevel, logging.STDOUT)
	}
}

func ensureSupportForModelLogLevel(handler logging.Logger) {
	if !handler.SupportsLogLevel(model.LogLevel) {
		handler.Override(model.LogLevel, logging.DISCARD)
	}
}

func (builder *LogHandlersBuilder) deriveLogHandler(currConfig LoggerConfig) (logging.Logger, error) {
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

func (builder *LogHandlersBuilder) deriveLogLevelAndDestination(configLogLevel string, configDestination string) (logging.Level, logging.Destination) {
	logLevel := builder.deriveLogLevel(configLogLevel)
	destination := builder.deriveDestination(configDestination, configLogLevel)
	return logLevel, destination
}

func (builder *LogHandlersBuilder) deriveLogLevel(configLogLevel string) logging.Level {
	var derivedLogLevel logging.Level
	switch configLogLevel {
	case "Debugging":
		derivedLogLevel = logging.DEBUG
	case "Information":
		derivedLogLevel = logging.INFO
	case "Warnings":
		derivedLogLevel = logging.WARN
	case "Errors":
		derivedLogLevel = logging.ERROR
	default:
		derivedLogLevel = logging.Level(configLogLevel)
	}
	return derivedLogLevel
}

func (builder *LogHandlersBuilder) deriveDestination(configDestination string, configLogLevel string) logging.Destination {
	var derivedDestination logging.Destination
	switch configDestination {
	case "StandardOutput":
		derivedDestination = logging.STDOUT
	case "StandardError":
		derivedDestination = logging.STDERR
	case "Discarded":
		derivedDestination = logging.DISCARD
	default:
		builder.errors.Add(
			fmt.Errorf("attempted to map log level [%s] to unrecognised destination [%s]",
				configLogLevel, configDestination))
	}
	return derivedDestination
}

func deriveLogFormatter(loggerConfig LoggerConfig) logging.Formatter {
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
