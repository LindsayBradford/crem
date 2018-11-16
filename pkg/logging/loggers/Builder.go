// Copyright (c) 2018 Australian Rivers Institute.

package loggers

import (
	cremerrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/formatters"
)

// Builder is a fluent constructor of a Logger, allowing the caller to specify various formatters and
// log-level destinations to best suit their needs.
type Builder struct {
	logHandler  logging.Logger
	buildErrors *cremerrors.CompositeError
}

func (builder *Builder) ForDefaultLogHandler() *Builder {
	return builder.
		ForNativeLibraryLogHandler().
		WithName("DefaultLogHandler").
		WithFormatter(new(formatters.RawMessageFormatter))
}

// ForNativeLibraryLogHandler instructs Builder to use the native built-in go library wrapper as its
// Logger
func (builder *Builder) ForNativeLibraryLogHandler() *Builder {
	builder.buildErrors = cremerrors.New("Failed to build valid Logger")

	newHandler := new(NativeLibraryLogger)

	defaultDestinations := new(logging.Destinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(formatters.NullFormatter))
	newHandler.Initialise()

	builder.logHandler = newHandler
	return builder
}

// ForNativeLibraryLogHandler instructs Builder to use the native built-in go library wrapper as its
// Logger
func (builder *Builder) ForBareBonesLogHandler() *Builder {
	builder.buildErrors = cremerrors.New("Failed to build valid Logger")

	newHandler := new(BareBonesLogger)

	defaultDestinations := new(logging.Destinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)
	newHandler.SetFormatter(new(formatters.NullFormatter))
	newHandler.Initialise()

	builder.logHandler = newHandler
	return builder
}

// WithName instructs Builder to label the Logger being built with the specified human-friendly name.
func (builder *Builder) WithName(name string) *Builder {
	handlerBeingBuilt := builder.logHandler
	handlerBeingBuilt.SetName(name)
	return builder
}

// WithFormatter instructs Builder to ensure that the Logger constructed will use formatter for its log
// entry formatters. If not called, the default NullFormatter will be used.
func (builder *Builder) WithFormatter(formatter logging.Formatter) *Builder {
	formatter.Initialise()

	handlerBeingBuilt := builder.logHandler
	handlerBeingBuilt.SetFormatter(formatter)

	return builder
}

// WithLogLevelDestination instructs Builder to override the existing Destinations with a new
// destination for the given logLevel.
func (builder *Builder) WithLogLevelDestination(logLevel logging.Level, destination logging.Destination) *Builder {
	handlerBeingBuilt := builder.logHandler

	handlerDestinations := handlerBeingBuilt.Destinations()
	handlerDestinations.Override(logLevel, destination)
	if nativeLibraryHandler, ok := handlerBeingBuilt.(*NativeLibraryLogger); ok {
		nativeLibraryHandler.addLogLevel(logLevel)
	}

	return builder
}

// Build instructs Builder to finalise building its Logger, and return it to he caller.
func (builder *Builder) Build() (logging.Logger, error) {
	handlerBeingBuilt := builder.logHandler
	if builder.buildErrors.Size() == 0 {
		return handlerBeingBuilt, nil
	} else {
		return handlerBeingBuilt, builder.buildErrors
	}
}
