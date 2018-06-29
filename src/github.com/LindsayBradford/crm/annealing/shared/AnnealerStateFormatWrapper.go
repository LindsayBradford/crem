// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"fmt"
)

type AnnealerStateFormatWrapper struct {
	AnnealerToFormat Annealer

	MethodFormats map[string]string
}

const default_float64_format = "%f"
const default_uint_format = "%d"

func (this *AnnealerStateFormatWrapper) Initialise() *AnnealerStateFormatWrapper {
	this.MethodFormats = map[string]string{
		"Temperature":      default_float64_format,
		"CoolingFactor":    default_float64_format,
		"MaxIterations":    default_uint_format,
		"CurrentIteration": default_uint_format,
	}
	return this
}

func (this *AnnealerStateFormatWrapper) Wrapping(annealer Annealer) *AnnealerStateFormatWrapper {
	this.Wrap(annealer)
	return this
}

func (this *AnnealerStateFormatWrapper) Wrap(annealer Annealer) {
	this.AnnealerToFormat = annealer
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
