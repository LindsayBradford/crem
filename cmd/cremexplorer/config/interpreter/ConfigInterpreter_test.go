// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"io/ioutil"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	. "github.com/onsi/gomega"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
)

func TestConfigInterpreter_NewConfigInterpreter_MinimalConfigNoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configText := readTestFileAsText("testdata/MinimalValidConfig.toml")
	configUnderTest, configError := data.RetrieveConfigFromString(configText)

	// when
	interpreterUnderTest := NewInterpreter().Interpret(configUnderTest)

	// then
	g.Expect(configError).To(BeNil())
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
	actualScenario := interpreterUnderTest.Scenario()
	g.Expect(actualScenario).To(BeAssignableToTypeOf(&scenario.BaseScenario{}))
}

func TestConfigInterpreter_BrokenMandatoryFields_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configText := readTestFileAsText("testdata/BrokenMandatoryFieldConfig.toml")
	configUnderTest, configError := data.RetrieveConfigFromString(configText)
	g.Expect(configError).To(Not(BeNil()))

	// when
	interpreterUnderTest := NewInterpreter().Interpret(configUnderTest)

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(Not(BeNil()))
}

func readTestFileAsText(filePath string) string {
	if b, err := ioutil.ReadFile(filePath); err == nil {
		return string(b)
	}
	return "error reading file"
}
