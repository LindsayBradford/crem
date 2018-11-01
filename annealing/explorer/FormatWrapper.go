// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"fmt"
)

type FormatWrapper struct {
	StateToFormat AnnealableExplorer
	MethodFormats map[string]string
}

const defaultFloat64Format = "%f"
const defaultBoolFormat = "%y"
const defaultPercentFormat = "%f"

func (wrapper *FormatWrapper) Initialise() *FormatWrapper {
	wrapper.MethodFormats = map[string]string{
		"ObjectiveValue":         defaultFloat64Format,
		"ChangeInObjectiveValue": defaultFloat64Format,
		"ChangeIsDesirable":      defaultBoolFormat,
		"ChangeAccepted":         defaultBoolFormat,
		"AcceptanceProbability":  defaultPercentFormat,
	}
	return wrapper
}

func (wrapper *FormatWrapper) Wrapping(explorer Explorer) *FormatWrapper {
	wrapper.Wrap(explorer)
	return wrapper
}

func (wrapper *FormatWrapper) Wrap(explorer Explorer) {
	wrapper.StateToFormat = explorer
}

func (wrapper *FormatWrapper) ObjectiveValue() string {
	return wrapper.applyFormatting("ObjectiveValue", wrapper.StateToFormat.ObjectiveValue())
}

func (wrapper *FormatWrapper) ChangeInObjectiveValue() string {
	return wrapper.applyFormatting("ChangeInObjectiveValue", wrapper.StateToFormat.ChangeInObjectiveValue())
}

func (wrapper *FormatWrapper) ChangeIsDesirable() string {
	return wrapper.applyFormatting("ChangeIsDesirable", wrapper.StateToFormat.ChangeIsDesirable())
}

func (wrapper *FormatWrapper) ChangeAccepted() string {
	return wrapper.applyFormatting("ChangeAccepted", wrapper.StateToFormat.ChangeAccepted())
}

func (wrapper *FormatWrapper) AcceptanceProbability() string {
	return wrapper.applyFormatting("AcceptanceProbability", wrapper.StateToFormat.AcceptanceProbability())
}

func (wrapper *FormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := wrapper.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
