// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package loggers

import (
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/formatters"
)

var DefaultNullLogLogger = new(NullLogger)

type NullLogger struct{}

func (handler *NullLogger) IsDefault() bool                                        { return false }
func (handler *NullLogger) SetAsDefault(isDefault bool)                            {}
func (handler *NullLogger) Name() string                                           { return "NULL" }
func (handler *NullLogger) SetName(name string)                                    {}
func (handler *NullLogger) Debug(message interface{})                              {}
func (handler *NullLogger) Info(message interface{})                               {}
func (handler *NullLogger) Warn(message interface{})                               {}
func (handler *NullLogger) Error(message interface{})                              {}
func (handler *NullLogger) LogAtLevel(logLevel logging.Level, message interface{}) {}

func (handler *NullLogger) Initialise()                                                      {}
func (handler *NullLogger) SetDestinations(*logging.Destinations)                            {}
func (handler *NullLogger) Destinations() *logging.Destinations                              { return nil }
func (handler *NullLogger) SetFormatter(formatter logging.Formatter)                         {}
func (handler *NullLogger) Formatter() logging.Formatter                                     { return &formatters.NullFormatter{} }
func (handler *NullLogger) BeingDiscarded(logLevel logging.Level) bool                       { return true }
func (handler *NullLogger) SupportsLogLevel(logLevel logging.Level) bool                     { return true }
func (handler *NullLogger) Override(logLevel logging.Level, destination logging.Destination) {}
