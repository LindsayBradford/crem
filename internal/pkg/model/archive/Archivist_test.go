// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestArchivist_Store(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)
	testModel := modumb.NewModel()

	expectedVariableSize := len(*testModel.DecisionVariables())
	expectedActionSize := len(testModel.ManagementActions())

	// when
	actualArchive := archivistUnderTest.Store(testModel)

	// then
	g.Expect(expectedVariableSize).To(BeNumerically(equalTo, len(*actualArchive.Variables())))
	g.Expect(expectedActionSize).To(BeNumerically(equalTo, actualArchive.Actions().Len()))
}

func TestArchivist_Retrieve_InitialModel(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)

	modelToStore := modumb.NewModel()
	modelToStore.Initialise()
	modelToRetrieve := modelToStore.DeepClone()

	g.Expect(modelToStore.DecisionVariables()).To(Equal(modelToRetrieve.DecisionVariables()))

	// when
	storedArchive := archivistUnderTest.Store(modelToStore)
	archivistUnderTest.Retrieve(storedArchive, modelToRetrieve)

	// then
	g.Expect(modelToStore.DecisionVariables()).To(Equal(modelToRetrieve.DecisionVariables()))
}

func testArchivist_Retrieve_AlteredModel(t *testing.T) {
	// TODO: Test is failing on model equality... why?
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)

	modelToStore := modumb.NewModel()

	modelToStore.SetEventNotifier(new(observer.SynchronousAnnealingEventNotifier))
	modelToStore.Initialise()

	numberOfRandomChanges := 7
	for change := 0; change < numberOfRandomChanges; change++ {
		modelToStore.TryRandomChange()
	}

	modelToRetrieve := modelToStore.DeepClone()

	// when
	storedArchive := archivistUnderTest.Store(modelToStore)
	archivistUnderTest.Retrieve(storedArchive, modelToRetrieve)

	// then
	g.Expect(modelToStore).To(Equal(modelToRetrieve))
}
