package set

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"regexp"
	"strings"
)

type Summary map[string]solution.Summary

func (s Summary) FileNameSafeId() string {
	baseId := s.justSomeId()

	safeId := strings.Replace(baseId, " ", "", -1)

	iterationMatcher := regexp.MustCompile(`\(.+/.+\)`)
	safeId = iterationMatcher.ReplaceAllString(safeId, "")

	return safeId
}

func (s Summary) justSomeId() string {
	for key := range s {
		return key
	}
	return ""
}
