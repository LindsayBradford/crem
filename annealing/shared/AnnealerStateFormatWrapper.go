// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"fmt"
)

type AnnealerFormatWrapper struct {
	AnnealerToFormat Annealer
	MethodFormats    map[string]string
}

const defaultStringFormat = "%s"
const defaultFloat64Format = "%f"
const defaultUintFormat = "%d"

func (wrapper *AnnealerFormatWrapper) Initialise() *AnnealerFormatWrapper {
	wrapper.MethodFormats = map[string]string{
		"iD":               defaultStringFormat,
		"Temperature":      defaultFloat64Format,
		"CoolingFactor":    defaultFloat64Format,
		"MaxIterations":    defaultUintFormat,
		"CurrentIteration": defaultUintFormat,
	}
	return wrapper
}

func (wrapper *AnnealerFormatWrapper) Wrapping(annealer Annealer) *AnnealerFormatWrapper {
	wrapper.Wrap(annealer)
	return wrapper
}

func (wrapper *AnnealerFormatWrapper) Wrap(annealer Annealer) {
	wrapper.AnnealerToFormat = annealer
}

func (wrapper *AnnealerFormatWrapper) Id() string {
	return wrapper.applyFormatting("Id", wrapper.AnnealerToFormat.Id())
}

func (wrapper *AnnealerFormatWrapper) Temperature() string {
	return wrapper.applyFormatting("Temperature", wrapper.AnnealerToFormat.Temperature())
}

func (wrapper *AnnealerFormatWrapper) CoolingFactor() string {
	return wrapper.applyFormatting("CoolingFactor", wrapper.AnnealerToFormat.CoolingFactor())
}

func (wrapper *AnnealerFormatWrapper) MaxIterations() string {
	return wrapper.applyFormatting("MaxIterations", wrapper.AnnealerToFormat.MaxIterations())
}

func (wrapper *AnnealerFormatWrapper) CurrentIteration() string {
	return wrapper.applyFormatting("CurrentIteration", wrapper.AnnealerToFormat.CurrentIteration())
}

func (wrapper *AnnealerFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := wrapper.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
