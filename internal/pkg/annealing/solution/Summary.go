package solution

type Summary struct {
	SortIndex uint64 `json:"-"`
	Id        string
	Variables VariableSetSummary
	Actions   ActionSummary
	Note      string
}

func (s *Solution) Summarise() *Summary {
	return &Summary{
		SortIndex: 0,
		Id:        "",
		Variables: s.produceVariableSummary(),
		Actions:   s.produceActionSummary(),
		Note:      "",
	}
}

func (s *Summary) WithId(id string) *Summary {
	s.Id = id
	return s
}

func (s *Summary) Noting(note string) *Summary {
	s.Note = note
	return s
}

func (s *Summary) WithSortOrder(sortIndex uint64) *Summary {
	s.SortIndex = sortIndex
	return s
}

type VariableSetSummary []VariableSummary

type VariableSummary struct {
	Name  string
	Value float64
}

func (s *Solution) produceVariableSummary() VariableSetSummary {
	summary := make(VariableSetSummary, 0)

	for _, variable := range s.DecisionVariables {
		variableSummary := VariableSummary{
			Name:  variable.Name,
			Value: variable.Value,
		}
		summary = append(summary, variableSummary)
	}
	return summary
}

type ActionSummary string

func (s *Solution) produceActionSummary() ActionSummary {
	return ActionSummary(s.EncodedActions)
}
