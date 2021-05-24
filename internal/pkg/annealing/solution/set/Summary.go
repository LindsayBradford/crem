package set

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"regexp"
	"strings"
)

type Summary map[string]solution.Summary

func (s Summary) Id() string {
	baseId := s.justSomeId()

	iterationMatcher := regexp.MustCompile(`Solution \(.+\)`)
	id := iterationMatcher.ReplaceAllString(baseId, "VariableSetSummary")

	return id
}

func (s Summary) FileNameSafeId() string {
	baseId := s.justSomeId()

	safeId := strings.Replace(baseId, " ", "", -1)

	iterationMatcher := regexp.MustCompile(`Solution\(.+\)`)
	safeId = iterationMatcher.ReplaceAllString(safeId, "")
	safeId = strings.Replace(safeId, "/", "_of_", -1)

	return safeId
}

func (s Summary) justSomeId() string {
	for key := range s {
		return key
	}
	return ""
}
