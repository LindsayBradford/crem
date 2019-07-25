// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"fmt"

	annealingObserver "github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/formatters"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

type LoggingConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	loggerBuilder *loggers.Builder
	logger        logging.Logger
}

func NewLoggingConfigInterpreter() *LoggingConfigInterpreter {
	interpreter := new(LoggingConfigInterpreter).initialise()
	return interpreter
}

func (i *LoggingConfigInterpreter) initialise() *LoggingConfigInterpreter {
	i.errors = compositeErrors.New("Logging Configuration")

	i.loggerBuilder = new(loggers.Builder)

	var buildErrors error

	i.logger, buildErrors = i.loggerBuilder.ForNativeLibraryLogHandler().
		WithName("DefaultLogHandler").
		WithFormatter(new(formatters.RawMessageFormatter)).
		WithLogLevelDestination(annealingObserver.AnnealerLogLevel, logging.STDOUT).
		WithLogLevelDestination(model.LogLevel, logging.DISCARD).
		Build()

	if buildErrors != nil {
		i.errors.Add(buildErrors)
	}

	return i
}

func (i *LoggingConfigInterpreter) Interpret(config *data.LoggingConfig) *LoggingConfigInterpreter {
	i.deriveLogHandler(config)
	i.deriveLogLevelDestinations(config)
	i.logger, _ = i.loggerBuilder.Build()

	return i
}

func (i *LoggingConfigInterpreter) deriveLogHandler(config *data.LoggingConfig) {
	formatter := deriveLogFormatter(config.LoggingFormatter)
	switch config.LoggingType {
	case data.NativeLibrary, data.UnspecifiedLoggerType:
		i.loggerBuilder.
			ForNativeLibraryLogHandler().
			WithFormatter(formatter).
			WithLogLevelDestination(annealingObserver.AnnealerLogLevel, logging.STDOUT).
			WithLogLevelDestination(model.LogLevel, logging.DISCARD)
	case data.BareBones:
		i.loggerBuilder.
			ForBareBonesLogHandler().
			WithFormatter(formatter).
			WithLogLevelDestination(annealingObserver.AnnealerLogLevel, logging.STDOUT).
			WithLogLevelDestination(model.LogLevel, logging.DISCARD)
	default:
		panic("Should not reach here")
	}
}

func (i *LoggingConfigInterpreter) deriveLogLevelDestinations(config *data.LoggingConfig) {
	for configLogLevel, configDestination := range config.LogLevelDestinations {
		logLevel, destination := i.deriveLogLevelAndDestination(configLogLevel, configDestination)
		i.loggerBuilder.WithLogLevelDestination(logLevel, destination)
	}
}

func deriveLogFormatter(formatterType data.FormatterType) logging.Formatter {
	switch formatterType {
	case data.RawMessage, data.UnspecifiedFormatterType:
		return new(formatters.RawMessageFormatter)
	case data.Json:
		return new(formatters.JsonFormatter)
	case data.NameValuePair:
		return new(formatters.NameValuePairFormatter)
	default:
		panic("Should not reach here")
	}
}

func (i *LoggingConfigInterpreter) deriveLogLevelAndDestination(configLogLevel string, configDestination string) (logging.Level, logging.Destination) {
	logLevel := i.deriveLogLevel(configLogLevel)
	destination := i.deriveDestination(configDestination, configLogLevel)
	return logLevel, destination
}

func (i *LoggingConfigInterpreter) deriveLogLevel(configLogLevel string) logging.Level {
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

func (i *LoggingConfigInterpreter) deriveDestination(configDestination string, configLogLevel string) logging.Destination {
	var derivedDestination logging.Destination
	switch configDestination {
	case "StandardOutput":
		derivedDestination = logging.STDOUT
	case "StandardError":
		derivedDestination = logging.STDERR
	case "Discarded":
		derivedDestination = logging.DISCARD
	default:
		i.errors.Add(
			fmt.Errorf("attempted to map log level [%s] to unrecognised destination [%s]",
				configLogLevel, configDestination))
	}
	return derivedDestination
}

func (i *LoggingConfigInterpreter) LogHandler() logging.Logger {
	return i.logger
}

func (i *LoggingConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
