// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"fmt"
)

type FormatWrapper struct {
	StateToFormat Observable
	MethodFormats map[string]string
}

const defaultFloat64Format = "%f"
const defaultBoolFormat = "%y"
const defaultPercentFormat = "%f"

func (fw *FormatWrapper) Initialise() *FormatWrapper {
	fw.MethodFormats = map[string]string{
		"ObjectiveValue":         defaultFloat64Format,
		"ChangeInObjectiveValue": defaultFloat64Format,
		"ChangeIsDesirable":      defaultBoolFormat,
		"ChangeAccepted":         defaultBoolFormat,
		"AcceptanceProbability":  defaultPercentFormat,
	}
	return fw
}

func (fw *FormatWrapper) Wrapping(explorer Observable) *FormatWrapper {
	fw.Wrap(explorer)
	return fw
}

func (fw *FormatWrapper) Wrap(explorer Observable) {
	fw.StateToFormat = explorer
}

func (fw *FormatWrapper) ObjectiveValue() string {
	return fw.applyFormatting("ObjectiveValue", fw.StateToFormat.ObjectiveValue())
}

func (fw *FormatWrapper) ChangeInObjectiveValue() string {
	return fw.applyFormatting("ChangeInObjectiveValue", fw.StateToFormat.ChangeInObjectiveValue())
}

func (fw *FormatWrapper) ChangeIsDesirable() string {
	return fw.applyFormatting("ChangeIsDesirable", fw.StateToFormat.ChangeIsDesirable())
}

func (fw *FormatWrapper) ChangeAccepted() string {
	return fw.applyFormatting("ChangeAccepted", fw.StateToFormat.ChangeAccepted())
}

func (fw *FormatWrapper) AcceptanceProbability() string {
	return fw.applyFormatting("AcceptanceProbability", fw.StateToFormat.AcceptanceProbability())
}

func (fw *FormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := fw.MethodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}
