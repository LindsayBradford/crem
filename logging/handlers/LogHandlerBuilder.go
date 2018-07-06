// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package handlers

import (
	crmerrors "github.com/LindsayBradford/crm/errors"
	. "github.com/LindsayBradford/crm/logging/formatters"
	. "github.com/LindsayBradford/crm/logging/shared"
)

// LogHandlerBuilder is a fluent constructor of a LogHandler, allowing the caller to specify various formatters and
// log-level destinations to best suit their needs.
type LogHandlerBuilder struct {
	logHandler  LogHandler
	buildErrors *crmerrors.CompositeError
}

// ForNativeLibraryLogHandler instructs LogHandlerBuilder to use the native built-in go library wrapper as its
// LogHandler
func (this *LogHandlerBuilder) ForNativeLibraryLogHandler() *LogHandlerBuilder {
	this.buildErrors = crmerrors.NewComposite("Failed to build valid LogHandler")

	newHandler := new(NativeLibraryLogHandler)

	defaultDestinations := new(LogLevelDestinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(NullFormatter))
	newHandler.Initialise()

	this.logHandler = newHandler
	return this
}

// ForNativeLibraryLogHandler instructs LogHandlerBuilder to use the native built-in go library wrapper as its
// LogHandler
func (this *LogHandlerBuilder) ForBareBonesLogHandler() *LogHandlerBuilder {
	this.buildErrors = crmerrors.NewComposite("Failed to build valid LogHandler")

	newHandler := new(BareBonesLogHandler)

	defaultDestinations := new(LogLevelDestinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(NullFormatter))
	newHandler.Initialise()

	this.logHandler = newHandler
	return this
}

// WithFormatter instructs LogHandlerBuilder to ensure that the LogHandler constructed will use formatter for its log
// entry formatting. If not called, the default NullFormatter will be used.
func (this *LogHandlerBuilder) WithFormatter(formatter LogFormatter) *LogHandlerBuilder {
	formatter.Initialise()
	this.logHandler.SetFormatter(formatter)
	return this
}

// WithLogLevelDestination instructs LogHandlerBuilder to override the existing LogLevelDestinations with a new
// destination for the given logLevel.
func (this *LogHandlerBuilder) WithLogLevelDestination(logLevel LogLevel, destination LogDestination) *LogHandlerBuilder {
	this.logHandler.Destinations().Override(logLevel, destination)
	if nativeLibraryHandler, ok := this.logHandler.(*NativeLibraryLogHandler); ok {
		nativeLibraryHandler.AddLogLevel(logLevel)
	}

	return this
}

// Build instructs LogHandlerBuilder to finalise building its LogHandler, and return it to he caller.
func (this *LogHandlerBuilder) Build() (LogHandler, error) {
	if this.buildErrors.Size() == 0 {
		return this.logHandler, nil
	} else {
		return this.logHandler, this.buildErrors
	}
}