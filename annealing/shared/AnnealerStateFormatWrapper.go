// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"fmt"
)

type AnnealerFormatWrapper struct {
	AnnealerToFormat Annealer
	MethodFormats    map[string]string
}

const default_float64_format = "%f"
const default_uint_format = "%d"

func (this *AnnealerFormatWrapper) Initialise() *AnnealerFormatWrapper {
	this.MethodFormats = map[string]string{
		"Temperature":      default_float64_format,
		"CoolingFactor":    default_float64_format,
		"MaxIterations":    default_uint_format,
		"CurrentIteration": default_uint_format,
	}
	return this
}

func (this *AnnealerFormatWrapper) Wrapping(annealer Annealer) *AnnealerFormatWrapper {
	this.Wrap(annealer)
	return this
}

func (this *AnnealerFormatWrapper) Wrap(annealer Annealer) {
	this.AnnealerToFormat = annealer
}

func (this *AnnealerFormatWrapper) Temperature() string {
	return this.applyFormatting("Temperature", this.AnnealerToFormat.Temperature())
}

func (this *AnnealerFormatWrapper) CoolingFactor() string {
	return this.applyFormatting("CoolingFactor", this.AnnealerToFormat.CoolingFactor())
}

func (this *AnnealerFormatWrapper) MaxIterations() string {
	return this.applyFormatting("MaxIterations", this.AnnealerToFormat.MaxIterations())
}

func (this *AnnealerFormatWrapper) CurrentIteration() string {
	return this.applyFormatting("CurrentIteration", this.AnnealerToFormat.CurrentIteration())
}

func (this *AnnealerFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := this.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
