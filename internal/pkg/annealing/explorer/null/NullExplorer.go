// Copyright (c) 2018 Australian Rivers Institute.

package null

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

var NullExplorer = new(Explorer)

type Explorer struct {
	explorer.SingleObjectiveAnnealableExplorer
}

func (e *Explorer) Initialise() {
	e.SetObjectiveValue(0)
}

func (e *Explorer) WithName(name string) *Explorer {
	e.SetName(name)
	return e
}

func (e *Explorer) SetObjectiveValue(temperature float64) {}
func (e *Explorer) TryRandomChange(temperature float64)   {}
func (e *Explorer) AcceptLastChange()                     {}
func (e *Explorer) RevertLastChange()                     {}
