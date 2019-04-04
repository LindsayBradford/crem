// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"fmt"
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

var variableHeadings = []string{nameHeading, valueHeading, unitOfMeasureHeading}

const (
	planningUnitHeading = "PlanningUnit"
	inactiveActionValue = "0"
	activeActionValue   = "1"
)

type CsvDecisionVariableMarshaler struct{}

func (cm *CsvDecisionVariableMarshaler) Marshal(solution *Solution) ([]byte, error) {
	return cm.marshalDecisionVariables(solution.DecisionVariables)
}

func (cm *CsvDecisionVariableMarshaler) marshalDecisionVariables(variables variable.EncodeableDecisionVariables) ([]byte, error) {
	csvStringAsBytes := ([]byte)(cm.decisionVariablesToCsvString(variables))
	return csvStringAsBytes, nil
}

func (cm *CsvDecisionVariableMarshaler) decisionVariablesToCsvString(variables variable.EncodeableDecisionVariables) string {
	builder := new(strings.FluentBuilder)
	builder.Add(join(variableHeadings...)).Add(newline)

	for _, variable := range variables {
		joinedVariableAttributes := joinAttributes(variable)
		builder.Add(joinedVariableAttributes).Add(newline)
	}

	return builder.String()
}

func joinAttributes(variable variable.EncodeableDecisionVariable) string {
	return join(
		variable.Name,
		toString(variable.Value),
		variable.Measure.String(),
	)
}

func join(entries ...string) string {
	return strings2.Join(entries, separator)
}

func toString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

type CsvManagementActionMarshaler struct{}

func (cm *CsvManagementActionMarshaler) Marshal(solution *Solution) ([]byte, error) {
	return cm.marshalManagementActions(solution)
}

func (cm *CsvManagementActionMarshaler) marshalManagementActions(solution *Solution) ([]byte, error) {
	csvStringAsBytes := ([]byte)(cm.csvEncodeManagementActions(solution))
	return csvStringAsBytes, nil
}

func (cm *CsvManagementActionMarshaler) csvEncodeManagementActions(solution *Solution) string {
	builder := new(strings.FluentBuilder)

	headings := csvEncodeActionHeadings(solution.ActiveManagementActions)
	builder.Add(join(headings...)).Add(newline)

	for _, planningUnit := range solution.PlanningUnits {
		actionRowValues := cm.buildActionCsvValuesForPlanningUnit(headings, planningUnit, solution)
		builder.Add(join(actionRowValues...)).Add(newline)
	}

	return builder.String()
}

func csvEncodeActionHeadings(planningUnitActions map[PlanningUnitId]ManagementActions) []string {
	headings := make([]string, 1)
	headings[0] = planningUnitHeading

	headingsAdded := make(map[ManagementActionType]bool, 0)
	for _, actions := range planningUnitActions {
		for _, action := range actions {
			if _, hasEntry := headingsAdded[action]; !hasEntry {
				headings = append(headings, string(action))
				headingsAdded[action] = true
			}
		}
	}

	return headings
}

func (cm *CsvManagementActionMarshaler) buildActionCsvValuesForPlanningUnit(
	actionHeadings []string, planningUnit PlanningUnitId, solution *Solution) []string {

	values := make([]string, len(actionHeadings))

	values[0] = string(planningUnit)

	if activeActions, unitHasActiveActions := solution.ActiveManagementActions[planningUnit]; unitHasActiveActions {
		for headingIndex, csvHeading := range actionHeadings {
			if shouldSkipColumnWith(csvHeading) {
				continue
			}

			actionValue := inactiveActionValue
			for _, action := range activeActions {
				if actionMatchesColumnNamed(action, csvHeading) {
					actionValue = activeActionValue
				}
			}

			values[headingIndex] = actionValue
		}
	} else {
		for headingIndex, csvHeading := range actionHeadings {
			if shouldSkipColumnWith(csvHeading) {
				continue
			}
			values[headingIndex] = inactiveActionValue
		}
	}
	return values
}

func shouldSkipColumnWith(csvHeading string) bool {
	return csvHeading == planningUnitHeading
}

func actionMatchesColumnNamed(action ManagementActionType, csvHeading string) bool {
	return string(action) == csvHeading
}
