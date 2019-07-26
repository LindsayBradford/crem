// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
)

type ReportingConfig struct {
	ReportEveryNumberOfIterations uint64
	data.LoggingConfig
}
