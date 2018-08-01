// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	. "github.com/LindsayBradford/crm/annealing/logging"
	"github.com/LindsayBradford/crm/config"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
)

func BuildLogHandlers(crmConfig *config.CRMConfig) []LogHandler {
	loggingConfig := crmConfig.Loggers

	handlerList := make([]LogHandler, len(loggingConfig))
	for index, currConfig := range loggingConfig {
		logBuilder := new(LogHandlerBuilder)

		switch currConfig.Type {
		case "NativeLibrary", "":
			logBuilder.ForNativeLibraryLogHandler().WithName(currConfig.Name).AsDefault(currConfig.Default)
		case "BareBones":
			logBuilder.ForBareBonesLogHandler().WithName(currConfig.Name).AsDefault(currConfig.Default)
		}

		switch currConfig.Formatter {
		case "JSON":
			logBuilder.WithFormatter(new(JsonFormatter))
		case "NameValuePair":
			logBuilder.WithFormatter(new(NameValuePairFormatter))
		case "RawMessage", "":
			logBuilder.WithFormatter(new(RawMessageFormatter))
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

		// TODO: I'm throwing away errors... bad.. fix it.
		handlerList[index], _ = logBuilder.Build()
	}
	return handlerList
}
