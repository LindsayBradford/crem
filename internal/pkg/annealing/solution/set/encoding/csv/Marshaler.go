// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	"github.com/LindsayBradford/crem/pkg/strings"
	"regexp"
	"sort"
	"strconv"
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

var numberMatcher *regexp.Regexp

func init() {
	numberMatcher = regexp.MustCompile(`(\d+) of `)
}

type sortableSummaries []string

func (v sortableSummaries) Len() int {
	return len(v)
}

func (v sortableSummaries) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v sortableSummaries) Less(i, j int) bool {
	const indexOfNumberMatch = 1
	numberMatchAtI := numberMatcher.FindStringSubmatch(v[i])
	numberMatchAtJ := numberMatcher.FindStringSubmatch(v[j])

	// As-Is entry should always be first, and will not be caught be a number matching regular expression

	if numberMatchAtI == nil {
		return true
	} else if numberMatchAtJ == nil {
		return false
	}

	numberAtI, _ := strconv.ParseInt(numberMatchAtI[indexOfNumberMatch], 10, 32)
	numberAtJ, _ := strconv.ParseInt(numberMatchAtJ[indexOfNumberMatch], 10, 32)

	return numberAtI < numberAtJ
}

func (cm *SummaryMarshaler) summaryToCsvString(summary *set.Summary) string {
	headers := deriveHeaders(summary)

	builder := new(strings.FluentBuilder)
	builder.
		Add(join(headers...)).
		Add(newline)

	summarySet := make([]string, 0)
	for id, variables := range *summary {
		trimmedId := trimId(id)
		summarySet = append(summarySet, joinAttributes(trimmedId, variables))
	}

	sort.Sort(sortableSummaries(summarySet))

	for _, sortedSummary := range summarySet {
		builder.Add(sortedSummary).Add(newline)
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
	if strings2.Contains(id, "As-Is") {
		return trimAsIsId(id)
	}
	return trimNumberedId(id)
}

func trimAsIsId(id string) string {
	return "As-Is"
}

func trimNumberedId(id string) string {
	iterationMatcher := regexp.MustCompile("\\d+/\\d+")
	trimmedId := iterationMatcher.FindString(id)
	prettifiedMatcher := regexp.MustCompile("/")
	prettifiedId := prettifiedMatcher.ReplaceAllString(trimmedId, " of ")

	return prettifiedId
}
