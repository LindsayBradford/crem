// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestNonDominanceModelArchive_EmptyArchive_IsNonDominant(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	archiveUnderTest := New()
	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
}

func TestNonDominanceModelArchive_ArchiveSetPropertyMaintained(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()

	// when
	modelToChange.DoRandomChange()
	storageResult := archiveUnderTest.AttemptToArchive(modelToChange)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))

	// when
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(RejectedWithDuplicateEntryDetected))
}

func TestNonDominanceModelArchive_DominatorReplaceDominated(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()

	// when
	storageResult := archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))

	// when
	modelToChange.SetManagementAction(0, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredReplacingDominatedEntries))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))

	// when
	modelToChange.SetManagementAction(3, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredReplacingDominatedEntries))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))
}

func TestNonDominanceModelArchive_ArchiveAttemptOfDominatedRejected(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()

	// when
	storageResult := archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))

	// when
	modelToChange.SetManagementAction(0, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredReplacingDominatedEntries))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))

	// when
	modelToChange.SetManagementAction(0, false)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(RejectedWithStoredEntryDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))
}

func TestNonDominanceModelArchive_ForceArchiveAttemptStored(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()

	// when
	storageResult := archiveUnderTest.AttemptToArchive(modelToChange)

	modelToChange.SetManagementAction(0, true)
	modelToChange.AcceptChange()
	archiveUnderTest.AttemptToArchive(modelToChange)

	modelToChange.SetManagementAction(3, true)
	modelToChange.AcceptChange()
	archiveUnderTest.AttemptToArchive(modelToChange)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))
	showArchiveState(t, archiveUnderTest)

	// when

	modelToChange.SetManagementAction(3, false)
	modelToChange.AcceptChange()

	storageResult = archiveUnderTest.ForceIntoArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredForcingDominatingStateRemoval))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))
}

func TestNonDominanceModelArchive_ChangesPreserveNonDominance(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()

	// when
	desiredArchiveSize := 7
	actualArchiveSize := len(archiveUnderTest.archive)
	var changesRequired uint
	for actualArchiveSize < desiredArchiveSize {
		modelToChange.DoRandomChange()
		changesRequired++

		archiveUnderTest.AttemptToArchive(modelToChange)
		actualArchiveSize = len(archiveUnderTest.archive)

		// then
		g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	}

	t.Logf("Model cnanges required for [%d] archive entrries = %d\n", actualArchiveSize, changesRequired)

	showArchiveState(t, archiveUnderTest)
}

func buildSilentMultiObjectiveDumbModel() *modumb.Model {
	model := modumb.NewModel().WithId("Test Mo Dumb Model")
	model.SetEventNotifier(loggers.NullTestingEventNotifier)
	model.Initialise()
	return model
}

func showArchiveState(t *testing.T, archive *NonDominanceModelArchive) {
	for _, entry := range archive.archive {
		t.Logf("Variables: %v, action SHA265: %s\n", entry.Variables, entry.Actions.Sha256())
	}
}
