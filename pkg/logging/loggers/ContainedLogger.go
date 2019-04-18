// Copyright (c) 2019 Australian Rivers Institute.

package loggers

import "github.com/LindsayBradford/crem/pkg/logging"

// ContainedLogger is a struct offering a default container implementation of a Log Handler
type ContainedLogger struct {
	logHandler logging.Logger
}

func (c *ContainedLogger) SetLogHandler(logHandler logging.Logger) {
	if logHandler == nil {
		logHandler = NewNullLogger()
	}
	c.logHandler = logHandler
}

func (c *ContainedLogger) LogHandler() logging.Logger {
	return c.logHandler
}
