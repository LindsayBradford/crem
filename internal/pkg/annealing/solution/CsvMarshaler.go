// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"fmt"
	"sort"
	strings2 "strings"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// https://tools.ietf.org/html/rfc4180

const (
	nameHeading          = "Name"
	valueHeading         = "Value"
	unitOfMeasureHeading = "UnitOfMeasure"
	separator            = ", "
	newline              = "\n"
)

var headings = []string{nameHeading, valueHeading, unitOfMeasureHeading}

type CsvDecisionVariableMarshaler struct{}

func (cm *CsvDecisionVariableMarshaler) Marshal(solution *Solution) ([]byte, error) {
	return cm.marshalDecisionVariables(solution.DecisionVariables)
}

func (cm *CsvDecisionVariableMarshaler) marshalDecisionVariables(variables variable.EncodeableDecisionVariables) ([]byte, error) {
	builder := new(strings.FluentBuilder)
	builder.Add(join(headings...)).Add(newline)

	for _, variable := range variables {
		joinedAttributes := join(
			variable.Name,
			toString(variable.Value),
			variable.Measure.String(),
		)
		builder.Add(joinedAttributes).Add(newline)
	}

	variableAsBytes := ([]byte)(builder.String())
	return variableAsBytes, nil
}

func quote(textToQuote string) string {
	return "\"" + textToQuote + "\""
}

func join(entries ...string) string {
	return strings2.Join(entries, separator)
}

func toString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

type CsvManagementActionMarshaler struct{}

func (cm *CsvManagementActionMarshaler) Marshal(solution *Solution) ([]byte, error) {
	return cm.marshalManagementActions(solution.PlanningUnitManagementActionsMap)
}

func (cm *CsvManagementActionMarshaler) marshalManagementActions(planningUnitActions map[PlanningUnitId]ManagementActions) ([]byte, error) {
	builder := new(strings.FluentBuilder)
	csvHeadings := cm.buildHeadings(planningUnitActions)

	joinedHeadings := join(csvHeadings...)
	builder.Add(joinedHeadings).Add(newline)

	sortedKeys := cm.sortPlanningUnitKeys(planningUnitActions) // TODO: sometimes missing a key if no actions active for PU

	for _, sortedKey := range sortedKeys {
		typedKey := PlanningUnitId(sortedKey)

		values := make([]string, len(headings))
		values[0] = sortedKey

		actions := planningUnitActions[typedKey]

		for headingIndex, heading := range headings {
			if heading == "PlanningUnit" {
				continue
			}
			actionValue := 0
			for _, action := range actions {
				actionAsString := string(action)
				if heading == actionAsString {
					actionValue = 1
				}
			}
			values[headingIndex] = toString(actionValue)
		}

		joinedActions := join(values...)
		builder.Add(joinedActions).Add(newline)
	}

	variableAsBytes := ([]byte)(builder.String())
	return variableAsBytes, nil
}

func (cm *CsvManagementActionMarshaler) sortPlanningUnitKeys(planningUnitActions map[PlanningUnitId]ManagementActions) []string {
	sortedKeys := make([]string, 0)
	for key := range planningUnitActions {
		sortedKeys = append(sortedKeys, string(key))
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func (cm *CsvManagementActionMarshaler) buildHeadings(planningUnitActions map[PlanningUnitId]ManagementActions) []string {
	headings := make([]string, 0)

	headingsAdded := make(map[string]bool, 0)

	headings = append(headings, "PlanningUnit")

	for _, actions := range planningUnitActions {
		for _, action := range actions {
			actionAsString := string(action)
			if _, hasEntry := headingsAdded[string(actionAsString)]; !hasEntry {
				headings = append(headings, actionAsString)
				headingsAdded[actionAsString] = true
			}
		}
	}

	return headings
}
