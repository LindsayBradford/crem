// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	"fmt"
)

type SolutionExplorerFormatWrapper struct {
	StateToFormat SolutionExplorer
	MethodFormats map[string]string
}

const default_float64_format = "%f"
const default_bool_format = "%y"
const default_percent_format = "%f"

func (this *SolutionExplorerFormatWrapper) Initialise() *SolutionExplorerFormatWrapper {
	this.MethodFormats = map[string]string{
		"ObjectiveValue":         default_float64_format,
		"ChangeInObjectiveValue": default_float64_format,
		"ChangeIsDesirable":      default_bool_format,
		"ChangeAccepted":         default_bool_format,
		"AcceptanceProbability":  default_percent_format,
	}
	return this
}

func (this *SolutionExplorerFormatWrapper) Wrapping(explorer SolutionExplorer) *SolutionExplorerFormatWrapper {
	this.Wrap(explorer)
	return this
}

func (this *SolutionExplorerFormatWrapper) Wrap(explorer SolutionExplorer) {
	this.StateToFormat = explorer
}

func (this *SolutionExplorerFormatWrapper) ObjectiveValue() string {
	return this.applyFormatting("ObjectiveValue", this.StateToFormat.ObjectiveValue())
}

func (this *SolutionExplorerFormatWrapper) ChangeInObjectiveValue() string {
	return this.applyFormatting("ChangeInObjectiveValue", this.StateToFormat.ChangeInObjectiveValue())
}

func (this *SolutionExplorerFormatWrapper) ChangeIsDesirable() string {
	return this.applyFormatting("ChangeIsDesirable", this.StateToFormat.ChangeIsDesirable())
}

func (this *SolutionExplorerFormatWrapper) ChangeAccepted() string {
	return this.applyFormatting("ChangeAccepted", this.StateToFormat.ChangeAccepted())
}

func (this *SolutionExplorerFormatWrapper) AcceptanceProbability() string {
	return this.applyFormatting("AcceptanceProbability", this.StateToFormat.AcceptanceProbability())
}

func (this *SolutionExplorerFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := this.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
