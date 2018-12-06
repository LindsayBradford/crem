// Copyright (c) 2018 Australian Rivers Institute.

package wrapper

import (
	"fmt"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
)

type FormatWrapper struct {
	AnnealerToFormat annealing.Observable
	MethodFormats    map[string]string
}

const defaultStringFormat = "%s"
const defaultFloat64Format = "%f"
const defaultUintFormat = "%d"

func (wrapper *FormatWrapper) Initialise() *FormatWrapper {
	wrapper.MethodFormats = map[string]string{
		"iD":                defaultStringFormat,
		"Temperature":       defaultFloat64Format,
		"CoolingFactor":     defaultFloat64Format,
		"MaximumIterations": defaultUintFormat,
		"CurrentIteration":  defaultUintFormat,
	}
	return wrapper
}

func (wrapper *FormatWrapper) Wrapping(annealer annealing.Annealer) *FormatWrapper {
	wrapper.Wrap(annealer)
	return wrapper
}

func (wrapper *FormatWrapper) Wrap(annealer annealing.Observable) {
	wrapper.AnnealerToFormat = annealer
}

func (wrapper *FormatWrapper) Id() string {
	return wrapper.applyFormatting("Id", wrapper.AnnealerToFormat.Id())
}

func (wrapper *FormatWrapper) Temperature() string {
	return wrapper.applyFormatting("Temperature", wrapper.AnnealerToFormat.Temperature())
}

func (wrapper *FormatWrapper) CoolingFactor() string {
	return wrapper.applyFormatting("CoolingFactor", wrapper.AnnealerToFormat.CoolingFactor())
}

func (wrapper *FormatWrapper) MaximumIterations() string {
	return wrapper.applyFormatting("MaximumIterations", wrapper.AnnealerToFormat.MaximumIterations())
}

func (wrapper *FormatWrapper) CurrentIteration() string {
	return wrapper.applyFormatting("CurrentIteration", wrapper.AnnealerToFormat.CurrentIteration())
}

func (wrapper *FormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := wrapper.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
