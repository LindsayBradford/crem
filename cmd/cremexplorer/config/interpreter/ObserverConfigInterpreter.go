// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	data2 "github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
	annealingObserver "github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/config/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type ObserverConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	loggingInterpreter *interpreter.LoggingConfigInterpreter

	observer observer.Observer
}

func NewObserverConfigInterpreter() *ObserverConfigInterpreter {
	interpreter := new(ObserverConfigInterpreter).initialise()
	return interpreter
}

func (i *ObserverConfigInterpreter) initialise() *ObserverConfigInterpreter {
	i.errors = compositeErrors.New("Observer Configuration")
	i.loggingInterpreter = interpreter.NewLoggingConfigInterpreter()
	i.initialiseObserving()
	return i
}

func (i *ObserverConfigInterpreter) initialiseObserving() {
	i.observer = new(annealingObserver.AnnealingMessageObserver).
		WithLogHandler(i.LogHandler()).
		WithFilter(new(filters.IterationCountFilter).WithModulo(1))
}

func (i *ObserverConfigInterpreter) Interpret(config *data2.ObserverConfig) *ObserverConfigInterpreter {
	i.interpretLogger(&config.LoggingConfig)
	i.interpretObserver(config)
	return i
}

func (i *ObserverConfigInterpreter) interpretLogger(config *data.LoggingConfig) {
	i.loggingInterpreter.Interpret(config)
	if i.loggingInterpreter.Errors() != nil {
		i.errors.Add(i.loggingInterpreter.Errors())
	}
}

func (i *ObserverConfigInterpreter) interpretObserver(config *data2.ObserverConfig) {
	i.observer = new(annealingObserver.AnnealingMessageObserver).
		WithLogHandler(i.LogHandler()).
		WithFilter(new(filters.IterationCountFilter).WithModulo(1))
}

func (i *ObserverConfigInterpreter) Observer() observer.Observer {
	return i.observer
}

func (i *ObserverConfigInterpreter) LogHandler() logging.Logger {
	return i.loggingInterpreter.LogHandler()
}

func (i *ObserverConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
