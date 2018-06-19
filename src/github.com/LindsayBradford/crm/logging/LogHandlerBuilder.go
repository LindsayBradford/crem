// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	crmerrors "github.com/LindsayBradford/crm/errors"
)

type LogHandlerBuilder struct{
	logHandler  LogHandler
	buildErrors *crmerrors.CompositeError
}

func (this *LogHandlerBuilder) ForNativeLibraryLogHandler() *LogHandlerBuilder {
	this.buildErrors = crmerrors.NewComposite("Failed to build valid LogHandler")

	newHandler := new(NativeLibraryLogHandler)

	defaultDestinations := new(LogLevelDestinations).Initialise()
	newHandler.SetDestinations(defaultDestinations)

	newHandler.Initialise().WithFormatter(new (NullFormatter))

	this.logHandler = newHandler
	return this
}

func (this *LogHandlerBuilder) WithFormatter(formatter LogAttributeFormatter) *LogHandlerBuilder {
	formatter.Initialise()
	this.logHandler.SetFormatter(formatter)
	return this
}
func (this *LogHandlerBuilder) WithLogLevelDestination(logLevel LogLevel, destination LogDestination) *LogHandlerBuilder {
	this.logHandler.Destinations().Override(logLevel, destination)
	return this
}

func (this *LogHandlerBuilder) Build() (LogHandler, error) {
	if this.buildErrors.Size() == 0 {
		return this.logHandler, nil
	} else {
		return this.logHandler, this.buildErrors
	}
}



