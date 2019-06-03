// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/threading"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestArchivist_Store(t *testing.T) {
	// TODO: Rid myself of my reliance on CatchmentModel for testing this.  Too much prep-work for the model.
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)
	testModel := catchment.NewModel()

	expectedVariableSize := len(*testModel.DecisionVariables())
	expectedActionSize := len(testModel.ManagementActions())

	// when
	actualArchive := archivistUnderTest.Store(testModel)

	// then
	g.Expect(expectedVariableSize).To(BeNumerically(equalTo, len(*actualArchive.Variables())))
	g.Expect(expectedActionSize).To(BeNumerically(equalTo, actualArchive.Actions().Len()))
}

func TestArchivist_Retrieve_InitialModel(t *testing.T) {
	// TODO: Rid myself of my reliance on CatchmentModel for testing this.  Too much prep-work for the model.
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)

	modelToStore := catchment.NewModel()
	modelToRetrieve := modelToStore.DeepClone()

	// when
	storedArchive := archivistUnderTest.Store(modelToStore)
	archivistUnderTest.Retrieve(storedArchive, modelToRetrieve)

	// then
	g.Expect(modelToStore).To(Equal(modelToRetrieve))
}

func testArchivist_Retrieve_AlteredModel(t *testing.T) {
	// TODO: Rid myself of my reliance on CatchmentModel for testing this.  Too much prep-work for the model.
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)

	modelToStore := catchment.NewModel().
		WithName("test").
		WithOleFunctionWrapper(threading.GetMainThreadChannel().Call)

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
