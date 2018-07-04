// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
)

type NullLogHandler struct {}

func (this *NullLogHandler) Debug(message string) {}
func (this *NullLogHandler) DebugWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) Info(message string) {}
func (this *NullLogHandler) InfoWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) Warn(message string) {}
func (this *NullLogHandler) WarnWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) Error(message string) {}
func (this *NullLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {}
func (this *NullLogHandler) ErrorWithError(err error) {}
func (this *NullLogHandler) LogAtLevel(logLevel LogLevel, message string) {}
func (this *NullLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {}
func (this *NullLogHandler) Initialise() {}
func (this *NullLogHandler) SetDestinations(*LogLevelDestinations) {}
func (this *NullLogHandler) Destinations() *LogLevelDestinations {return nil}
func (this *NullLogHandler) SetFormatter(formatter LogFormatter) {}
func (this *NullLogHandler) Formatter() LogFormatter{return &NullFormatter{}}
func (this *NullLogHandler) BeingDiscarded(logLevel LogLevel) bool { return true }
