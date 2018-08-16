// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	"fmt"
)

type ExplorerFormatWrapper struct {
	StateToFormat AnnealableExplorer
	MethodFormats map[string]string
}

const defaultFloat64Format = "%f"
const defaultBoolFormat = "%y"
const defaultPercentFormat = "%f"

func (wrapper *ExplorerFormatWrapper) Initialise() *ExplorerFormatWrapper {
	wrapper.MethodFormats = map[string]string{
		"ObjectiveValue":         defaultFloat64Format,
		"ChangeInObjectiveValue": defaultFloat64Format,
		"ChangeIsDesirable":      defaultBoolFormat,
		"ChangeAccepted":         defaultBoolFormat,
		"AcceptanceProbability":  defaultPercentFormat,
	}
	return wrapper
}

func (wrapper *ExplorerFormatWrapper) Wrapping(explorer Explorer) *ExplorerFormatWrapper {
	wrapper.Wrap(explorer)
	return wrapper
}

func (wrapper *ExplorerFormatWrapper) Wrap(explorer Explorer) {
	wrapper.StateToFormat = explorer
}

func (wrapper *ExplorerFormatWrapper) ObjectiveValue() string {
	return wrapper.applyFormatting("ObjectiveValue", wrapper.StateToFormat.ObjectiveValue())
}

func (wrapper *ExplorerFormatWrapper) ChangeInObjectiveValue() string {
	return wrapper.applyFormatting("ChangeInObjectiveValue", wrapper.StateToFormat.ChangeInObjectiveValue())
}

func (wrapper *ExplorerFormatWrapper) ChangeIsDesirable() string {
	return wrapper.applyFormatting("ChangeIsDesirable", wrapper.StateToFormat.ChangeIsDesirable())
}

func (wrapper *ExplorerFormatWrapper) ChangeAccepted() string {
	return wrapper.applyFormatting("ChangeAccepted", wrapper.StateToFormat.ChangeAccepted())
}

func (wrapper *ExplorerFormatWrapper) AcceptanceProbability() string {
	return wrapper.applyFormatting("AcceptanceProbability", wrapper.StateToFormat.AcceptanceProbability())
}

func (wrapper *ExplorerFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := wrapper.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
