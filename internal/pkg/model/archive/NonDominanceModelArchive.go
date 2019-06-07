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

func (a *NonDominanceModelArchive) Archive(model model.Model) StorageResult {
	compressedModelState := a.compressor.Compress(model)
	return a.ArchiveState(compressedModelState)
}

func (a *NonDominanceModelArchive) ArchiveState(modelState *CompressedModelState) StorageResult {
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

func (a *NonDominanceModelArchive) newModelStateIsRejectedByArchive(modelState *CompressedModelState) bool {
	for currentIndex := 0; currentIndex < len(a.archive); currentIndex++ {
		if a.archive[currentIndex].Variables.Dominates(&modelState.Variables) {
			return true
		}
	}
	return false
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

func (a *NonDominanceModelArchive) IsNonDominant() bool {
	if a.IsEmpty() {
		return true
	}

	for currentIndex := 0; currentIndex < len(a.archive); currentIndex++ {
		for downstreamIndex := currentIndex + 1; downstreamIndex < len(a.archive)-1; downstreamIndex++ {
			if a.archive[currentIndex].Variables.DominancePresent(&a.archive[downstreamIndex].Variables) {
				return false
			}
		}
	}
	return true
}

func (a *NonDominanceModelArchive) ArchiveArray() []*CompressedModelState {
	return a.archive
}
