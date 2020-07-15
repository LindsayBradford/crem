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

func RetrieveConfigFromFile(configFilePath string) (*EngineConfig, error) {
	summary := decoderSummary{
		content:     configFilePath,
		contentType: file,
		decoder:     toml.DecodeFile,
	}
	return retrieveConfig(summary)
}

func RetrieveConfigFromString(tomlString string) (*EngineConfig, error) {
	summary := decoderSummary{
		content:     tomlString,
		contentType: text,
		decoder:     toml.Decode,
	}
	return retrieveConfig(summary)
}

func retrieveConfig(source decoderSummary) (*EngineConfig, error) {
	allErrors := errors2.New("configuration retrieval")

	var conf = defaultConfig()
	metaData, decodeErr := source.decoder(source.content, &conf)
	if decodeErr != nil {
		allErrors.Add(errors.Wrap(decodeErr, "failed retrieving config from "+source.contentType.String()))
	}
	if len(metaData.Undecoded()) > 0 {
		errorMsg := fmt.Sprintf("unrecognised configuration key(s) %q", metaData.Undecoded())
		allErrors.Add(errors.New(errorMsg))
	}
	conf.MetaData.FilePath = deriveFilePathFromSource(source)

	if allErrors.Size() > 0 {
		return nil, allErrors
	}

	return &conf, nil
}

func defaultConfig() EngineConfig {
	config := EngineConfig{
		data.MetaDataConfig{
			FilePath: "<no file path specified>",
		},
		data.HttpServerConfig{
			ApiPort:   8080,
			AdminPort: 8081,
		},
	}
	return config
}

func deriveFilePathFromSource(source decoderSummary) string {
	switch source.contentType {
	case file:
		return source.content
	default:
		return "<unspecified>"
	}
}
