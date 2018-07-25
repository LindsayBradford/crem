// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	"fmt"
)

type SolutionTourerFormatWrapper struct {
	StateToFormat SolutionTourer
	MethodFormats map[string]string
}

const default_float64_format = "%f"
const default_bool_format = "%y"
const default_percent_format = "%f"

func (this *SolutionTourerFormatWrapper) Initialise() *SolutionTourerFormatWrapper {
	this.MethodFormats = map[string]string{
		"ObjectiveValue":         default_float64_format,
		"ChangeInObjectiveValue": default_float64_format,
		"ChangeIsDesirable":      default_bool_format,
		"ChangeAccepted":         default_bool_format,
		"AcceptanceProbability":  default_percent_format,
	}
	return this
}

func (this *SolutionTourerFormatWrapper) Wrapping(tourer SolutionTourer) *SolutionTourerFormatWrapper {
	this.Wrap(tourer)
	return this
}

func (this *SolutionTourerFormatWrapper) Wrap(tourer SolutionTourer) {
	this.StateToFormat = tourer
}

func (this *SolutionTourerFormatWrapper) ObjectiveValue() string {
	return this.applyFormatting("ObjectiveValue", this.StateToFormat.ObjectiveValue())
}

func (this *SolutionTourerFormatWrapper) ChangeInObjectiveValue() string {
	return this.applyFormatting("ChangeInObjectiveValue", this.StateToFormat.ChangeInObjectiveValue())
}

func (this *SolutionTourerFormatWrapper) ChangeIsDesirable() string {
	return this.applyFormatting("ChangeIsDesirable", this.StateToFormat.ChangeIsDesirable())
}

func (this *SolutionTourerFormatWrapper) ChangeAccepted() string {
	return this.applyFormatting("ChangeAccepted", this.StateToFormat.ChangeAccepted())
}

func (this *SolutionTourerFormatWrapper) AcceptanceProbability() string {
	return this.applyFormatting("AcceptanceProbability", this.StateToFormat.AcceptanceProbability())
}

func (this *SolutionTourerFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := this.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
