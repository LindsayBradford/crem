// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"errors"
	"fmt"

	. "github.com/LindsayBradford/crm/annealing/logging"
	"github.com/LindsayBradford/crm/config"
	. "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
)

var logBuilder = new(LogHandlerBuilder)

func BuildLogHandlers(loggingConfig []config.LoggerConfig) ([]LogHandler, error) {
	listError := new(CompositeError)

	handlerList := make([]LogHandler, 1)
	handlerList[defaultLoggerIndex] = buildDefaultLogHandler(listError)

	for _, currConfig := range loggingConfig {
		newLogger, newLogError := deriveLogHandler(currConfig, listError)
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

func buildDefaultLogHandler(listError *CompositeError) LogHandler {
	defaultLogger, defaultLogError :=
		logBuilder.ForDefaultLogHandler().
			WithLogLevelDestination(AnnealerLogLevel, STDOUT).
			Build()
	if defaultLogError != nil {
		// TODO: Prime candidate for error wrapping?
		listError.Add(errors.New("failed creating default log handler"))
	}
	return defaultLogger
}

func deriveLogHandler(currConfig config.LoggerConfig, listError *CompositeError) (LogHandler, error) {
	deriveBaseLogHandler(currConfig)
	logBuilder.WithLogLevelDestination(AnnealerLogLevel, STDOUT)
	deriveConfiguredLogLevelDestinations(currConfig, listError)
	return logBuilder.Build()
}

func deriveBaseLogHandler(currConfig config.LoggerConfig) {
	switch currConfig.Type {
	case config.NativeLibrary, config.UnspecifiedLoggerType:
		logBuilder.
			ForNativeLibraryLogHandler().
			WithFormatter(deriveLogFormatter(currConfig)).
			WithName(currConfig.Name)
	case config.BareBones:
		logBuilder.
			ForBareBonesLogHandler().
			WithFormatter(deriveLogFormatter(currConfig)).
			WithName(currConfig.Name)
	default:
		panic("Should not reach here")
	}
}

func deriveConfiguredLogLevelDestinations(currConfig config.LoggerConfig, listError *CompositeError) {
	for configLogLevel, configDestination := range currConfig.LogLevelDestinations {
		logLevel, destination := deriveLogLevelAndDestination(configLogLevel, configDestination, listError)
		logBuilder.WithLogLevelDestination(logLevel, destination)
	}
}

func deriveLogLevelAndDestination(configLogLevel string, configDestination string, listError *CompositeError) (LogLevel, LogDestination) {
	logLevel := deriveLogLevel(configLogLevel, listError)
	destination := deriveDestination(configDestination, listError, configLogLevel)
	return logLevel, destination
}

func deriveLogLevel(configLogLevel string, listError *CompositeError) LogLevel {
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
	case "Annealing":
		derivedLogLevel = AnnealerLogLevel
	default:
		listError.Add(fmt.Errorf("attempted to map to unrecognised log level [%s]", configLogLevel))
	}
	return derivedLogLevel
}

func deriveDestination(configDestination string, listError *CompositeError, key string) LogDestination {
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
			fmt.Errorf("attempted to map log level [%s] to unrecognised destination [%s]", key, configDestination),
		)
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
