// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"testing"

	"github.com/LindsayBradford/crem/annealing/annealers"
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/annealing/observer"
	. "github.com/onsi/gomega"
)

func TestAnnealerBuilder_MinimalDumbValidConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/DumbAnnealerMinimalValidConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(BeNil(), "Annealer build should not have failed.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build should have returned a valid logHandler.")

	dummyAnnealer := new(annealers.ElapsedTimeTrackingAnnealer)

	g.Expect(
		annealerUnderTest).To(BeAssignableToTypeOf(dummyAnnealer),
		"Annealer should have built with default annealer type")

	g.Expect(
		annealerUnderTest.Temperature()).To(BeNumerically("==", 10),
		"Annealer should have built with config supplied Temperature")

	g.Expect(
		annealerUnderTest.CoolingFactor()).To(BeNumerically("==", 0.99),
		"Annealer should have built with config supplied CoolingFactor")

	g.Expect(
		annealerUnderTest.MaxIterations()).To(BeNumerically("==", 5),
		"Annealer should have built with config supplied MaxIterations")

	solutionExplorerUnderTest := annealerUnderTest.SolutionExplorer()

	g.Expect(
		solutionExplorerUnderTest.Name()).To(Equal("validConfig"),
		"Annealer should have built with config supplied Explorer")

	dummyExplorer := new(explorer.DumbExplorer)

	g.Expect(
		solutionExplorerUnderTest).To(BeAssignableToTypeOf(dummyExplorer),
		"Annealer should have built with config supplied Explorer")
}

func TestAnnealerBuilder_MinimalNullValidConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	configUnderTest, retrieveError := RetrieveCremFromFile("testdata/NullAnnealerMinimalValidConfig.toml")
	g.Expect(retrieveError).To(BeNil(), "Config retrieval should not have failed.")

	builderUnderTest := new(AnnealerBuilder)

	annealerUnderTest, logHandler, buildError :=
		builderUnderTest.WithConfig(configUnderTest).Build()

	g.Expect(buildError).To(BeNil(), "Annealer build should not have failed.")
	g.Expect(logHandler).To(Not(BeNil()), "Annealer build should have returned a valid logHandler.")

	dummyAnnealer := new(annealers.SimpleAnnealer)

	g.Expect(
		annealerUnderTest).To(BeAssignableToTypeOf(dummyAnnealer),
		"Annealer should have built with default annealer type")

	g.Expect(
		annealerUnderTest.Temperature()).To(BeNumerically("==", 10),
		"Annealer should have built with config supplied Temperature")

	g.Expect(
		annealerUnderTest.CoolingFactor()).To(BeNumerically("==", 0.99),
		"Annealer should have built with config supplied CoolingFactor")

	g.Expect(
		annealerUnderTest.MaxIterations()).To(BeNumerically("==", 5),
		"Annealer should have built with config supplied MaxIterations")

	solutionExplorerUnderTest := annealerUnderTest.SolutionExplorer()

	g.Expect(
		solutionExplorerUnderTest.Name()).To(Equal("validConfig"),
		"Annealer should have built with config supplied Explorer")

	dummyExplorer := new(explorer.NullExplorer)

	g.Expect(
		solutionExplorerUnderTest).To(BeAssignableToTypeOf(dummyExplorer),
		"Annealer should have built with config supplied Explorer")
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
	g.Expect(buildError.Error()).To(ContainSubstring("configuration specifies a non-existent explorer explorer"))
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

	g.Expect(
		annealerUnderTest.Temperature()).To(BeNumerically("==", 50),
		"Annealer should have built with config supplied Temperature")

	g.Expect(
		annealerUnderTest.CoolingFactor()).To(BeNumerically("==", 0.995),
		"Annealer should have built with config supplied CoolingFactor")

	g.Expect(
		annealerUnderTest.MaxIterations()).To(BeNumerically("==", 2000),
		"Annealer should have built with config supplied MaxIterations")

	solutionExplorerUnderTest := annealerUnderTest.SolutionExplorer()

	g.Expect(
		solutionExplorerUnderTest.Name()).To(Equal("DoraTheExplorer"),
		"Annealer should have built with config supplied Explorer")

	dummyExplorer := new(explorer.DumbExplorer)

	g.Expect(
		solutionExplorerUnderTest).To(BeAssignableToTypeOf(dummyExplorer),
		"Annealer should have built with config supplied Explorer")

	actualObservers := annealerUnderTest.Observers()

	g.Expect(
		len(actualObservers)).To(BeNumerically("==", 3),
		"Annealer should have built with config supplied annealing observers")

	dummyMessageObserver := new(observer.AnnealingMessageObserver)

	g.Expect(actualObservers[0]).To(BeAssignableToTypeOf(dummyMessageObserver),
		"Annealer should have built with config supplied annealing message observer")

	dummyAttributeObserver := new(observer.AnnealingAttributeObserver)

	g.Expect(actualObservers[1]).To(BeAssignableToTypeOf(dummyAttributeObserver),
		"Annealer should have built with config supplied annealing attribute observer")
}

type TestRegistereableExplorer struct {
	explorer.BaseExplorer
}

func (tre *TestRegistereableExplorer) WithName(name string) *TestRegistereableExplorer {
	tre.BaseExplorer.WithName(name)
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
				return new(TestRegistereableExplorer).WithName(config.Name)
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
