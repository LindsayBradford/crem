// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	appData "github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
	annealingObserver "github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/config/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
)

const defaultReportingIterationNumber = 1

type ReportingConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	loggingInterpreter *interpreter.LoggingConfigInterpreter

	observer observer.Observer
}

func NewObserverConfigInterpreter() *ReportingConfigInterpreter {
	interpreter := new(ReportingConfigInterpreter).initialise()
	return interpreter
}

func (i *ReportingConfigInterpreter) initialise() *ReportingConfigInterpreter {
	i.errors = compositeErrors.New("Reporting Configuration")
	i.loggingInterpreter = interpreter.NewLoggingConfigInterpreter()
	i.initialiseObserving()
	return i
}

func (i *ReportingConfigInterpreter) initialiseObserving() {
	i.observer = new(annealingObserver.AnnealingMessageObserver).
		WithLogHandler(i.LogHandler()).
		WithFilter(new(filters.IterationCountFilter).WithModulo(defaultReportingIterationNumber))
}

func (i *ReportingConfigInterpreter) Interpret(config *appData.ReportingConfig) *ReportingConfigInterpreter {
	i.interpretLogger(&config.LoggingConfig)
	i.interpretObserver(config)
	return i
}

func (i *ReportingConfigInterpreter) interpretLogger(config *data.LoggingConfig) {
	i.loggingInterpreter.Interpret(config)
	if i.loggingInterpreter.Errors() != nil {
		i.errors.Add(i.loggingInterpreter.Errors())
	}
}

func (i *ReportingConfigInterpreter) interpretObserver(config *appData.ReportingConfig) {
	i.observer = new(annealingObserver.AnnealingMessageObserver).
		WithLogHandler(i.LogHandler()).
		WithFilter(
			new(filters.IterationCountFilter).
				WithModulo(config.ReportEveryNumberOfIterations),
		)
}

func (i *ReportingConfigInterpreter) Observer() observer.Observer {
	return i.observer
}

func (i *ReportingConfigInterpreter) LogHandler() logging.Logger {
	return i.loggingInterpreter.LogHandler()
}

func (i *ReportingConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
