package solution

type VariableSummary struct {
	Name  string
	Value float64
}

type Summary []VariableSummary

func (s *Solution) Summarise() Summary {
	summary := make(Summary, 0)

	for _, variable := range s.DecisionVariables {
		variableSummary := VariableSummary{
			Name:  variable.Name,
			Value: variable.Value,
		}
		summary = append(summary, variableSummary)
	}
	return summary
}
