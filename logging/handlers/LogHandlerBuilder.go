// Copyright (c) 2018 Australian Rivers Institute.

package handlers

import (
	cremerrors "github.com/LindsayBradford/crem/errors"
	. "github.com/LindsayBradford/crem/logging/formatters"
	. "github.com/LindsayBradford/crem/logging/shared"
)

// LogHandlerBuilder is a fluent constructor of a LogHandler, allowing the caller to specify various formatters and
// log-level destinations to best suit their needs.
type LogHandlerBuilder struct {
	logHandler  LogHandler
	buildErrors *cremerrors.CompositeError
}

func (builder *LogHandlerBuilder) ForDefaultLogHandler() *LogHandlerBuilder {
	return builder.
		ForNativeLibraryLogHandler().
		WithName("DefaultLogHandler").
		WithFormatter(new(RawMessageFormatter))
}

// ForNativeLibraryLogHandler instructs LogHandlerBuilder to use the native built-in go library wrapper as its
// LogHandler
func (builder *LogHandlerBuilder) ForNativeLibraryLogHandler() *LogHandlerBuilder {
	builder.buildErrors = cremerrors.NewComposite("Failed to build valid LogHandler")

	newHandler := new(NativeLibraryLogHandler)

	defaultDestinations := new(LogLevelDestinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(NullFormatter))
	newHandler.Initialise()

	builder.logHandler = newHandler
	return builder
}

// ForNativeLibraryLogHandler instructs LogHandlerBuilder to use the native built-in go library wrapper as its
// LogHandler
func (builder *LogHandlerBuilder) ForBareBonesLogHandler() *LogHandlerBuilder {
	builder.buildErrors = cremerrors.NewComposite("Failed to build valid LogHandler")

	newHandler := new(BareBonesLogHandler)

	defaultDestinations := new(LogLevelDestinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(NullFormatter))
	newHandler.Initialise()

	builder.logHandler = newHandler
	return builder
}

// WithName instructs LogHandlerBuilder to label the LogHandler being built with the specified human-friendly name.
func (builder *LogHandlerBuilder) WithName(name string) *LogHandlerBuilder {
	handlerBeingBuilt := builder.logHandler
	handlerBeingBuilt.SetName(name)
	return builder
}

// WithFormatter instructs LogHandlerBuilder to ensure that the LogHandler constructed will use formatter for its log
// entry formatting. If not called, the default NullFormatter will be used.
func (builder *LogHandlerBuilder) WithFormatter(formatter LogFormatter) *LogHandlerBuilder {
	formatter.Initialise()

	handlerBeingBuilt := builder.logHandler
	handlerBeingBuilt.SetFormatter(formatter)

	return builder
}

// WithLogLevelDestination instructs LogHandlerBuilder to override the existing LogLevelDestinations with a new
// destination for the given logLevel.
func (builder *LogHandlerBuilder) WithLogLevelDestination(logLevel LogLevel, destination LogDestination) *LogHandlerBuilder {
	handlerBeingBuilt := builder.logHandler

	handlerDestinations := handlerBeingBuilt.Destinations()
	handlerDestinations.Override(logLevel, destination)
	if nativeLibraryHandler, ok := handlerBeingBuilt.(*NativeLibraryLogHandler); ok {
		nativeLibraryHandler.addLogLevel(logLevel)
	}

	return builder
}

// Build instructs LogHandlerBuilder to finalise building its LogHandler, and return it to he caller.
func (builder *LogHandlerBuilder) Build() (LogHandler, error) {
	handlerBeingBuilt := builder.logHandler
	if builder.buildErrors.Size() == 0 {
		return handlerBeingBuilt, nil
	} else {
		return handlerBeingBuilt, builder.buildErrors
	}
}
