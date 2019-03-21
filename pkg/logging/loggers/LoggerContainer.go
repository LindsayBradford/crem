// Copyright (c) 2019 Australian Rivers Institute.

package loggers

import "github.com/LindsayBradford/crem/pkg/logging"

// LoggerContainer is a struct offering a default container implementation of a Log Handler
type LoggerContainer struct {
	logHandler logging.Logger
}

func (c *LoggerContainer) SetLogHandler(logHandler logging.Logger) {
	if logHandler == nil {
		logHandler = NewNullLogger()
	}
	c.logHandler = logHandler
}

func (c *LoggerContainer) LogHandler() logging.Logger {
	return c.logHandler
}
