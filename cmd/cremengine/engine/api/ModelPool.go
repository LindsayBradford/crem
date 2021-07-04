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

type ModelPoolLabel string

const (
	AsIs       ModelPoolLabel = "As-Is"
	Scratchpad ModelPoolLabel = "Scratchpad"
)

func NewModelPool(referenceModel *catchment.Model) ModelPool {
	newPool := ModelPool{cache: make(map[ModelPoolLabel]ModelContainer, 2)}
	newPool.initialise(referenceModel)
	return newPool
}

type ModelPool struct {
	cache map[ModelPoolLabel]ModelContainer
}

func (mp *ModelPool) initialise(referenceModel *catchment.Model) {
	mp.assignAsIsModelFrom(referenceModel)
	mp.generateScratchpadModel()
}

func (mp *ModelPool) Size() int {
	return len(mp.cache)
}

func (mp *ModelPool) assignAsIsModelFrom(referenceModel *catchment.Model) {
	asIsModel := generateAsIsModel(referenceModel)
	asIsModel.AddAttribute("Summary", "As-Is Model. No active management actions.")
	mp.cache[AsIs] = NewModelContainer(asIsModel)
}

func generateAsIsModel(referenceModel *catchment.Model) *catchment.Model {
	asIsClone := referenceModel.DeepClone()

	return toCatchmentModel(asIsClone)
}

func (mp *ModelPool) generateScratchpadModel() {
	asIsModel := mp.Model(AsIs)

	scratchpadModel := asIsModel.DeepClone()

	scratchpadCatchmentModel := toCatchmentModel(scratchpadModel)
	scratchpadCatchmentModel.AddAttribute("Summary", "Scratchpad")
	mp.cache[Scratchpad] = NewModelContainer(scratchpadCatchmentModel)
}

func toCatchmentModel(thisModel model.Model) *catchment.Model {
	catchmentModel, isCatchmentModel := thisModel.(*catchment.Model)
	if isCatchmentModel {
		return catchmentModel
	}
	assert.That(false).WithFailureMessage("Should not get here").Holds()
	return nil
}

func (mp *ModelPool) HasModel(label ModelPoolLabel) bool {
	if _, hasContainer := mp.cache[label]; hasContainer {
		return true
	}
	return false
}

func (mp *ModelPool) Model(label ModelPoolLabel) *catchment.Model {
	return mp.cache[label].Model
}

func (mp *ModelPool) Solution(label ModelPoolLabel) *solution.Solution {
	return mp.cache[label].Solution
}

func (mp *ModelPool) Update(label ModelPoolLabel) {
	containerToUpdate := mp.cache[label]
	containerToUpdate.Update()
}

func (mp *ModelPool) InstantiateModel(label ModelPoolLabel, modelEncoding string, summary string) {
	if label == AsIs {
		return
	}

	newModel := mp.Model(AsIs).DeepClone()

	compressedModel := modelCompressor.Compress(newModel)
	compressedModel.Decode(modelEncoding)
	modelCompressor.Decompress(compressedModel, newModel)

	newCatchmentModel := toCatchmentModel(newModel)
	newCatchmentModel.AddAttribute("Summary", summary)
	mp.cache[label] = NewModelContainer(newCatchmentModel)
}

func NewModelContainer(model *catchment.Model) ModelContainer {
	container := ModelContainer{Model: model, LastUpdated: time.Now()}
	container.Update()
	return container
}

type ModelContainer struct {
	Model    *catchment.Model
	Solution *solution.Solution

	LastUpdated time.Time
}

func (mc *ModelContainer) Update() {
	mc.Solution = solutionBuilder.
		WithId(mc.Model.Id()).
		ForModel(mc.Model).
		Build()
	mc.LastUpdated = time.Now()
}
