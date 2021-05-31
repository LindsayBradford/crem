// Copyright (c) 2019 Australian Rivers Institute.

package scenario

import (
	"fmt"
	solutionset "github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	encoding2 "github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding"
	"os"
	"sync"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/archive"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/pkg/errors"
)

const (
	Solution          = "Solution"
	SolutionSet       = "SolutionSet"
	defaultOutputPath = "solutions"
)

type CallableSaver interface {
	observer.Observer
	logging.Container
	SetDecompressionModel(model model.Model)
}

type Saver struct {
	loggers.ContainedLogger
	decompressionModel model.Model
	outputType         encoding.OutputType
	outputPath         string

	decompressionMutex sync.Mutex
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

func (s *Saver) WithLogHandler(logHandler logging.Logger) *Saver {
	s.SetLogHandler(logHandler)
	return s
}

func (s *Saver) SetDecompressionModel(model model.Model) {
	clone := model.DeepClone()
	clone.Initialise()
	s.decompressionModel = clone
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
	if event.EventType != observer.FinishedAnnealing {
		return
	}
	if event.HasAttribute(Solution) {
		s.LogHandler().Info("Saving annealing optimised solution")
		solution := event.Attribute(Solution).(solution.Solution)
		s.saveOptimisedSolution(&solution)
	}
	if event.HasAttribute(SolutionSet) {
		s.LogHandler().Info("Saving annealing solution set")
		solutionSet := event.Attribute(SolutionSet).(archive.NonDominanceModelArchive)
		s.saveSolutionSet(solutionSet)
	}
}

func (s *Saver) saveOptimisedSolution(solution *solution.Solution) {
	s.ensureOutputPathIsUsable()
	s.encodeOptimisedSolution(solution)
}

func (s *Saver) encodeOptimisedSolution(optimisedSolution *solution.Solution) {
	summary := make(solutionset.Summary, 0)

	asIsSolution := s.deriveASsIsSolutionForOptimised(optimisedSolution)
	s.encodeSolution(*asIsSolution)
	s.summarise(&summary, asIsSolution)

	optimisedSolution.Id = optimisedSolution.Id + " Solution (1/1)" // Necessary for summary parsing later.

	s.encodeSolution(*optimisedSolution)
	s.summarise(&summary, optimisedSolution)
	s.encodeSummary(&summary)
}

func (s *Saver) deriveASsIsSolutionForOptimised(optimisedSolution *solution.Solution) *solution.Solution {
	s.decompressionMutex.Lock()
	defer s.decompressionMutex.Unlock()

	s.decompressionModel.Initialise()
	asIsSolutionId := s.deriveAsIsOptimisedSolutionId(optimisedSolution)

	decompressedModelSolution := new(solution.SolutionBuilder).
		WithId(asIsSolutionId).
		ForModel(s.decompressionModel).
		Build()

	return decompressedModelSolution
}

func (s *Saver) deriveAsIsOptimisedSolutionId(optimisedSolution *solution.Solution) string {
	return optimisedSolution.Id + " Solution (As-Is)"
}

func (s *Saver) encodeSolution(modelSolution solution.Solution) {
	encoder := new(encoding.Builder).
		ForOutputType(s.outputType).
		WithOutputPath(s.outputPath).
		WithLogHandler(s.LogHandler()).
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

func (s *Saver) saveSolutionSet(solutionSet archive.NonDominanceModelArchive) {
	s.ensureOutputPathIsUsable()
	s.encodeSolutionSet(solutionSet)
}

func (s *Saver) encodeSolutionSet(solutionSet archive.NonDominanceModelArchive) {
	summary := make(solutionset.Summary, 0)

	asIsSolution := s.deriveASsIsSolution(solutionSet)
	s.encodeSolution(*asIsSolution)
	s.summarise(&summary, asIsSolution)

	for solutionIndex, compressedModel := range solutionSet.Archive() {
		solution := s.deriveModelSolution(solutionSet, solutionIndex, compressedModel)
		s.encodeSolution(*solution)
		s.summarise(&summary, solution)
	}
	s.encodeSummary(&summary)
}

func (s *Saver) deriveASsIsSolution(solutionSet archive.NonDominanceModelArchive) *solution.Solution {
	s.decompressionMutex.Lock()
	defer s.decompressionMutex.Unlock()

	s.decompressionModel.Initialise()
	asIsSolutionId := s.deriveAsIsSolutionId(solutionSet)

	decompressedModelSolution := new(solution.SolutionBuilder).
		WithId(asIsSolutionId).
		ForModel(s.decompressionModel).
		Build()

	return decompressedModelSolution
}

func (s *Saver) deriveAsIsSolutionId(solutionSet archive.NonDominanceModelArchive) string {
	return solutionSet.Id() + " Solution (As-Is)"
}

func (s *Saver) deriveModelSolution(solutionSet archive.NonDominanceModelArchive, solutionIndex int, compressedModel *archive.CompressedModelState) *solution.Solution {
	solutionId := s.deriveSolutionId(solutionSet, solutionIndex+1)
	return s.deriveSolutionFrom(solutionSet, compressedModel, solutionId)
}

func (s *Saver) deriveSolutionFrom(solutionSet archive.NonDominanceModelArchive, compressedModel *archive.CompressedModelState, solutionId string) *solution.Solution {
	s.decompressionMutex.Lock()
	defer s.decompressionMutex.Unlock()

	solutionSet.Decompress(compressedModel, s.decompressionModel)

	decompressedModelSolution := new(solution.SolutionBuilder).
		WithId(solutionId).
		ForModel(s.decompressionModel).
		Build()

	return decompressedModelSolution
}

func (s *Saver) deriveSolutionId(solutionSet archive.NonDominanceModelArchive, currentSolution int) string {
	solutionSetSize := solutionSet.Len()
	solutionId := fmt.Sprintf("%s Solution (%d/%d)", solutionSet.Id(), currentSolution, solutionSetSize)
	return solutionId
}

func (s *Saver) summarise(summary *solutionset.Summary, solution *solution.Solution) {
	baseMap := *summary
	baseMap[solution.Id] = solution.Summarise()
}

func (s *Saver) encodeSummary(summary *solutionset.Summary) {
	encoder := new(encoding2.Builder).
		ForOutputType(s.outputType).
		WithOutputPath(s.outputPath).
		WithLogHandler(s.LogHandler()).
		Build()

	if encodingError := encoder.Encode(summary); encodingError != nil {
		s.LogHandler().Error(encodingError)
	}
}
