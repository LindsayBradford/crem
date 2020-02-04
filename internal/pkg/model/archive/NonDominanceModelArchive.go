// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/model"
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
		return "Stored with no dominance detected"
	case StoredReplacingDominatedEntries:
		return "Stored replacing dominated archive entries"
	case RejectedWithStoredEntryDominanceDetected:
		return "Rejected with archive reporting dominance of stored entries"
	case RejectedWithDuplicateEntryDetected:
		return "Rejected with archive reporting solution already archived"
	case StoredForcingDominatingStateRemoval:
		return "Stored, forcing dominance entries out of archive"
	}
	return "Can be stored -- but why are you seeing this?!?"
}

func New() *NonDominanceModelArchive {
	return new(NonDominanceModelArchive).Initialise()
}

type NonDominanceModelArchive struct {
	archive    []*CompressedModelState
	compressor ModelCompressor
}

func (a *NonDominanceModelArchive) Initialise() *NonDominanceModelArchive {
	a.archive = make([]*CompressedModelState, 0)
	return a
}

func (a *NonDominanceModelArchive) AttemptToArchive(model model.Model) StorageResult {
	compressedModelState := a.compressor.Compress(model)
	return a.AttemptToArchiveState(compressedModelState)
}

func (a *NonDominanceModelArchive) Compress(model model.Model) *CompressedModelState {
	return a.compressor.Compress(model)
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

func (a *NonDominanceModelArchive) ArchiveSummary() ArchiveSummary {
	summary := a.buildSummary()

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

func (a *NonDominanceModelArchive) buildSummary() ArchiveSummary {
	summary := make(ArchiveSummary, 0)

	if a.IsEmpty() {
		return nil
	}

	for index := range a.archive[0].Variables {
		summary[index] = &VariableSummary{}
	}

	return summary
}

type ArchiveSummary map[int]*VariableSummary

type VariableSummary struct {
	Minimum float64
	Maximum float64
	Range   float64
}
