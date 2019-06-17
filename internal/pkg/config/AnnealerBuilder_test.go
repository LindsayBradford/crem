// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	observer2 "github.com/LindsayBradford/crem/internal/pkg/observer"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestAnnealerBuilder_MinimalDumbValidConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerMinimalValidConfig.toml")
	g.Expect(retrieveError).To(BeNil())

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(BeNil())
	g.Expect(logHandler).To(Not(BeNil()))

	dummyAnnealer := new(annealers.ElapsedTimeTrackingAnnealer)

	g.Expect(annealerUnderTest).To(BeAssignableToTypeOf(dummyAnnealer))

	annealerAttributes := annealerUnderTest.EventAttributes(observer2.FinishedAnnealing)
	actualMaximumIterations := annealerAttributes.Value(annealers.MaximumIterations).(uint64)

	g.Expect(actualMaximumIterations).To(BeNumerically(equalTo, 5))

	solutionExplorerUnderTest := annealerUnderTest.SolutionExplorer()

	g.Expect(solutionExplorerUnderTest.Name()).To(Equal("validConfig"))

	dummyExplorer := kirkpatrick.New()

	g.Expect(solutionExplorerUnderTest).To(BeAssignableToTypeOf(dummyExplorer))
}

func TestAnnealerBuilder_MinimalNullValidConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/NullAnnealerMinimalValidConfig.toml")
	g.Expect(retrieveError).To(BeNil())

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(BeNil())
	g.Expect(logHandler).To(Not(BeNil()))

	dummyAnnealer := new(annealers.SimpleAnnealer)

	g.Expect(
		annealerUnderTest).To(BeAssignableToTypeOf(dummyAnnealer))

	annealerAttributes := annealerUnderTest.EventAttributes(observer2.FinishedAnnealing)
	actualMaximumIterations := annealerAttributes.Value(annealers.MaximumIterations).(uint64)

	g.Expect(actualMaximumIterations).To(BeNumerically(equalTo, 5))

	solutionExplorerUnderTest := annealerUnderTest.SolutionExplorer()

	g.Expect(
		solutionExplorerUnderTest.Name()).To(Equal("validConfig"))

	dummyExplorer := null.NullExplorer

	g.Expect(
		solutionExplorerUnderTest).To(BeAssignableToTypeOf(dummyExplorer))
}

func TestAnnealerBuilder_DumbAnnealerInvalidAnnealerTypeConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidAnnealerTypeConfig.toml")
	g.Expect(retrieveError).To(Not(BeNil()), "Config retrieval should not have failed.")
	g.Expect(configUnderTest).To(BeNil(), "Config retrieves should have been nil.")

	g.Expect(retrieveError.Error()).To(ContainSubstring("invalid value \"SomeTotallyUnknownAnnealerType\""))
	t.Logf("config retrieval error reported: %s", retrieveError)
}

func TestAnnealerBuilder_DumbAnnealerInvalidEventNotifierConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidEventNotifierConfig.toml")
	g.Expect(retrieveError).To(Not(BeNil()), "Config retrieval should not have failed.")
	g.Expect(configUnderTest).To(BeNil(), "Config retrieves should have been nil.")

	g.Expect(retrieveError.Error()).To(ContainSubstring("invalid value \"NotTheEventNotifierYouAreLookingFor\""))
	t.Logf("config retrieval error reported: %s", retrieveError)
}

func TestAnnealerBuilder_MissingSolutionExplorerTypeConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerMissingSolutionExplorerTypeConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(Not(BeNil()), "Annealer build should have failed.")
	g.Expect(buildError.Error()).To(ContainSubstring(" no explorers are registered for that type"))
	t.Logf("Annealer build error reported: %s", buildError)

	g.Expect(annealerUnderTest).To(BeNil(), "Annealer build failure should have returned nil annealer.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build failure should have returned valid logHandler.")
}

func TestAnnealerBuilder_MismatchedSolutionExplorerNamesConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerMismatchedSolutionExplorerNamesConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(Not(BeNil()), "Annealer build should have failed.")
	g.Expect(buildError.Error()).To(ContainSubstring("configuration specifies a non-existent explorer"))
	t.Logf("Annealer build error reported: %s", buildError)

	g.Expect(annealerUnderTest).To(BeNil(), "Annealer build failure should have returned nil annealer.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build failure should have returned valid logHandler.")
}

func TestAnnealerBuilder_DumbAnnealerNoSolutionExplorerConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerNoSolutionExplorerConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(Not(BeNil()), "Annealer build should have failed.")
	g.Expect(buildError.Error()).To(ContainSubstring("configuration failed to specify any explorers"))
	t.Logf("Annealer build error reported: %s", buildError)

	g.Expect(annealerUnderTest).To(BeNil(), "Annealer build failure should have returned nil annealer.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build failure should have returned valid logHandler.")
}

func TestAnnealerBuilder_DumbAnnealerInvalidAnnealingObserversTypeConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidAnnealingObserversTypeConfig.toml")
	g.Expect(retrieveError).To(Not(BeNil()), "Config retrieval should not have failed.")
	g.Expect(configUnderTest).To(BeNil(), "Config retrieves should have been nil.")

	g.Expect(retrieveError.Error()).To(ContainSubstring("invalid value \"NoKnownValidAnnealingObserverType\""))
	t.Logf("config retrieval error reported: %s", retrieveError)
}

func TestAnnealerBuilder_DumbAnnealerInvalidIterationFilterConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidIterationFilterConfig.toml")
	g.Expect(retrieveError).To(Not(BeNil()), "Config retrieval should not have failed.")
	g.Expect(configUnderTest).To(BeNil(), "Config retrieves should have been nil.")

	g.Expect(retrieveError.Error()).To(ContainSubstring("invalid value \"ThereAintNoSuchFilter\""))
	t.Logf("config retrieval error reported: %s", retrieveError)
}

func TestAnnealerBuilder_InvalidAnnealingObserversLoggerConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidAnnealingObserversLoggerConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(Not(BeNil()), "Annealer build should have failed.")
	g.Expect(buildError.Error()).To(ContainSubstring("configuration specifies a non-existent logger"))
	t.Logf("Annealer build error reported: %s", buildError)

	g.Expect(annealerUnderTest).To(BeNil(), "Annealer build failure should have returned nil annealer.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build failure should have returned valid logHandler.")
}

func TestAnnealerBuilder_DumbAnnealerInvalidLoggersTypeConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidLoggersTypeConfig.toml")
	g.Expect(retrieveError).To(Not(BeNil()), "Config retrieval should not have failed.")
	g.Expect(configUnderTest).To(BeNil(), "Config retrieves should have been nil.")

	g.Expect(retrieveError.Error()).To(ContainSubstring("invalid value \"BorkedLoggerType\""))
	t.Logf("config retrieval error reported: %s", retrieveError)
}

func TestAnnealerBuilder_DumbAnnealerInvalidLoggersFormatterConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidLoggersFormatterConfig.toml")
	g.Expect(retrieveError).To(Not(BeNil()), "Config retrieval should not have failed.")
	g.Expect(configUnderTest).To(BeNil(), "Config retrieves should have been nil.")

	g.Expect(retrieveError.Error()).To(ContainSubstring("invalid value \"That's no ordinary formatter! It's got great big gnashy teeth!\""))
	t.Logf("config retrieval error reported: %s", retrieveError)
}

func TestAnnealerBuilder_DumbAnnealerInvalidLogLevelDestinationConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerInvalidLogLevelDestinationConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(Not(BeNil()), "Annealer build should have failed.")
	g.Expect(buildError.Error()).To(ContainSubstring("unrecognised destination [The Hedgehog Song]"))
	t.Logf("Annealer build error reported: %s", buildError)

	g.Expect(annealerUnderTest).To(BeNil(), "Annealer build failure should have returned nil annealer.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build failure should have returned valid logHandler.")
}

func TestAnnealerBuilder_DumbAnnealerRichValidConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerRichValidConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(BeNil(), "Annealer build should not have failed.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build should have returned a valid logHandler.")

	dummyAnnealer := new(annealers.ElapsedTimeTrackingAnnealer)

	g.Expect(
		annealerUnderTest).To(BeAssignableToTypeOf(dummyAnnealer),
		"Annealer should have built with config supplied annealer type")

	annealerAttributes := annealerUnderTest.EventAttributes(observer2.FinishedAnnealing)
	actualMaximumIterations := annealerAttributes.Value(annealers.MaximumIterations).(uint64)

	expectedMaximumIterations := 2000
	g.Expect(actualMaximumIterations).To(BeNumerically(equalTo, expectedMaximumIterations))

	solutionExplorerUnderTest := annealerUnderTest.SolutionExplorer()

	g.Expect(solutionExplorerUnderTest.Name()).To(Equal("DoraTheExplorer"))

	explorer := new(kirkpatrick.Explorer)

	g.Expect(solutionExplorerUnderTest).To(BeAssignableToTypeOf(explorer))

	actualObservers := annealerUnderTest.Observers()

	g.Expect(len(actualObservers)).To(BeNumerically(equalTo, 3))

	dummyMessageObserver := new(observer.AnnealingMessageObserver)

	g.Expect(actualObservers[0]).To(BeAssignableToTypeOf(dummyMessageObserver))

	dummyAttributeObserver := new(observer.AnnealingAttributeObserver)

	g.Expect(actualObservers[1]).To(BeAssignableToTypeOf(dummyAttributeObserver))
}

func NewTestRegisterableExplorer() *TestRegistereableExplorer {
	explorer := new(TestRegistereableExplorer)
	explorer.Explorer = *kirkpatrick.New()
	return explorer
}

type TestRegistereableExplorer struct {
	kirkpatrick.Explorer
}

func (tre *TestRegistereableExplorer) WithName(name string) *TestRegistereableExplorer {
	tre.SetName(name)
	return tre
}

func TestAnnealerBuilder_NullAnnealerWithCustomExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/NullAnnealerWithCustomExplorerConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder).RegisteringExplorer(
		ExplorerRegistration{
			ExplorerType: "TestDefinedSolutionExplorer",
			ConfigFunction: func(config SolutionExplorerConfig) explorer.Explorer {
				return NewTestRegisterableExplorer().WithName(config.Name)
			},
		},
	)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(BeNil(), "Annealer build should not have failed.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build should have returned a valid logHandler.")

	solutionExplorerUnderTest := annealerUnderTest.SolutionExplorer()

	g.Expect(
		solutionExplorerUnderTest.Name()).To(Equal("testyName"),
		"Annealer should have built with config supplied Explorer")

	dummyExplorer := new(TestRegistereableExplorer)

	g.Expect(
		solutionExplorerUnderTest).To(BeAssignableToTypeOf(dummyExplorer),
		"Annealer should have built with config supplied Explorer")
}
