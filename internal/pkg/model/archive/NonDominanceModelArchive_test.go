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
	numberOfRandomChanges := 200
	for change := 0; change < numberOfRandomChanges; change++ {
		modelToChange.DoRandomChange()

		archiveUnderTest.Archive(modelToChange)

		for _, entry := range archiveUnderTest.archive {
			t.Logf("%v\n", entry.Variables)
		}
		t.Log("----\n")

		// then
		g.Expect(archiveUnderTest.IsNonDominant()).To(BeTrue())
	}

	g.Expect(true).To(BeTrue())
}

func buildSilentMultiObjectiveDumbModel() *modumb.Model {
	model := modumb.NewModel().WithId("Test Mo Dumb Model")
	model.SetEventNotifier(loggers.NullTestingEventNotifier)
	model.Initialise()
	return model
}
