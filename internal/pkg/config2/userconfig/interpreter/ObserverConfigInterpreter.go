// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	annealingObserver "github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type ObserverConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	loggingInterpreter *LoggingConfigInterpreter

	observer observer.Observer
}

func NewObserverConfigInterpreter() *ObserverConfigInterpreter {
	interpreter := new(ObserverConfigInterpreter).initialise()
	return interpreter
}

func (i *ObserverConfigInterpreter) initialise() *ObserverConfigInterpreter {
	i.errors = compositeErrors.New("Observer Configuration")
	i.loggingInterpreter = NewLoggingConfigInterpreter().initialise()
	i.initialiseObserving()
	return i
}

func (i *ObserverConfigInterpreter) initialiseObserving() {
	i.observer = new(annealingObserver.AnnealingMessageObserver).
		WithLogHandler(i.LogHandler()).
		WithFilter(new(filters.IterationCountFilter).WithModulo(1))
}

func (i *ObserverConfigInterpreter) Interpret(config *data.ObserverConfig) *ObserverConfigInterpreter {
	i.loggingInterpreter.Interpret(&config.LoggingConfig)
	i.interpretObserverSpecific(config)
	return i
}

func (i *ObserverConfigInterpreter) interpretObserverSpecific(config *data.ObserverConfig) {
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
