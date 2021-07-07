package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/archive"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"time"
)

var solutionBuilder solution.SolutionBuilder
var modelCompressor archive.ModelCompressor

type SolutionPoolLabel string

const (
	AsIs SolutionPoolLabel = "As-Is"
)

func NewSolutionPool(referenceModel *catchment.Model) SolutionPool {
	newPool := SolutionPool{cache: make(map[SolutionPoolLabel]SolutionContainer, 1)}
	newPool.initialise(referenceModel)
	return newPool
}

type SolutionPool struct {
	referenceModel *catchment.Model
	builder        solution.SolutionBuilder
	cache          map[SolutionPoolLabel]SolutionContainer
}

func (sp *SolutionPool) initialise(referenceModel *catchment.Model) {
	sp.assignAsIsSolutionFrom(referenceModel)
}

func (sp *SolutionPool) Size() int {
	return len(sp.cache)
}

func (sp *SolutionPool) assignAsIsSolutionFrom(referenceModel *catchment.Model) {
	asIsModel := generateAsIsModel(referenceModel)
	sp.referenceModel = asIsModel
	asIsSolution := sp.deriveSolutionFrom(asIsModel)
	sp.cache[AsIs] = NewSolutionContainer(asIsSolution, "As-Is solution. No management actions active.")
}

func (sp *SolutionPool) deriveSolutionFrom(model *catchment.Model) *solution.Solution {
	return sp.builder.WithId(model.Id()).ForModel(model).Build()
}

func generateAsIsModel(referenceModel *catchment.Model) *catchment.Model {
	asIsClone := referenceModel.DeepClone()
	asIsClone.Initialise(model.AsIs)
	return toCatchmentModel(asIsClone)
}

func toCatchmentModel(thisModel model.Model) *catchment.Model {
	catchmentModel, isCatchmentModel := thisModel.(*catchment.Model)
	if isCatchmentModel {
		return catchmentModel
	}
	assert.That(false).WithFailureMessage("Should not get here").Holds()
	return nil
}

func (sp *SolutionPool) HasSolution(label SolutionPoolLabel) bool {
	if _, hasContainer := sp.cache[label]; hasContainer {
		return true
	}
	return false
}

func (sp *SolutionPool) Solution(label SolutionPoolLabel) *solution.Solution {
	return sp.cache[label].Solution
}

func (sp *SolutionPool) Summary(label SolutionPoolLabel) string {
	return sp.cache[label].Summary
}

func (sp *SolutionPool) AddSolution(label SolutionPoolLabel, modelEncoding string, summary string) {
	if label == AsIs {
		return
	}

	newModel := sp.referenceModel.DeepClone()

	compressedModel := modelCompressor.Compress(newModel)
	compressedModel.Decode(modelEncoding)
	modelCompressor.Decompress(compressedModel, newModel)

	newCatchmentModel := toCatchmentModel(newModel)
	newSolution := sp.deriveSolutionFrom(newCatchmentModel)
	sp.cache[label] = NewSolutionContainer(newSolution, summary)
}

func NewSolutionContainer(solution *solution.Solution, summary string) SolutionContainer {
	container := SolutionContainer{Solution: solution, Summary: summary, LastUpdated: time.Now()}
	return container
}

type SolutionContainer struct {
	Solution    *solution.Solution
	Summary     string
	LastUpdated time.Time
}
