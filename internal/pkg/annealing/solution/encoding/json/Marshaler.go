// Copyright (c) 2019 Australian Rivers Institute.

package json

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"regexp"
)

type Marshaler struct{}

const (
	newLinePrefix = ""
	indent        = "  "
)

func (m *Marshaler) Marshal(s *solution.Solution) ([]byte, error) {
	unalteredMarshalling, marshalError := json.MarshalIndent(s, newLinePrefix, indent)
	if marshalError != nil {
		return unalteredMarshalling, marshalError
	}

	alteredMarshalling := replaceText(
		string(unalteredMarshalling),
		solution.DefaultPlanningUnitHeading,
		s.PlanningUnitHeading(),
	)

	return []byte(alteredMarshalling), nil
}

func replaceText(originalText string, stringToReplace string, replacementString string) string {
	if stringToReplace == replacementString {
		return originalText
	}
	replaceExpression := regexp.MustCompile(stringToReplace)
	return replaceExpression.ReplaceAllString(originalText, replacementString)
}
