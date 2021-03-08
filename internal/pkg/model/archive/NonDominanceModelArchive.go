// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"github.com/LindsayBradford/crem/pkg/name"
	"math"
	"sort"
)

type StorageResult uint

const (
	StoredReplacingDominatedEntries StorageResult = iota
	StoredWithNoDominanceDetected
	RejectedWithStoredEntryDominanceDetected
	RejectedWithDuplicateEntryDetected
	canBeStored
	StoredForcingDominatingStateRemoval
)

func (sr StorageResult) String() string {
	switch sr {
	case StoredWithNoDominanceDetected:
		return "Stored, no dominance of archive entries detected"
	case StoredReplacingDominatedEntries:
		return "Stored, replacing dominated solutions in archive"
	case RejectedWithStoredEntryDominanceDetected:
		return "Rejected, solution(s) in archive would dominate"
	case RejectedWithDuplicateEntryDetected:
		return "Rejected, solution is already archived"
	case StoredForcingDominatingStateRemoval:
		return "Stored, forcing dominating solutions out of archive"
	}
	return "Can be stored -- but why are you seeing this?!?"
}

func New() *NonDominanceModelArchive {
	return new(NonDominanceModelArchive).Initialise()
}

const notIsolated = math.MaxFloat64

type NonDominanceModelArchive struct {
	name.IdentifiableContainer
	archive   []*CompressedModelState
	isolation []float64

	compressor ModelCompressor
	rand.RandContainer
}

func (a *NonDominanceModelArchive) Initialise() *NonDominanceModelArchive {
	a.archive = make([]*CompressedModelState, 0)
	a.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return a
}

func (a *NonDominanceModelArchive) AttemptToArchive(model model.Model) StorageResult {
	compressedModelState := a.compressor.Compress(model)
	return a.AttemptToArchiveState(compressedModelState)
}

func (a *NonDominanceModelArchive) Archive() []*CompressedModelState {
	return a.archive
}

func (a *NonDominanceModelArchive) Compress(model model.Model) *CompressedModelState {
	return a.compressor.Compress(model)
}

func (a *NonDominanceModelArchive) Decompress(condensedModelState *CompressedModelState, model model.Model) {
	assert.That(condensedModelState != nil && model != nil).
		WithFailureMessage("Model state, or model is nil").
		Holds()
	a.compressor.Decompress(condensedModelState, model)
}

func (a *NonDominanceModelArchive) AttemptToArchiveState(modelState *CompressedModelState) StorageResult {
	storageState := a.newModelStateCannotBeArchived(modelState)

	if storageState != canBeStored {
		return storageState
	}

	storageState = StoredWithNoDominanceDetected

	nonDominatedArray := make([]*CompressedModelState, 0)
	for currentIndex := 0; currentIndex < len(a.archive); currentIndex++ {
		if modelState.Variables.Dominates(&a.archive[currentIndex].Variables) {
			storageState = StoredReplacingDominatedEntries
		} else {
			nonDominatedArray = append(nonDominatedArray, a.archive[currentIndex])
		}
	}

	if storageState == StoredReplacingDominatedEntries {
		a.archive = nonDominatedArray
	}

	a.archive = append(a.archive, modelState)
	return storageState
}

func (a *NonDominanceModelArchive) ForceIntoArchive(model model.Model) StorageResult {
	compressedModelState := a.compressor.Compress(model)
	return a.ForceModelStateIntoArchive(compressedModelState)
}

func (a *NonDominanceModelArchive) ForceModelStateIntoArchive(modelState *CompressedModelState) StorageResult {
	nonDominatedArray := make([]*CompressedModelState, 0)
	for currentIndex := 0; currentIndex < len(a.archive); currentIndex++ {
		currentState := &a.archive[currentIndex].Variables
		if !currentState.Dominates(&modelState.Variables) {
			nonDominatedArray = append(nonDominatedArray, a.archive[currentIndex])
		}
	}

	a.archive = nonDominatedArray

	a.archive = append(a.archive, modelState)
	return StoredForcingDominatingStateRemoval
}

func (a *NonDominanceModelArchive) newModelStateCannotBeArchived(modelState *CompressedModelState) StorageResult {
	for currentIndex := 0; currentIndex < len(a.archive); currentIndex++ {
		if a.archive[currentIndex].Variables.Dominates(&modelState.Variables) {
			return RejectedWithStoredEntryDominanceDetected
		}
		if a.archive[currentIndex].Actions.IsEquivalentTo(&modelState.Actions) {
			return RejectedWithDuplicateEntryDetected
		}
	}
	return canBeStored
}

func (a *NonDominanceModelArchive) IsEmpty() bool {
	return len(a.archive) == 0
}

func (a *NonDominanceModelArchive) Len() int {
	return len(a.archive)
}

func (a *NonDominanceModelArchive) IsNonDominant() bool {
	if a.IsEmpty() {
		return true
	}

	for currentIndex := 0; currentIndex < len(a.archive); currentIndex++ {
		for downstreamIndex := currentIndex + 1; downstreamIndex < len(a.archive)-1; downstreamIndex++ {
			if a.archive[currentIndex].Variables.DominancePresent(&a.archive[downstreamIndex].Variables) {
				return false // Shouldn't happen if AttemptToArchiveState() is successfully ensuring non-dominance holds regardless.
			}
		}
	}
	return true
}

func (a *NonDominanceModelArchive) ArchiveSummary() Summary {
	summary := a.buildSummary()

	if summary == nil {
		return nil
	}

	for variableIndex, variableValue := range a.archive[0].Variables {
		summary[variableIndex].Minimum = variableValue
		summary[variableIndex].Maximum = variableValue
		summary[variableIndex].Range = 0
	}

	for _, entry := range a.archive {
		for variableIndex, variableValue := range entry.Variables {
			rangeChanged := false
			if variableValue < summary[variableIndex].Minimum {
				summary[variableIndex].Minimum = variableValue
				rangeChanged = true
			}
			if variableValue > summary[variableIndex].Maximum {
				summary[variableIndex].Maximum = variableValue
				rangeChanged = true
			}
			if rangeChanged {
				summary[variableIndex].Range = math.Abs(summary[variableIndex].Maximum - summary[variableIndex].Minimum)
			}
		}
	}

	return summary
}

func (a *NonDominanceModelArchive) SelectRandomModel() *CompressedModelState {
	// This is the Engrand approach, and constant time in computational complexity, with unlikely isolation model selection..
	fullArchiveRange := len(a.archive)
	return a.selectRandomModel(fullArchiveRange)
}

func (a *NonDominanceModelArchive) SelectRandomIsolatedModel(selectionRange int) *CompressedModelState {
	// This is Suppapitmarm approach, and archiveSize^2 in computational complexity for good isolation model selection.
	a.calculateArchiveIsolation()
	sort.Sort(a)
	return a.selectRandomModel(selectionRange)
}

func (a *NonDominanceModelArchive) calculateArchiveIsolation() {
	summary := a.ArchiveSummary()
	a.isolation = make([]float64, len(a.archive))
	for index, modelState := range a.archive {
		a.isolation[index] = a.calculatedModelStateIsolation(modelState, summary)
	}
}

func (a *NonDominanceModelArchive) calculatedModelStateIsolation(state *CompressedModelState, summary Summary) float64 {
	if a.isIsolated(summary, state) {
		return a.deriveIsolation(summary, state)
	}
	return notIsolated
}

func (a *NonDominanceModelArchive) isIsolated(summary Summary, state *CompressedModelState) bool {
	for index, variableSummary := range summary {
		if variableSummary.Minimum == state.Variables[index] {
			return false
		}
	}
	return true
}

func (a *NonDominanceModelArchive) deriveIsolation(summary Summary, state *CompressedModelState) float64 {
	var isolation float64
	for _, modelState := range a.archive {
		if state.IsEquivalentTo(modelState) {
			continue
		}
		for summaryIndex, _ := range summary {
			numerator := (state.Variables[summaryIndex] - modelState.Variables[summaryIndex]) / summary[summaryIndex].Range
			isolation += math.Pow(numerator, 2)
		}
	}
	return isolation
}

func (a *NonDominanceModelArchive) selectRandomModel(selectionRange int) *CompressedModelState {
	selectedIndex := a.selectRandomIndex(selectionRange)
	return a.archive[selectedIndex]
}

func (a *NonDominanceModelArchive) selectRandomIndex(selectionRange int) int {
	assert.That(selectionRange > 0 && selectionRange <= len(a.archive)).
		WithFailureMessage("Selection range invalid for size of archive").
		Holds()

	randomIndex := a.RandomNumberGenerator().Intn(selectionRange)
	return randomIndex
}

func (a *NonDominanceModelArchive) buildSummary() Summary {
	summary := make(Summary, 0)

	if a.IsEmpty() {
		return nil
	}

	for index := range a.archive[0].Variables {
		summary[index] = &VariableSummary{}
	}

	return summary
}

func (a NonDominanceModelArchive) Swap(i, j int) {
	a.archive[i], a.archive[j] = a.archive[j], a.archive[i]
	a.isolation[i], a.isolation[j] = a.isolation[j], a.isolation[i]
}

func (a NonDominanceModelArchive) Less(i, j int) bool {
	return a.isolation[i] < a.isolation[j]
}

type Summary map[int]*VariableSummary

type VariableSummary struct {
	Minimum float64
	Maximum float64
	Range   float64
}
