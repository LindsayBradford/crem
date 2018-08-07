// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
)

var DefaultNullLogHandler = new(NullLogHandler)

type NullLogHandler struct{}

func (handler *NullLogHandler) IsDefault() bool                                 { return false }
func (handler *NullLogHandler) SetAsDefault(isDefault bool)                     {}
func (handler *NullLogHandler) Name() string                                    { return "NULL" }
func (handler *NullLogHandler) SetName(name string)                             {}
func (handler *NullLogHandler) Debug(message string)                            {}
func (handler *NullLogHandler) DebugWithAttributes(logAttributes LogAttributes) {}
func (handler *NullLogHandler) Info(message string)                             {}
func (handler *NullLogHandler) InfoWithAttributes(logAttributes LogAttributes)  {}
func (handler *NullLogHandler) Warn(message string)                             {}
func (handler *NullLogHandler) WarnWithAttributes(logAttributes LogAttributes)  {}
func (handler *NullLogHandler) Error(message string)                            {}
func (handler *NullLogHandler) ErrorWithAttributes(logAttributes LogAttributes) {}
func (handler *NullLogHandler) ErrorWithError(err error)                        {}
func (handler *NullLogHandler) LogAtLevel(logLevel LogLevel, message string)    {}
func (handler *NullLogHandler) LogAtLevelWithAttributes(logLevel LogLevel, logAttributes LogAttributes) {
}
func (handler *NullLogHandler) Initialise()                                            {}
func (handler *NullLogHandler) SetDestinations(*LogLevelDestinations)                  {}
func (handler *NullLogHandler) Destinations() *LogLevelDestinations                    { return nil }
func (handler *NullLogHandler) SetFormatter(formatter LogFormatter)                    {}
func (handler *NullLogHandler) Formatter() LogFormatter                                { return &NullFormatter{} }
func (handler *NullLogHandler) BeingDiscarded(logLevel LogLevel) bool                  { return true }
func (handler *NullLogHandler) SupportsLogLevel(logLevel LogLevel) bool                { return true }
func (handler *NullLogHandler) Override(logLevel LogLevel, destination LogDestination) {}
