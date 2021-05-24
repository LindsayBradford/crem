package solution

type Summary struct {
	Variables VariableSetSummary
	Actions   ActionSummary
}

type VariableSetSummary []VariableSummary

type VariableSummary struct {
	Name  string
	Value float64
}

type ActionSummary string

func (s *Solution) Summarise() Summary {
	return Summary{
		Variables: s.produceVariableSummary(),
		Actions:   s.produceActionSummary(),
	}
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

func (s *Solution) produceActionSummary() ActionSummary {
	return "#ReplaceMe" // TODO: get action summary out of solution.
}
