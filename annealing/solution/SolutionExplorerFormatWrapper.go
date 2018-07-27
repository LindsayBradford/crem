// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	"fmt"
)

type SolutionExplorerFormatWrapper struct {
	StateToFormat SolutionExplorer
	MethodFormats map[string]string
}

const defaultFloat64Format = "%f"
const defaultBoolFormat = "%y"
const defaultPercentFormat = "%f"

func (wrapper *SolutionExplorerFormatWrapper) Initialise() *SolutionExplorerFormatWrapper {
	wrapper.MethodFormats = map[string]string{
		"ObjectiveValue":         defaultFloat64Format,
		"ChangeInObjectiveValue": defaultFloat64Format,
		"ChangeIsDesirable":      defaultBoolFormat,
		"ChangeAccepted":         defaultBoolFormat,
		"AcceptanceProbability":  defaultPercentFormat,
	}
	return wrapper
}

func (wrapper *SolutionExplorerFormatWrapper) Wrapping(explorer SolutionExplorer) *SolutionExplorerFormatWrapper {
	wrapper.Wrap(explorer)
	return wrapper
}

func (wrapper *SolutionExplorerFormatWrapper) Wrap(explorer SolutionExplorer) {
	wrapper.StateToFormat = explorer
}

func (wrapper *SolutionExplorerFormatWrapper) ObjectiveValue() string {
	return wrapper.applyFormatting("ObjectiveValue", wrapper.StateToFormat.ObjectiveValue())
}

func (wrapper *SolutionExplorerFormatWrapper) ChangeInObjectiveValue() string {
	return wrapper.applyFormatting("ChangeInObjectiveValue", wrapper.StateToFormat.ChangeInObjectiveValue())
}

func (wrapper *SolutionExplorerFormatWrapper) ChangeIsDesirable() string {
	return wrapper.applyFormatting("ChangeIsDesirable", wrapper.StateToFormat.ChangeIsDesirable())
}

func (wrapper *SolutionExplorerFormatWrapper) ChangeAccepted() string {
	return wrapper.applyFormatting("ChangeAccepted", wrapper.StateToFormat.ChangeAccepted())
}

func (wrapper *SolutionExplorerFormatWrapper) AcceptanceProbability() string {
	return wrapper.applyFormatting("AcceptanceProbability", wrapper.StateToFormat.AcceptanceProbability())
}

func (wrapper *SolutionExplorerFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := wrapper.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
