// Copyright (c) 2019 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
)

type NullAnnealer struct {
	name.IdentifiableContainer

	explorer.ContainedExplorer
	model.ContainedModel

	loggers.ContainedLogger

	observer.ContainedEventNotifier
}

func (sa *NullAnnealer) Initialise() {
	sa.SetSolutionExplorer(null.NullExplorer)
	sa.SetLogHandler(new(loggers.NullLogger))
	sa.SetEventNotifier(new(observer.SynchronousAnnealingEventNotifier))

	sa.SetId("Null Annealer")
}

func (sa *NullAnnealer) SetId(title string) {
	sa.IdentifiableContainer.SetId(title)
}

func (sa *NullAnnealer) DeepClone() annealing.Annealer {
	clone := *sa
	explorerClone := sa.SolutionExplorer().DeepClone()
	clone.SetSolutionExplorer(explorerClone)
	return &clone
}

func (sa *NullAnnealer) EventAttributes(eventType observer.EventType) attributes.Attributes {
	return nil
}
func (sa *NullAnnealer) SetParameters(params parameters.Map) error    { return nil }
func (sa *NullAnnealer) ParameterErrors() error                       { return nil }
func (sa *NullAnnealer) Model() model.Model                           { return model.NullModel }
func (sa *NullAnnealer) Anneal()                                      {}
func (sa *NullAnnealer) AddObserver(observer observer.Observer) error { return nil }
func (sa *NullAnnealer) Observers() []observer.Observer               { return nil }
