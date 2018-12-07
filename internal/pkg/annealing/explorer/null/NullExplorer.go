// Copyright (c) 2018 Australian Rivers Institute.

package null

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

var NullExplorer = New()

type Explorer struct {
	explorer.BaseExplorer
	logging.ContainedLogger
}

func New() *Explorer {
	newExplorer := new(Explorer)
	newExplorer.Initialise()
	return newExplorer
}

func (e *Explorer) Initialise() {
	e.SetObjectiveValue(0)
	e.BaseExplorer.SetLogHandler(loggers.DefaultNullLogger)
}

func (e *Explorer) WithName(name string) *Explorer {
	e.SetName(name)
	return e
}

func (e *Explorer) SetObjectiveValue(temperature float64) {}
func (e *Explorer) TryRandomChange(temperature float64)   {}
func (e *Explorer) AcceptLastChange()                     {}
func (e *Explorer) RevertLastChange()                     {}
