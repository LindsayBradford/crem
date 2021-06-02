// Copyright (c) 2019 Australian Rivers Institute.

package json

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	"regexp"
)

type Marshaler struct{}

const (
	newLinePrefix = ""
	indent        = "  "
)

type SolutionSummaries struct {
	SolutionSet string
	Solutions   []solution.Summary
}

func (m *Marshaler) Marshal(summary *set.Summary) ([]byte, error) {
	dataToMarshal := deriveSolutionSummaries(summary)

	marshalling, marshalError := json.MarshalIndent(dataToMarshal, newLinePrefix, indent)
	if marshalError != nil {
		return marshalling, marshalError
	}

	return marshalling, nil
}

func deriveSolutionSummaries(summary *set.Summary) SolutionSummaries {
	return SolutionSummaries{
		SolutionSet: deriveSetNameFor(summary),
		Solutions:   summary.AsSortedArray(),
	}
}

var nameMatcher = regexp.MustCompile("(.*) Solution.*")

func deriveSetNameFor(summary *set.Summary) string {
	firstKey := getFirstKey(summary)
	justTheName := nameMatcher.FindStringSubmatch(firstKey)[1]
	return justTheName
}

func getFirstKey(summary *set.Summary) string {
	for key := range *summary {
		return key
	}
	return ""
}
