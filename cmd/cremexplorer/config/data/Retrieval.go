// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	errors2 "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
)

type decoderSummary struct {
	contentType contentType
	content     string

	decoder func(data string, v interface{}) (toml.MetaData, error)
}

type contentType int

const (
	file contentType = iota
	text
)

func (st contentType) String() string {
	switch st {
	case file:
		return "file"
	case text:
		return "decoderSummary"
	default:
		return "undefined"
	}
}

func RetrieveConfigFromFile(configFilePath string) (*Config, error) {
	summary := decoderSummary{
		content:     configFilePath,
		contentType: file,
		decoder:     toml.DecodeFile,
	}
	return retrieveConfig(summary)
}

func RetrieveConfigFromString(tomlString string) (*Config, error) {
	summary := decoderSummary{
		content:     tomlString,
		contentType: text,
		decoder:     toml.Decode,
	}
	return retrieveConfig(summary)
}

func retrieveConfig(source decoderSummary) (*Config, error) {
	var conf = defaultConfig()
	metaData, decodeErr := source.decoder(source.content, &conf)
	if decodeErr != nil {
		return nil, errors.Wrap(decodeErr, "failed retrieving config from "+source.contentType.String())
	}
	if len(metaData.Undecoded()) > 0 {
		errorMsg := fmt.Sprintf("unrecognised configuration key(s) %q", metaData.Undecoded())
		return nil, errors.New(errorMsg)
	}
	conf.MetaData.FilePath = deriveFilePathFromSource(source)

	if checkErrors := checkMandatoryFields(&conf); checkErrors != nil {
		return nil, checkErrors
	}

	return &conf, nil
}

func defaultConfig() Config {
	config := Config{
		Scenario: ScenarioConfig{
			RunNumber:                  1,
			MaximumConcurrentRunNumber: 1,
			OutputPath:                 ".",
			Reporting: ReportingConfig{
				ReportEveryNumberOfIterations: 1,
			},
		},
	}
	return config
}

func checkMandatoryFields(config *Config) error {
	errors := errors2.New("Missing mandatory configuration")

	if config.Scenario.Name == "" {
		errors.AddMessage("Scenario.Name must be supplied")
	}

	if config.Scenario.RunNumber < 1 {
		errors.AddMessage("Scenario.RunNumber must be supplied with a value >= 1")
	}

	if config.Annealer.Type == data.UnspecifiedAnnealerType {
		errors.AddMessage("Annealer.Type must be supplied")
	}

	if config.Model.Type == "" {
		errors.AddMessage("Model.Type  must be supplied")
	}

	if errors.Size() > 0 {
		return errors
	}
	return nil
}

func deriveFilePathFromSource(source decoderSummary) string {
	switch source.contentType {
	case file:
		return source.content
	default:
		return "<unspecified>"
	}
}
