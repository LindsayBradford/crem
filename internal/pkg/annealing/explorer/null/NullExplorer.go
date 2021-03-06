// Copyright (c) 2018 Australian Rivers Institute.

package null

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
)

var NullExplorer = new(Explorer)

type Explorer struct {
	name.NameContainer
	name.IdentifiableContainer
}

func (e *Explorer) WithName(name string) *Explorer {
	e.SetName(name)
	return e
}

func (e *Explorer) DeepClone() explorer.Explorer              { return e }
func (e *Explorer) Initialise()                               {}
func (e *Explorer) TearDown()                                 {}
func (e *Explorer) SetParameters(params parameters.Map) error { return nil }
func (e *Explorer) ParameterErrors() error                    { return nil }

func (e *Explorer) TryRandomChange() {}
func (e *Explorer) EventAttributes(eventType observer.EventType) attributes.Attributes {
	return nil
}
func (e *Explorer) CoolDown() {}

func (e *Explorer) Model() model.Model                  { return model.NullModel }
func (e *Explorer) SetModel(model model.Model)          {}
func (e *Explorer) LogHandler() logging.Logger          { return new(loggers.NullLogger) }
func (e *Explorer) SetLogHandler(logger logging.Logger) {}
