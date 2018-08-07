// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	"fmt"
	"log"

	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
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

func (handler *NativeLibraryLogHandler) Debug(message string) {
	handler.LogAtLevel(DEBUG, message)
}

func (handler *NativeLibraryLogHandler) DebugWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(DEBUG, logAttributes)
}

func (handler *NativeLibraryLogHandler) Info(message string) {
	handler.LogAtLevel(INFO, message)
}

func (handler *NativeLibraryLogHandler) InfoWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(INFO, logAttributes)
}

func (handler *NativeLibraryLogHandler) Warn(message string) {
	handler.LogAtLevel(WARN, message)
}

func (handler *NativeLibraryLogHandler) WarnWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(WARN, logAttributes)
}

func (handler *NativeLibraryLogHandler) Error(message string) {
	handler.LogAtLevel(ERROR, message)
}

func (handler *NativeLibraryLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {
	handler.LogAtLevelWithAttributes(ERROR, logAttributes)
}

func (handler *NativeLibraryLogHandler) ErrorWithError(err error) {
	handler.LogAtLevel(ERROR, fmt.Sprintf(err.Error()))
}

func (handler *NativeLibraryLogHandler) LogAtLevel(logLevel LogLevel, message string) {
	logAttributes := LogAttributes{NameValuePair{Name: MessageNameLabel, Value: message}}
	handler.deriveDestination(logLevel).Println("[" + string(logLevel) + "] " + handler.formatter.Format(logAttributes))
}

func (handler *NativeLibraryLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {
	handler.deriveDestination(logLevel).Println("[" + string(logLevel) + "] " + handler.formatter.Format(logAttributes))
}

func (handler *NativeLibraryLogHandler) deriveDestination(logLevel LogLevel) *log.Logger {
	return handler.loggerMap[logLevel]
}
