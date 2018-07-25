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

func (this *NativeLibraryLogHandler) Initialise() {
	this.loggerMap = make(map[LogLevel]*log.Logger)
	this.AddLogLevel(DEBUG).AddLogLevel(INFO).AddLogLevel(WARN).AddLogLevel(ERROR)
}

func (this *NativeLibraryLogHandler) AddLogLevel(logLevel LogLevel) *NativeLibraryLogHandler {
	this.loggerMap[logLevel] = log.New(this.destinations.Destinations[logLevel], "", metadataMask)
	return this
}

func (this *NativeLibraryLogHandler) WithFormatter(formatter LogFormatter) *NativeLibraryLogHandler {
	formatter.Initialise()
	this.formatter = formatter
	return this
}

func (this *NativeLibraryLogHandler) Debug(message string) {
	this.LogAtLevel(DEBUG, message)
}

func (this *NativeLibraryLogHandler) DebugWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(DEBUG, logAttributes)
}

func (this *NativeLibraryLogHandler) Info(message string) {
	this.LogAtLevel(INFO, message)
}

func (this *NativeLibraryLogHandler) InfoWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(INFO, logAttributes)
}

func (this *NativeLibraryLogHandler) Warn(message string) {
	this.LogAtLevel(WARN, message)
}

func (this *NativeLibraryLogHandler) WarnWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(WARN, logAttributes)
}

func (this *NativeLibraryLogHandler) Error(message string) {
	this.LogAtLevel(ERROR, message)
}

func (this *NativeLibraryLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {
	this.LogAtLevelWithAttributes(ERROR, logAttributes)
}

func (this *NativeLibraryLogHandler) ErrorWithError(err error) {
	this.LogAtLevel(ERROR, fmt.Sprintf(err.Error()))
}

func (this *NativeLibraryLogHandler) LogAtLevel(logLevel LogLevel, message string) {
	logAttributes := LogAttributes{NameValuePair{MESSAGE_LABEL, message}}
	this.loggerMap[logLevel].Println("[" + string(logLevel) + "] " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {
	this.loggerMap[logLevel].Println("[" + string(logLevel) + "] " + this.formatter.Format(logAttributes))
}
