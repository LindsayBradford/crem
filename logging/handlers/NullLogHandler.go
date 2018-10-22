// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	. "github.com/LindsayBradford/crem/logging/formatters"
	. "github.com/LindsayBradford/crem/logging/shared"
)

var DefaultNullLogHandler = new(NullLogHandler)

type NullLogHandler struct{}

func (handler *NullLogHandler) IsDefault() bool                                   { return false }
func (handler *NullLogHandler) SetAsDefault(isDefault bool)                       {}
func (handler *NullLogHandler) Name() string                                      { return "NULL" }
func (handler *NullLogHandler) SetName(name string)                               {}
func (handler *NullLogHandler) Debug(message interface{})                         {}
func (handler *NullLogHandler) Info(message interface{})                          {}
func (handler *NullLogHandler) Warn(message interface{})                          {}
func (handler *NullLogHandler) Error(message interface{})                         {}
func (handler *NullLogHandler) LogAtLevel(logLevel LogLevel, message interface{}) {}

func (handler *NullLogHandler) Initialise()                                            {}
func (handler *NullLogHandler) SetDestinations(*LogLevelDestinations)                  {}
func (handler *NullLogHandler) Destinations() *LogLevelDestinations                    { return nil }
func (handler *NullLogHandler) SetFormatter(formatter LogFormatter)                    {}
func (handler *NullLogHandler) Formatter() LogFormatter                                { return &NullFormatter{} }
func (handler *NullLogHandler) BeingDiscarded(logLevel LogLevel) bool                  { return true }
func (handler *NullLogHandler) SupportsLogLevel(logLevel LogLevel) bool                { return true }
func (handler *NullLogHandler) Override(logLevel LogLevel, destination LogDestination) {}
