// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

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

		archiveUnderTest.Archive(modelToChange)
		actualArchiveSize = len(archiveUnderTest.archive)

		// then
		g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	}
	t.Logf("Model cnanges required for [%d] archive entrries = %d\n", actualArchiveSize, changesRequired)

	for _, entry := range archiveUnderTest.archive {
		t.Logf("%v\n", entry.Variables)
	}

	g.Expect(true).To(BeTrue())
}

func buildSilentMultiObjectiveDumbModel() *modumb.Model {
	model := modumb.NewModel().WithId("Test Mo Dumb Model")
	model.SetEventNotifier(loggers.NullTestingEventNotifier)
	model.Initialise()
	return model
}
