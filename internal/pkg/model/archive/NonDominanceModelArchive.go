// Copyright (c) 2019 Australian Rivers Institute.

package archive

import "github.com/LindsayBradford/crem/internal/pkg/model"

type StorageResult uint

const (
	StoredReplacingDominatedEntries StorageResult = iota
	StoredWithNoDominanceDetected
	RejectedWithStoredEntryDominanceDetected
	RejectedWithDuplicateEntryDetected
	canBeStored
	StoredForcingDominatingStateRemoval
)

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
