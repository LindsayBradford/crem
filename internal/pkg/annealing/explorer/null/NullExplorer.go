// Copyright (c) 2018 Australian Rivers Institute.

package null

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

var NullExplorer = new(Explorer)

type Explorer struct {
	name.ContainedName
}

func (e *Explorer) WithName(name string) *Explorer {
	e.SetName(name)
	return e
}

func (e *Explorer) Initialise()                                   {}
func (e *Explorer) TearDown()                                     {}
func (e *Explorer) SetObjectiveValue(temperature float64)         {}
func (e *Explorer) TryRandomChange(temperature float64)           {}
func (e *Explorer) AcceptLastChange()                             {}
func (e *Explorer) RevertLastChange()                             {}
func (e *Explorer) ScenarioId() string                            { return "" }
func (e *Explorer) SetScenarioId(id string)                       {}
func (e *Explorer) DeepClone() explorer.Explorer                  { return e }
func (e *Explorer) CloneObservable() explorer.Explorer            { return e }
func (e *Explorer) Model() model.Model                            { return nil }
func (e *Explorer) SetModel(model model.Model)                    {}
func (e *Explorer) SetLogHandler(logHandler logging.Logger) error { return nil }
func (e *Explorer) LogHandler() logging.Logger                    { return nil }
func (e *Explorer) SetRandomNumberGenerator(generator *rand.Rand) {}
func (e *Explorer) RandomNumberGenerator() *rand.Rand             { return nil }

func (e *Explorer) ObjectiveValue() float64         { return 0 }
func (e *Explorer) ChangeInObjectiveValue() float64 { return 0 }
func (e *Explorer) AcceptanceProbability() float64  { return 0 }
func (e *Explorer) ChangeIsDesirable() bool         { return false }
func (e *Explorer) ChangeAccepted() bool            { return false }