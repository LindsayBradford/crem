// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"io/ioutil"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	. "github.com/onsi/gomega"
)

const (
	emptyTestFile        = "testdata/EmptyConfig.toml"
	minimalValidTestFile = "testdata/MinimalValidConfig.toml"
	richValidTestFile    = "testdata/RichValidConfig.toml"
	richInvalidTestFile  = "testdata/RichInvalidConfig.toml"

	expectedScenarioName = "testScenario"
	testAnnealerType     = "Kirkpatrick"
	expectedModelType    = "TestModel"
)

func TestRetrieveConfigFromFile_EmptyConfig_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	config, retrieveError := RetrieveConfigFromFile(emptyTestFile)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(Not(BeNil()))
	g.Expect(config).To(BeNil())
}

func TestRetrieveConfigFromFile_MinimalValidConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedAnnealerType := data.AnnealerType{Value: testAnnealerType}

	// when
	config, retrieveError := RetrieveConfigFromFile(minimalValidTestFile)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(BeNil())
	g.Expect(config.Scenario.Name).To(Equal(expectedScenarioName))
	g.Expect(config.Annealer.Type).To(Equal(expectedAnnealerType))
	g.Expect(config.Model.Type).To(Equal(expectedModelType))
}

func TestRetrieveConfigFromString_MinimalValidConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configText := readTestFileAsText(minimalValidTestFile)
	expectedAnnealerType := data.AnnealerType{testAnnealerType}

	// when
	config, retrieveError := RetrieveConfigFromString(configText)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(BeNil())
	g.Expect(config.Scenario.Name).To(Equal(expectedScenarioName))
	g.Expect(config.Annealer.Type).To(Equal(expectedAnnealerType))
	g.Expect(config.Model.Type).To(Equal(expectedModelType))
}

func TestRetrieveConfigFromFile_RichValidConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	config, retrieveError := RetrieveConfigFromFile(richValidTestFile)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(BeNil())
	g.Expect(config.Scenario.Name).To(Equal(expectedScenarioName))
}

func TestRetrieveConfigFromString_RichValidConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configText := readTestFileAsText(richValidTestFile)

	// when
	config, retrieveError := RetrieveConfigFromString(configText)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(BeNil())
	g.Expect(config.Scenario.Name).To(Equal(expectedScenarioName))
}

func TestRetrieveConfigFromString_RichInvalidConfig_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configText := readTestFileAsText(richInvalidTestFile)

	// when
	_, retrieveError := RetrieveConfigFromString(configText)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(Not(BeNil()))
}

func readTestFileAsText(filePath string) string {
	if b, err := ioutil.ReadFile(filePath); err == nil {
		return string(b)
	}
	return "error reading file"
}
