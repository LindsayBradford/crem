// Copyright (c) 2019 Australian Rivers Institute.

package data

import "github.com/LindsayBradford/crem/internal/pkg/config/data"

type EngineConfig struct {
	MetaData data.MetaDataConfig
	Engine   data.HttpServerConfig
}
