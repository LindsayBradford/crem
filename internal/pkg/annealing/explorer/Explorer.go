// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"errors"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

const Guaranteed = 1

const (
	ChangeIsDesirable = "ChangeIsDesirable"

	Temperature   = "Temperature"
	CoolingFactor = "CoolingFactor"

	AcceptanceProbability = "AcceptanceProbability"
	ChangeAccepted        = "ChangeAccepted"
	ChangeInvalid         = "ChangeInvalid"
	ReasonChangeInvalid   = "ReasonChangeInvalid"
)

type Explorer interface {
	name.Nameable
	name.Identifiable
	model.Container
	parameters.Container
	logging.Container

	DeepClone() Explorer
	Initialise()
	TearDown()

	TryRandomChange()

	CoolDown()
	EventAttributes(eventType observer.EventType) attributes.Attributes
}

// Container defines an interface embedding an Explorer
type Container interface {
	SolutionExplorer() Explorer
	SetSolutionExplorer(explorer Explorer) error
}

// ContainedExplorer is a struct offering a default implementation of ContainedLogger
type ContainedExplorer struct {
	explorer Explorer
}

func (e *ContainedExplorer) SolutionExplorer() Explorer {
	return e.explorer
}

func (e *ContainedExplorer) SetSolutionExplorer(explorer Explorer) error {
	if explorer == nil {
		return errors.New("invalid attempt to set Solution Explorer to nil value")
	}
	e.explorer = explorer
	return nil
}
