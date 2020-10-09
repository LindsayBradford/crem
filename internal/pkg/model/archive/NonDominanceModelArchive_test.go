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
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 2))

	// when
	modelToChange.SetManagementAction(1, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))

	// when
	modelToChange.SetManagementAction(2, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredReplacingDominatedEntries))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))
}

func TestNonDominanceModelArchive_ArchiveAttemptOfDominatedRejectedForcedAccepted(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()
	storageResult := archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 1))

	// when
	modelToChange.SetManagementAction(0, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 2))

	modelToChange.SetManagementAction(1, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredWithNoDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))

	modelToChange.SetManagementAction(2, true)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredReplacingDominatedEntries))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))

	modelToChange.SetManagementAction(2, false)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(RejectedWithDuplicateEntryDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))

	modelToChange.SetManagementAction(1, false)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(RejectedWithDuplicateEntryDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))

	modelToChange.SetManagementAction(0, false)
	modelToChange.AcceptChange()
	storageResult = archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(RejectedWithStoredEntryDominanceDetected))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))

	// when

	t.Logf("Forcing last change into archive")
	storageResult = archiveUnderTest.ForceIntoArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// then
	g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	g.Expect(storageResult).To(Equal(StoredForcingDominatingStateRemoval))
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, 3))
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

func TestNonDominanceModelArchiveSummary_EmptyArchive_NoSummary(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	archiveUnderTest := New()
	// then
	g.Expect(archiveUnderTest.ArchiveSummary()).To(BeNil())
}

func TestNonDominanceModelArchiveSummary_OneModelArchive_NoRange(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()
	archiveUnderTest.AttemptToArchive(modelToChange)

	// when
	summary := archiveUnderTest.ArchiveSummary()

	// then
	for _, entry := range summary {
		g.Expect(entry.Minimum).To(BeNumerically(equalTo, entry.Maximum))
		g.Expect(entry.Range).To(BeNumerically(equalTo, 0))
	}
}

func TestNonDominanceModelArchiveSummary_Archive_SummaryValid(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	modelToChange := buildSilentMultiObjectiveDumbModel()
	archiveUnderTest := New()
	archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	// when
	summary := archiveUnderTest.ArchiveSummary()

	// then
	for _, entry := range summary {
		g.Expect(entry.Minimum).To(BeNumerically(equalTo, entry.Maximum))
		g.Expect(entry.Range).To(BeNumerically(equalTo, 0))
	}

	// when
	modelToChange.SetManagementAction(0, true)
	modelToChange.AcceptChange()
	archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	summary = archiveUnderTest.ArchiveSummary()

	// then
	g.Expect(len(archiveUnderTest.archive)).To(BeNumerically(equalTo, 2))

	for index, entry := range summary {
		switch index {
		case 0:
			g.Expect(entry.Minimum).To(BeNumerically(equalTo, 999))
			g.Expect(entry.Maximum).To(BeNumerically(equalTo, 1000))
			g.Expect(entry.Range).To(BeNumerically(equalTo, 1))
		default:
			g.Expect(entry.Minimum).To(BeNumerically(equalTo, entry.Maximum))
			g.Expect(entry.Range).To(BeNumerically(equalTo, 0))
		}
	}

	// when
	modelToChange.SetManagementAction(0, false)
	modelToChange.AcceptChange()
	archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	modelToChange.SetManagementAction(1, true)
	modelToChange.AcceptChange()
	archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	summary = archiveUnderTest.ArchiveSummary()

	// then
	g.Expect(len(archiveUnderTest.archive)).To(BeNumerically(equalTo, 3))

	g.Expect(summary[0].Minimum).To(BeNumerically(equalTo, 999))
	g.Expect(summary[0].Maximum).To(BeNumerically(equalTo, 1000))
	g.Expect(summary[0].Range).To(BeNumerically(equalTo, 1))

	g.Expect(summary[1].Minimum).To(BeNumerically(equalTo, 1998))
	g.Expect(summary[1].Maximum).To(BeNumerically(equalTo, 2000))
	g.Expect(summary[1].Range).To(BeNumerically(equalTo, 2))

	// when
	modelToChange.SetManagementAction(1, false)
	modelToChange.AcceptChange()
	archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	modelToChange.SetManagementAction(2, true)
	modelToChange.AcceptChange()
	archiveUnderTest.AttemptToArchive(modelToChange)
	showArchiveState(t, archiveUnderTest)

	summary = archiveUnderTest.ArchiveSummary()

	// then
	g.Expect(len(archiveUnderTest.archive)).To(BeNumerically(equalTo, 4))

	g.Expect(summary[0].Minimum).To(BeNumerically(equalTo, 999))
	g.Expect(summary[0].Maximum).To(BeNumerically(equalTo, 1000))
	g.Expect(summary[0].Range).To(BeNumerically(equalTo, 1))

	g.Expect(summary[1].Minimum).To(BeNumerically(equalTo, 1998))
	g.Expect(summary[1].Maximum).To(BeNumerically(equalTo, 2000))
	g.Expect(summary[1].Range).To(BeNumerically(equalTo, 2))

	g.Expect(summary[2].Minimum).To(BeNumerically(equalTo, 2997))
	g.Expect(summary[2].Maximum).To(BeNumerically(equalTo, 3000))
	g.Expect(summary[2].Range).To(BeNumerically(equalTo, 3))
}

func buildSilentMultiObjectiveDumbModel() *modumb.Model {
	model := modumb.NewModel().WithId("Test Mo Dumb Model")
	model.AddObserver(loggers.DefaultTestingAnnealingObserver)
	model.Initialise()
	return model
}

func showArchiveState(t *testing.T, archive *NonDominanceModelArchive) {
	for _, entry := range archive.archive {
		t.Logf("Variables: %v, action SHA265: %s\n", entry.Variables, entry.Actions.Sha256())
	}
}
