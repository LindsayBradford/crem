// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package objectives

import (
	"fmt"
)

type ObjectiveManagerStateFormatWrapper struct {
	ObjectiveManagerToFormat ObjectiveManager
	MethodFormats map[string]string
}

const default_float64_format = "%f"
const default_bool_format = "%y"
const default_percent_format = "%f"

func (this *ObjectiveManagerStateFormatWrapper) Initialise() *ObjectiveManagerStateFormatWrapper {
	this.MethodFormats = map[string]string{
			"ObjectiveValue":         default_float64_format,
			"ChangeInObjectiveValue": default_float64_format,
			"ChangeIsDesirable":      default_bool_format,
			"ChangeAccepted":         default_bool_format,
			"AcceptanceProbability":  default_percent_format,
		}
	return this
}

func (this *ObjectiveManagerStateFormatWrapper) Wrapping(objectiveManager ObjectiveManager) *ObjectiveManagerStateFormatWrapper {
	this.Wrap(objectiveManager)
	return this
}

func (this *ObjectiveManagerStateFormatWrapper) Wrap(objectiveManager ObjectiveManager) {
	this.ObjectiveManagerToFormat = objectiveManager;
}

func (this *ObjectiveManagerStateFormatWrapper) ObjectiveValue() string {
	return this.applyFormatting("ObjectiveValue", this.ObjectiveManagerToFormat.ObjectiveValue())
}

func (this *ObjectiveManagerStateFormatWrapper) ChangeInObjectiveValue() string {
	return this.applyFormatting("ChangeInObjectiveValue", this.ObjectiveManagerToFormat.ChangeInObjectiveValue())
}

func (this *ObjectiveManagerStateFormatWrapper) ChangeIsDesirable() string {
	return this.applyFormatting("ChangeIsDesirable", this.ObjectiveManagerToFormat.ChangeIsDesirable())
}

func (this *ObjectiveManagerStateFormatWrapper) ChangeAccepted() string {
	return this.applyFormatting("ChangeAccepted", this.ObjectiveManagerToFormat.ChangeAccepted())
}

func (this *ObjectiveManagerStateFormatWrapper) AcceptanceProbability() string {
	return this.applyFormatting("AcceptanceProbability", this.ObjectiveManagerToFormat.AcceptanceProbability())
}

func (this *ObjectiveManagerStateFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := this.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
