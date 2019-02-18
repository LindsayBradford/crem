// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

const Guaranteed = 1

type Explorer interface {
	name.Nameable
	name.Identifiable
	model.Container
	rand.Container

	Observable
	logging.Container

	DeepClone() Explorer
	Initialise()
	TearDown()

	TryRandomChange(temperature float64)
	CloneObservable() Observable
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

// Container defines an interface embedding a Model
type ObservableContainer interface {
	ObservableExplorer() Observable
	SetObservableExplorer(explorer Observable) error
}

// ContainedExplorer is a struct offering a default implementation of Container
type ContainedObservableExplorer struct {
	explorer Observable
}

func (e *ContainedObservableExplorer) ObservableExplorer() Observable {
	return e.explorer
}

func (e *ContainedObservableExplorer) SetObservableExplorer(explorer Observable) error {
	if explorer == nil {
		return errors.New("invalid attempt to set Solution Explorer to nil value")
	}
	e.explorer = explorer
	return nil
}
