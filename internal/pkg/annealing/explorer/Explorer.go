// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"errors"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Explorer interface {
	Observable
	name.Nameable

	ScenarioId() string
	SetScenarioId(name string)

	Initialise()
	TearDown()

	DeepClone() Explorer
	CloneObservable() Explorer
	TryRandomChange(temperature float64)

	rand.Container
	model.Container
	logging.Container
}

type Observable interface {
	ObjectiveValue() float64
	ChangeInObjectiveValue() float64
	AcceptanceProbability() float64
	ChangeIsDesirable() bool
	ChangeAccepted() bool
}

// Container defines an interface embedding a Model
type Container interface {
	SolutionExplorer() Explorer
	SetSolutionExplorer(explorer Explorer) error
}

// ContainedExplorer is a struct offering a default implementation of Container
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
