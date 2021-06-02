// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	"github.com/LindsayBradford/crem/pkg/strings"
	strings2 "strings"
)

// https://tools.ietf.org/html/rfc4180

const (
	idHeading      = "Solution"
	actionsHeading = "Actions"
	summaryHeading = "Summary"
	separator      = ", "
	newline        = "\n"
)

var (
	variableHeadings = []string{idHeading}
	defaultConverter = strings.NewConverter().WithFloatingPointPrecision(3).PaddingZeros()
)

type SummaryMarshaler struct{}

func (cm *SummaryMarshaler) Marshal(summary *set.Summary) ([]byte, error) {
	return cm.marshalSummary(summary)
}

func (cm *SummaryMarshaler) marshalSummary(summary *set.Summary) ([]byte, error) {
	csvStringAsBytes := ([]byte)(cm.summaryToCsvString(summary))
	return csvStringAsBytes, nil
}

func (cm *SummaryMarshaler) summaryToCsvString(summary *set.Summary) string {
	headers := deriveHeaders(summary)

	builder := new(strings.FluentBuilder)
	builder.
		Add(join(headers...)).
		Add(newline)

	summarySet := make([]string, 0)
	for _, solutionSummary := range summary.AsSortedArray() {
		summaryId := solutionSummary.Id
		note := solutionSummary.Note
		summarySet = append(summarySet, joinAttributes(summaryId, solutionSummary.Variables, solutionSummary.Actions, note))
	}

	for _, sortedSummary := range summarySet {
		builder.Add(sortedSummary).Add(newline)
	}

	return builder.String()
}

func deriveHeaders(summary *set.Summary) []string {
	exampleVariables := justSomeVariables(summary)

	headingNumber := len(exampleVariables) + 3
	headers := make([]string, headingNumber)

	headers[0] = idHeading
	for index, variable := range exampleVariables {
		headers[index+1] = variable.Name
	}
	headers[headingNumber-2] = actionsHeading
	headers[headingNumber-1] = summaryHeading

	return headers
}

func justSomeVariables(summary *set.Summary) solution.VariableSetSummary {
	for _, currentSummary := range *summary {
		return currentSummary.Variables
	}
	return nil
}

func joinAttributes(id string, variables []solution.VariableSummary, actions solution.ActionSummary, note string) string {
	joinedVariableValues := join(variableValueList(variables)...)
	joinedAttributes := join(id, joinedVariableValues, string(actions), note)
	return joinedAttributes
}

func variableValueList(variables []solution.VariableSummary) []string {
	values := make([]string, len(variables))
	for index, variable := range variables {
		values[index] = defaultConverter.Convert(variable.Value)
	}
	return values
}

func join(entries ...string) string {
	return strings2.Join(entries, separator)
}
