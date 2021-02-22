package data

import (
	"fmt"
	"github.com/BurntSushi/toml"
	errors2 "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
	"strings"
)

func RetrieveScenarioConfigFromFile(configFilePath string) (*ScenarioConfig, error) {
	summary := decoderSummary{
		content:     configFilePath,
		contentType: file,
		decoder:     toml.DecodeFile,
	}
	return retrieveScenarioConfigFromFile(summary)
}

func RetrieveScenarioConfigFromString(tomlString string) (*ScenarioConfig, error) {
	summary := decoderSummary{
		content:     tomlString,
		contentType: text,
		decoder:     toml.Decode,
	}
	return retrieveScenarioConfig(summary)
}

func retrieveScenarioConfig(source decoderSummary) (*ScenarioConfig, error) {
	allErrors := errors2.New("configuration retrieval")

	var conf = defaultScenarioConfig()
	_, decodeErr := source.decoder(source.content, &conf)
	if decodeErr != nil {
		allErrors.Add(errors.Wrap(decodeErr, "failed retrieving config from "+source.contentType.String()))
	}

	if allErrors.Size() > 0 {
		return nil, allErrors
	}

	return &conf, nil
}

func defaultScenarioConfig() ScenarioConfig {
	return ScenarioConfig{}
}

func retrieveScenarioConfigFromFile(source decoderSummary) (*ScenarioConfig, error) {
	allErrors := errors2.New("scenario configuration retrieval")

	var conf = defaultScenarioConfig()
	metaData, decodeErr := source.decoder(source.content, &conf)
	if decodeErr != nil {
		allErrors.Add(errors.Wrap(decodeErr, "failed retrieving config from "+source.contentType.String()))
	}

	var filteredKeys []string
	for _, rawUnknownKey := range metaData.Undecoded() {
		unknownKey := rawUnknownKey.String()
		if unknownKey == "" || strings.TrimSpace(unknownKey) == "" {
			continue
		}
		if strings.HasPrefix(unknownKey, "Scenario") {
			continue
		}
		if strings.HasPrefix(unknownKey, "Annealer") {
			continue
		}
		filteredKeys = append(filteredKeys, unknownKey)
	}

	if len(filteredKeys) > 0 {
		errorMsg := fmt.Sprintf("unrecognised configuration key(s) %q", filteredKeys)
		allErrors.Add(errors.New(errorMsg))
	}

	if checkErrors := checkMandatoryFields(&conf); checkErrors != nil {
		allErrors.Add(checkErrors)
	}

	if allErrors.Size() > 0 {
		return nil, allErrors
	}

	return &conf, nil
}

func checkMandatoryFields(config *ScenarioConfig) error {
	errors := errors2.New("Missing mandatory configuration")

	if config.Scenario.Name == "" {
		errors.AddMessage("Scenario.Name must be supplied")
	}

	if config.Model.Type == "" {
		errors.AddMessage("Model.Type  must be supplied")
	}

	if errors.Size() > 0 {
		return errors
	}
	return nil
}
