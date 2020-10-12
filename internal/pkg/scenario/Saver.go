// Copyright (c) 2019 Australian Rivers Institute.

package scenario

import (
	"fmt"
	solutionset "github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	encoding2 "github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding"
	"os"

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
		solution := event.Attribute(Solution).(solution.Solution)
		s.saveSolution(solution)
	}
	if event.HasAttribute(SolutionSet) {
		solutionSet := event.Attribute(SolutionSet).(archive.NonDominanceModelArchive)
		s.saveSolutionSet(solutionSet)
	}
}

func (s *Saver) saveSolution(solution solution.Solution) {
	s.debugLogSolutionInJson(solution)
	s.ensureOutputPathIsUsable()
	s.encodeSolution(solution)
}

func (s *Saver) encodeSolution(modelSolution solution.Solution) {
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

func (s *Saver) saveSolutionSet(solutionSet archive.NonDominanceModelArchive) {
	// s.debugLogSolutionInJson(solution)
	s.ensureOutputPathIsUsable()
	s.encodeSolutionSet(solutionSet)
}

func (s *Saver) encodeSolutionSet(solutionSet archive.NonDominanceModelArchive) {
	summary := make(solutionset.Summary, 0)
	for solutionIndex, compressedModel := range solutionSet.Archive() {
		solution := s.deriveModelSolution(solutionSet, solutionIndex, compressedModel)
		s.encodeSolution(*solution)
		s.summarise(&summary, solution)
	}
	s.encodeSummary(&summary)
}

func (s *Saver) deriveModelSolution(solutionSet archive.NonDominanceModelArchive, solutionIndex int, compressedModel *archive.CompressedModelState) *solution.Solution {
	solutionId := s.deriveSolutionId(solutionSet, solutionIndex+1)
	return s.deriveSolutionFrom(solutionSet, compressedModel, solutionId)
}

func (s *Saver) deriveSolutionFrom(solutionSet archive.NonDominanceModelArchive, compressedModel *archive.CompressedModelState, solutionId string) *solution.Solution {
	decompressedModel := s.decompress(solutionSet, compressedModel)

	decompressedModelSolution := new(solution.SolutionBuilder).
		WithId(solutionId).
		ForModel(decompressedModel).
		Build()

	return decompressedModelSolution
}

func (s *Saver) decompress(solutionSet archive.NonDominanceModelArchive, compressedModel *archive.CompressedModelState) model.Model {
	decompressedModel := s.decompressionModel.DeepClone()
	decompressedModel.Initialise()
	solutionSet.Decompress(compressedModel, decompressedModel)
	return decompressedModel
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
		Build()

	if encodingError := encoder.Encode(summary); encodingError != nil {
		s.LogHandler().Error(encodingError)
	}
}
