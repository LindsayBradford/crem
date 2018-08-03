// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	. "github.com/LindsayBradford/crm/annealing/logging"
	"github.com/LindsayBradford/crm/config"
	. "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
)

func BuildLogHandlers(crmConfig *config.CRMConfig) ([]LogHandler, error) {
	loggingConfig := crmConfig.Loggers

	if len(loggingConfig) == 0 {
		return buildDefaultLogHandlers()
	} else {
		return buildLogHandlers(loggingConfig)
	}
}

func buildDefaultLogHandlers() ([]LogHandler, error) {
	defaultLogHandler := buildDefaultLogger()
	return []LogHandler{defaultLogHandler}, nil
}

func buildLogHandlers(loggingConfig []config.LoggerConfig) ([]LogHandler, error) {
	listError := new(CompositeError)
	logBuilder := new(LogHandlerBuilder)
	handlerList := make([]LogHandler, len(loggingConfig))
	for index, currConfig := range loggingConfig {

		switch currConfig.Type {
		case config.NativeLibrary, config.UnspecifiedLoggerType:
			logBuilder.
				ForNativeLibraryLogHandler().
				WithFormatter(newFormatterFromConfig(currConfig)).
				WithName(currConfig.Name)
		case config.BareBones:
			logBuilder.
				ForBareBonesLogHandler().
				WithFormatter(newFormatterFromConfig(currConfig)).
				WithName(currConfig.Name)
		default:
			panic("Should not reach here")
		}

		logBuilder.WithLogLevelDestination(AnnealerLogLevel, STDOUT) // default, may be overridden with currConfig below.

		for key, value := range currConfig.LogLevelDestinations {
			var mappedKey LogLevel
			switch key {
			case "Debugging":
				mappedKey = DEBUG
			case "Information":
				mappedKey = INFO
			case "Warnings":
				mappedKey = WARN
			case "Errors":
				mappedKey = ERROR
			case "Annealing":
				mappedKey = AnnealerLogLevel
			}

			var mappedValue LogDestination
			switch value {
			case "StandardOutput":
				mappedValue = STDOUT
			case "StandardError":
				mappedValue = STDERR
			case "Discarded":
				mappedValue = DISCARD
			}
			logBuilder.WithLogLevelDestination(mappedKey, mappedValue)
		}

		var newLogError error
		if handlerList[index], newLogError = logBuilder.Build(); newLogError != nil {
			listError.Add(newLogError)
		}
	}
	if listError.Size() > 0 {
		return nil, listError
	}
	return handlerList, nil
}

func newFormatterFromConfig(loggerConfig config.LoggerConfig) LogFormatter {
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
