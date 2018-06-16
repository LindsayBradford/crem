/*
 * Copyright (c) 2018 Australian Rivers Institure. Author: Lindsay Bradford
 */

package annealing

import (
	"fmt"
)

type AnnealerStateFormatWrapper struct {
	annealerToFormat Annealer

	methodFormats map[string]string
}

const DEFAULT_FLOAT64_FORMAT = "%f"
const DEFAULT_UINT_FORMAT = "%d"

func NewAnnealerStateFormatWrapper(annealer Annealer) *AnnealerStateFormatWrapper {
	wrapper := AnnealerStateFormatWrapper{
		annealerToFormat: annealer,
		methodFormats: map[string]string{
			"Temperature":      DEFAULT_FLOAT64_FORMAT,
			"CoolingFactor":    DEFAULT_FLOAT64_FORMAT,
			"MaxIterations":    DEFAULT_UINT_FORMAT,
			"CurrentIteration": DEFAULT_UINT_FORMAT,
		},
	}

	return &wrapper
}

func (this *AnnealerStateFormatWrapper) Temperature() string {
	return this.applyFormatting("Temperature", this.annealerToFormat.Temperature())
}

func (this *AnnealerStateFormatWrapper) CoolingFactor() string {
	return this.applyFormatting("CoolingFactor", this.annealerToFormat.CoolingFactor())
}

func (this *AnnealerStateFormatWrapper) MaxIterations() string {
	return this.applyFormatting("MaxIterations", this.annealerToFormat.MaxIterations())
}

func (this *AnnealerStateFormatWrapper) CurrentIteration() string {
	return this.applyFormatting("CurrentIteration", this.annealerToFormat.CurrentIteration())
}

func (this *AnnealerStateFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := this.methodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
