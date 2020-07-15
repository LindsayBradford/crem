// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	emptyTestFile             = "testdata/EmptyConfig.toml"
	richValidTestFile         = "testdata/RichValidConfig.toml"
	richInvalidSyntaxTestFile = "testdata/RichInvalidSyntaxConfig.toml"
)

func TestRetrieveConfigFromFile_MissingConfig_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	config, retrieveError := RetrieveConfigFromFile("ThisFileDoesNotExist.txt")
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError.Error()).To(ContainSubstring("failed retrieving config from file"))
	g.Expect(config).To(BeNil())
}

func TestRetrieveConfigFromFile_EmptyConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	config, retrieveError := RetrieveConfigFromFile(emptyTestFile)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(BeNil())

	g.Expect(config.Engine.ApiPort).To(Equal(data.DefaultApiPort))
	g.Expect(config.Engine.AdminPort).To(Equal(data.DefaultAdminPort))
}

func TestRetrieveConfigFromFile_RichValidConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	const (
		expectedCMetaDataFilePath        = richValidTestFile
		expectedApiPort                  = uint64(3030)
		expectedAdminPort                = uint64(3031)
		expectedCacheMaximumAgeInSeconds = uint64(5)
		expectedJobQueueLength           = uint64(10)
	)

	// when
	config, retrieveError := RetrieveConfigFromFile(richValidTestFile)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(BeNil())
	g.Expect(config.MetaData.FilePath).To(Equal(expectedCMetaDataFilePath))

	g.Expect(config.Engine.ApiPort).To(Equal(expectedApiPort))
	g.Expect(config.Engine.AdminPort).To(Equal(expectedAdminPort))
	g.Expect(config.Engine.CacheMaximumAgeInSeconds).To(Equal(expectedCacheMaximumAgeInSeconds))
	g.Expect(config.Engine.JobQueueLength).To(Equal(expectedJobQueueLength))
}

func TestRetrieveConfigFromString_RichValidConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configText := readTestFileAsText(richValidTestFile)
	const (
		expectedCMetaDataFilePath        = "<unspecified>"
		expectedApiPort                  = uint64(3030)
		expectedAdminPort                = uint64(3031)
		expectedCacheMaximumAgeInSeconds = uint64(5)
		expectedJobQueueLength           = uint64(10)
	)

	// when
	config, retrieveError := RetrieveConfigFromString(configText)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(BeNil())

	g.Expect(config.MetaData.FilePath).To(Equal(expectedCMetaDataFilePath))

	g.Expect(config.Engine.ApiPort).To(Equal(expectedApiPort))
	g.Expect(config.Engine.AdminPort).To(Equal(expectedAdminPort))
	g.Expect(config.Engine.CacheMaximumAgeInSeconds).To(Equal(expectedCacheMaximumAgeInSeconds))
	g.Expect(config.Engine.JobQueueLength).To(Equal(expectedJobQueueLength))
}

func TestRetrieveConfigFromString_RichInvalidSyntaxConfig_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configText := readTestFileAsText(richInvalidSyntaxTestFile)

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
