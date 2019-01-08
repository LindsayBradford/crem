// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package loggers

import (
	"log"

	"github.com/LindsayBradford/crem/pkg/logging"
)

type NativeLibraryLogger struct {
	LoggerBase

	loggerMap map[logging.Level]*log.Logger
}

const metadataMask = log.Ldate | log.Ltime | log.Lmicroseconds

func (nll *NativeLibraryLogger) Initialise() {
	nll.loggerMap = make(map[logging.Level]*log.Logger)
	nll.addLogLevel(logging.DEBUG).addLogLevel(logging.INFO).addLogLevel(logging.WARN).addLogLevel(logging.ERROR)
}

func (nll *NativeLibraryLogger) Override(logLevel logging.Level, destination logging.Destination) {
	if !nll.SupportsLogLevel(logLevel) {
		nll.LoggerBase.Override(logLevel, destination)
		nll.addLogLevel(logLevel)
	}
}

func (nll *NativeLibraryLogger) addLogLevel(logLevel logging.Level) *NativeLibraryLogger {
	nll.loggerMap[logLevel] = log.New(nll.destinations.Destinations[logLevel], "", metadataMask)
	return nll
}

func (nll *NativeLibraryLogger) WithFormatter(formatter logging.Formatter) *NativeLibraryLogger {
	nll.formatter = formatter
	return nll
}

func (nll *NativeLibraryLogger) Debug(message interface{}) {
	nll.LogAtLevel(logging.DEBUG, message)
}

func (nll *NativeLibraryLogger) Info(message interface{}) {
	nll.LogAtLevel(logging.INFO, message)
}

func (nll *NativeLibraryLogger) Warn(message interface{}) {
	nll.LogAtLevel(logging.WARN, message)
}

func (nll *NativeLibraryLogger) Error(message interface{}) {
	nll.LogAtLevel(logging.ERROR, message)
}

func (nll *NativeLibraryLogger) LogAtLevel(logLevel logging.Level, message interface{}) {
	messageAttributes := toLogAttributes(message)
	nll.deriveDestination(logLevel).Println("[" + string(logLevel) + "] " + nll.formatter.Format(messageAttributes))
}

func (nll *NativeLibraryLogger) LogAtLevelWithAttributes(logLevel logging.Level, logAttributes logging.Attributes) {
	nll.deriveDestination(logLevel).Println("[" + string(logLevel) + "] " + nll.formatter.Format(logAttributes))
}

func (nll *NativeLibraryLogger) deriveDestination(logLevel logging.Level) *log.Logger {
	return nll.loggerMap[logLevel]
}
