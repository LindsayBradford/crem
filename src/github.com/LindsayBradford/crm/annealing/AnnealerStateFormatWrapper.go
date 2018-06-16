/*
 * Copyright (c) 2017 Australian Rivers Institute. Author: Lindsay Bradford
 */

package annealing

import (
	"fmt"
)

type AnnealerStateFormatWrapper struct {
	AnnealerToFormat Annealer

	MethodFormats map[string]string
}

const DEFAULT_FLOAT64_FORMAT = "%f"
const DEFAULT_UINT_FORMAT = "%d"

func (this *AnnealerStateFormatWrapper) Initialise() *AnnealerStateFormatWrapper {
	this.MethodFormats = map[string]string{
			"Temperature":      DEFAULT_FLOAT64_FORMAT,
			"CoolingFactor":    DEFAULT_FLOAT64_FORMAT,
			"MaxIterations":    DEFAULT_UINT_FORMAT,
			"CurrentIteration": DEFAULT_UINT_FORMAT,
		}
	return this
}

func (this *AnnealerStateFormatWrapper) Wrapping(annealer Annealer) *AnnealerStateFormatWrapper {
	this.Wrap(annealer)
	return this
}

func (this *AnnealerStateFormatWrapper) Wrap(annealer Annealer) {
	this.AnnealerToFormat = annealer;
}

func (this *AnnealerStateFormatWrapper) Temperature() string {
	return this.applyFormatting("Temperature", this.AnnealerToFormat.Temperature())
}

func (this *AnnealerStateFormatWrapper) CoolingFactor() string {
	return this.applyFormatting("CoolingFactor", this.AnnealerToFormat.CoolingFactor())
}

func (this *AnnealerStateFormatWrapper) MaxIterations() string {
	return this.applyFormatting("MaxIterations", this.AnnealerToFormat.MaxIterations())
}

func (this *AnnealerStateFormatWrapper) CurrentIteration() string {
	return this.applyFormatting("CurrentIteration", this.AnnealerToFormat.CurrentIteration())
}

func (this *AnnealerStateFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := this.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
