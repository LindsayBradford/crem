// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	_ "embed"
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"testing"

	. "github.com/onsi/gomega"
)

const richValidConfigFile = "testdata/RichValidConfig.toml"

//go:embed testdata/RichValidConfig.toml
var richValidConfig string

//go:embed testdata/RichInvalidSyntaxConfig.toml
var richInvalidConfig string

const (
	emptyTestFile = "testdata/EmptyConfig.toml"
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
		expectedCMetaDataFilePath        = richValidConfigFile
		expectedApiPort                  = uint64(3030)
		expectedAdminPort                = uint64(3031)
		expectedCacheMaximumAgeInSeconds = uint64(5)
		expectedJobQueueLength           = uint64(10)
	)

	// when
	config, retrieveError := RetrieveConfigFromFile(richValidConfigFile)
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
	const (
		expectedCMetaDataFilePath        = "<unspecified>"
		expectedApiPort                  = uint64(3030)
		expectedAdminPort                = uint64(3031)
		expectedCacheMaximumAgeInSeconds = uint64(5)
		expectedJobQueueLength           = uint64(10)
	)

	// when
	config, retrieveError := RetrieveConfigFromString(richValidConfig)
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

	// when
	_, retrieveError := RetrieveConfigFromString(richInvalidConfig)
	if retrieveError != nil {
		t.Log(retrieveError)
	}

	// then
	g.Expect(retrieveError).To(Not(BeNil()))
}
