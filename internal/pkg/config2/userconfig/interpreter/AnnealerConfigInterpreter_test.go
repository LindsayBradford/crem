// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/onsi/gomega"
)

func TestConfigInterpreter_NewAnnealerConfigInterpreter_NullAnnealerNoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	interpreterUnderTest := NewAnnealerConfigInterpreter()

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())

	g.Expect(interpreterUnderTest.Annealer()).To(BeAssignableToTypeOf(&annealers.NullAnnealer{}))
}

func TestConfigInterpreter_NullAnnealer_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.AnnealerConfig{
		Type: data.UnspecifiedAnnealerType,
	}

	// when
	interpreterUnderTest := NewAnnealerConfigInterpreter().Interpret(&configUnderTest)

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())

	actualAnnealer := interpreterUnderTest.Annealer()
	expectedAnnealerType := &annealers.NullAnnealer{}
	g.Expect(actualAnnealer).To(BeAssignableToTypeOf(expectedAnnealerType))
}

func TestConfigInterpreter_KirkpatrickAnnealer_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	parametersUnderTest := parameters.Map{
		"MaximumIterations":     int64(100),
		"DecisionVariable":      "ObjectiveVariable",
		"OptimisationDirection": "Maximising",
		"CoolingFactor":         0.95,
		"StartingTemperature":   float64(10),
	}
	configUnderTest := data.AnnealerConfig{
		Type:       data.Kirkpatrick,
		Parameters: parametersUnderTest,
	}

	// when
	interpreterUnderTest := NewAnnealerConfigInterpreter().Interpret(&configUnderTest)

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())

	actualAnnealer := interpreterUnderTest.Annealer()
	expectedAnnealerType := &annealers.ElapsedTimeTrackingAnnealer{}
	g.Expect(actualAnnealer).To(BeAssignableToTypeOf(expectedAnnealerType))

	actualExplorer := interpreterUnderTest.Annealer().SolutionExplorer()
	expectedExplorerType := &kirkpatrick.Explorer{}
	g.Expect(actualExplorer).To(BeAssignableToTypeOf(expectedExplorerType))
}

func TestConfigInterpreter_SuppapitnarmAnnealer_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.AnnealerConfig{
		Type: data.Suppapitnarm,
	}

	// when
	interpreterUnderTest := NewAnnealerConfigInterpreter().Interpret(&configUnderTest)

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())

	actualAnnealer := interpreterUnderTest.Annealer()
	expectedAnnealerType := &annealers.ElapsedTimeTrackingAnnealer{}
	g.Expect(actualAnnealer).To(BeAssignableToTypeOf(expectedAnnealerType))

	actualExplorer := interpreterUnderTest.Annealer().SolutionExplorer()
	expectedExplorerType := &kirkpatrick.Explorer{} // TODO: Should be dedidicated Suppapitnarm explorer.
	g.Expect(actualExplorer).To(BeAssignableToTypeOf(expectedExplorerType))

}
