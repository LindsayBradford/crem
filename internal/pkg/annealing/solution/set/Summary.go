package set

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"regexp"
	"sort"
	"strings"
)

type Summary map[string]solution.Summary

func (s Summary) Id() string {
	baseId := s.justSomeId()

	iterationMatcher := regexp.MustCompile(`Solution \(.+\)`)
	id := iterationMatcher.ReplaceAllString(baseId, "Summary")

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

type sortableSummaries []solution.Summary

func (v sortableSummaries) Len() int {
	return len(v)
}

func (v sortableSummaries) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v sortableSummaries) Less(i, j int) bool {
	return v[i].SortIndex < v[j].SortIndex
}

func (s Summary) AsSortedArray() []solution.Summary {
	var justTheSummaries sortableSummaries
	for _, value := range s {
		justTheSummaries = append(justTheSummaries, value)
	}
	sort.Sort(justTheSummaries)
	return justTheSummaries
}
