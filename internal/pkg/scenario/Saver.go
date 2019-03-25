// Copyright (c) 2019 Australian Rivers Institute.

package scenario

import (
	"encoding/json"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

type CallableSaver interface {
	observer.Observer
	logging.Container
}

type Saver struct {
	loggers.LoggerContainer
}

func (s *Saver) ObserveEvent(event observer.Event) {
	if observableAnnealer, isAnnealer := event.Source().(annealing.Observable); isAnnealer {
		s.observeAnnealingEvent(observableAnnealer, event)
	}
}

func (s *Saver) observeAnnealingEvent(annealer annealing.Observable, event observer.Event) {
	if event.EventType != observer.FinishedAnnealing {
		return
	}

	s.saveModelSolution(annealer)
}

func (s *Saver) saveModelSolution(annealer annealing.Observable) {
	modelSolution := annealer.Solution()

	modelSolutionAsJson := s.toJson(&modelSolution)

	s.LogHandler().Debug(modelSolutionAsJson)

	// TODO: actaully save the solution state
}

func (s *Saver) toJson(structure *solution.Solution) string {
	structureJson, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		s.LogHandler().LogAtLevel(logging.ERROR, err)
		return "{}"
	}
	return string(structureJson)

}
