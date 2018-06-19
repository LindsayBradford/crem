// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	"fmt"
	"log"
)

type NativeLibraryLogHandler struct {
	LogHandlerBase
	debug    *log.Logger
	info    *log.Logger
	warn    *log.Logger
	error    *log.Logger
}

func (this *NativeLibraryLogHandler) Initialise() {
	this.debug = log.New(this.destinations.destinations[DEBUG], "", log.Lshortfile|log.Ldate|log.Ltime|log.Lmicroseconds)
	this.info = log.New(this.destinations.destinations[INFO], "", log.Ldate|log.Ltime|log.Lmicroseconds)
	this.warn = log.New(this.destinations.destinations[WARN], "", log.Ldate|log.Ltime|log.Lmicroseconds)
	this.error = log.New(this.destinations.destinations[ERROR], "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func (this *NativeLibraryLogHandler) WithFormatter(formatter LogAttributeFormatter) *NativeLibraryLogHandler {
	formatter.Initialise()
	this.formatter = formatter
	return this
}

func (this *NativeLibraryLogHandler) Debug(message string) {
	logAttributes := LogAttributes{ NameValuePair{ "message", message }}
	this.debug.Println("DEBUG " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) DebugWithAttributes(logAttributes LogAttributes) {
	this.debug.Println("DEBUG " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) Info(message string) {
	logAttributes := LogAttributes{ NameValuePair{ "message", message }}
	this.info.Println("INFO " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) InfoWithAttributes(logAttributes LogAttributes) {
	this.info.Println("INFO " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) Warn(message string) {
	logAttributes := LogAttributes{ NameValuePair{ "message", message }}
	this.warn.Println("WARN " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) WarnWithAttributes(logAttributes LogAttributes) {
	this.warn.Println("WARN " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) Error(message string) {
	logAttributes := LogAttributes{ NameValuePair{ "message", message }}
	this.error.Println("ERROR " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {
	this.error.Println("ERROR " + this.formatter.Format(logAttributes))
}

func (this *NativeLibraryLogHandler) ErrorWithError(err error) {
	this.error.Println("ERROR " + fmt.Sprintf(err.Error()))
}