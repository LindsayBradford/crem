// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	"log"

	. "github.com/LindsayBradford/crem/logging/formatters"
	. "github.com/LindsayBradford/crem/logging/shared"
)

type NativeLibraryLogHandler struct {
	LogHandlerBase

	loggerMap map[LogLevel]*log.Logger
}

const metadataMask = log.Ldate | log.Ltime | log.Lmicroseconds

func (handler *NativeLibraryLogHandler) Initialise() {
	handler.loggerMap = make(map[LogLevel]*log.Logger)
	handler.addLogLevel(DEBUG).addLogLevel(INFO).addLogLevel(WARN).addLogLevel(ERROR)
}

func (handler *NativeLibraryLogHandler) Override(logLevel LogLevel, destination LogDestination) {
	if !handler.SupportsLogLevel(logLevel) {
		handler.LogHandlerBase.Override(logLevel, destination)
		handler.addLogLevel(logLevel)
	}
}

func (handler *NativeLibraryLogHandler) addLogLevel(logLevel LogLevel) *NativeLibraryLogHandler {
	handler.loggerMap[logLevel] = log.New(handler.destinations.Destinations[logLevel], "", metadataMask)
	return handler
}

func (handler *NativeLibraryLogHandler) WithFormatter(formatter LogFormatter) *NativeLibraryLogHandler {
	formatter.Initialise()
	handler.formatter = formatter
	return handler
}

func (handler *NativeLibraryLogHandler) Debug(message interface{}) {
	handler.LogAtLevel(DEBUG, message)
}

func (handler *NativeLibraryLogHandler) Info(message interface{}) {
	handler.LogAtLevel(INFO, message)
}

func (handler *NativeLibraryLogHandler) Warn(message interface{}) {
	handler.LogAtLevel(WARN, message)
}

func (handler *NativeLibraryLogHandler) Error(message interface{}) {
	handler.LogAtLevel(ERROR, message)
}

func (handler *NativeLibraryLogHandler) LogAtLevel(logLevel LogLevel, message interface{}) {
	messageAttributes := toLogAttributes(message)
	handler.deriveDestination(logLevel).Println("[" + string(logLevel) + "] " + handler.formatter.Format(messageAttributes))
}

func (handler *NativeLibraryLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {
	handler.deriveDestination(logLevel).Println("[" + string(logLevel) + "] " + handler.formatter.Format(logAttributes))
}

func (handler *NativeLibraryLogHandler) deriveDestination(logLevel LogLevel) *log.Logger {
	return handler.loggerMap[logLevel]
}
