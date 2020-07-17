// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

const defaultServerConfigPath = "config/server.toml"

const (
	DefaultApiPort   = uint64(8080)
	DefaultAdminPort = uint64(8081)
)

type HttpServerConfig struct {
	AdminPort                uint64
	ApiPort                  uint64
	CacheMaximumAgeInSeconds uint64
	JobQueueLength           uint64

	Logger LoggingConfig
}

func RetrieveHttpServer(configFilePath string) (*HttpServerConfig, error) {
	if configFilePath == "" {
		configFilePath = defaultServerConfigPath
		if defaultServerConfigFileNotSupplied() {
			return embeddedDefaultHttpServerConfig()
		}
	}

	return retrieveHttpServerFromFile(configFilePath)
}

func defaultServerConfigFileNotSupplied() bool {
	_, err := os.Stat(defaultServerConfigPath)
	return os.IsNotExist(err)
}

func embeddedDefaultHttpServerConfig() (*HttpServerConfig, error) {
	config := HttpServerConfig{ApiPort: DefaultApiPort, AdminPort: DefaultAdminPort}
	return &config, nil
}

func retrieveHttpServerFromFile(configFilePath string) (*HttpServerConfig, error) {
	var conf HttpServerConfig
	metaData, decodeErr := toml.DecodeFile(configFilePath, &conf)

	if decodeErr != nil {
		return nil, errors.Wrap(decodeErr, "failed retrieving config from file")
	}
	if len(metaData.Undecoded()) > 0 {
		errorMsg := fmt.Sprintf("unrecognised configuration key(s) %q", metaData.Undecoded())
		return nil, errors.New(errorMsg)
	}
	return &conf, nil
}
