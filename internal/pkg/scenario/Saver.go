// Copyright (c) 2019 Australian Rivers Institute.

package scenario

import (
	"os"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/pkg/errors"
)

const defaultOutputPath = "solutions"

type CallableSaver interface {
	observer.Observer
	logging.Container
}

type Saver struct {
	loggers.LoggerContainer
	outputType encoding.OutputType
	outputPath string
}

func NewSaver() *Saver {
	saver := new(Saver).WithOutputPath(defaultOutputPath)
	return saver
}

func (s *Saver) WithOutputType(outputType encoding.OutputType) *Saver {
	s.outputType = outputType
	return s
}

func (s *Saver) WithOutputPath(outputPath string) *Saver {
	if outputPath == "" {
		outputPath = defaultOutputPath
	}

	s.outputPath = outputPath
	return s
}

func (s *Saver) ensureOutputPathIsUsable() {
	if fileInfo, err := os.Stat(s.outputPath); err == nil {
		s.ensureExistingOutputPathIsUsable(fileInfo)
	} else if os.IsNotExist(err) {
		s.createOutputPath()
	} else {
		panic(errors.Wrap(err, "scenario saver cannot get file info of output path"))
	}
}

func (s *Saver) ensureExistingOutputPathIsUsable(fileInfo os.FileInfo) {
	if !fileInfo.IsDir() {
		panic(errors.New("scenario saver output path specified not a directory"))
	}
}

func (s *Saver) createOutputPath() {
	if mkDirError := os.MkdirAll(s.outputPath, os.ModePerm); mkDirError != nil {
		panic(errors.New("scenario saver failed making output path specified"))
	}
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
	s.debugLogSolutionInJson(modelSolution)
	s.ensureOutputPathIsUsable()
	s.encode(modelSolution)
}

func (s *Saver) encode(modelSolution solution.Solution) {
	encoder := new(encoding.Builder).
		ForOutputType(s.outputType).
		WithOutputPath(s.outputPath).
		Build()

	if encodingError := encoder.Encode(&modelSolution); encodingError != nil {
		s.LogHandler().Error(encodingError)
	}
}

func (s *Saver) debugLogSolutionInJson(modelSolution solution.Solution) {
	s.LogHandler().Debug("JSON Encoding of solution after annealing finished:")

	modelSolutionAsJson := s.toJson(&modelSolution)
	s.LogHandler().Debug(modelSolutionAsJson)
}

func (s *Saver) toJson(structure *solution.Solution) string {
	marshaler := new(json.Marshaler)
	solutionAsJson, marshalError := marshaler.Marshal(structure)

	if marshalError != nil {
		panic(errors.Wrap(marshalError, "failed marshalling solution to JSON"))
	}

	return string(solutionAsJson)
}
