// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	"github.com/LindsayBradford/crem/pkg/strings"
	"regexp"
	strings2 "strings"
)

// https://tools.ietf.org/html/rfc4180

const (
	idHeading = "Solution"
	separator = ", "
	newline   = "\n"
)

var variableHeadings = []string{idHeading}

var defaultConverter = strings.NewConverter().WithFloatingPointPrecision(3).PaddingZeros()

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

	for id, variables := range *summary {
		trimmedId := trimId(id)
		joinedVariableAttributes := joinAttributes(trimmedId, variables)
		builder.Add(joinedVariableAttributes).Add(newline)
	}

	return builder.String()
}

func deriveHeaders(summary *set.Summary) []string {
	exampleVariables := justSomeVariables(summary)
	headers := make([]string, len(exampleVariables)+1)

	headers[0] = idHeading
	for index, variable := range exampleVariables {
		headers[index+1] = variable.Name
	}

	return headers
}

func justSomeVariables(summary *set.Summary) []solution.VariableSummary {
	for _, variables := range *summary {
		return variables
	}
	return nil
}

func joinAttributes(id string, variables []solution.VariableSummary) string {
	joinedVariableValues := join(variableValueList(variables)...)
	joinedAttributes := join(id, joinedVariableValues)
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

func trimId(id string) string {
	iterationMatcher := regexp.MustCompile("\\d+/\\d+")
	trimmedId := iterationMatcher.FindString(id)
	prettifiedMatcher := regexp.MustCompile("/")
	prettifiedId := prettifiedMatcher.ReplaceAllString(trimmedId, " of ")
	return prettifiedId
}
