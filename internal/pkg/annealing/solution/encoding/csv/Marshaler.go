// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"fmt"
	strings2 "strings"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
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

type DecisionVariableMarshaler struct{}

func (cm *DecisionVariableMarshaler) Marshal(solution *solution.Solution) ([]byte, error) {
	return cm.marshalDecisionVariables(solution)
}

func (cm *DecisionVariableMarshaler) marshalDecisionVariables(solution *solution.Solution) ([]byte, error) {
	csvStringAsBytes := ([]byte)(cm.decisionVariablesToCsvString(solution))
	return csvStringAsBytes, nil
}

func (cm *DecisionVariableMarshaler) decisionVariablesToCsvString(solution *solution.Solution) string {

	planningUnits := planningUnitsAsHeaders(solution.PlanningUnits)

	builder := new(strings.FluentBuilder)
	builder.
		Add(join(variableHeadings...)).
		Add(separator).
		Add(join(planningUnits...)).
		Add(newline)

	variables := solution.DecisionVariables

	for _, variable := range variables {
		joinedVariableAttributes := joinAttributes(variable, solution.PlanningUnits)
		builder.Add(joinedVariableAttributes).Add(newline)
	}

	return builder.String()
}

func planningUnitsAsHeaders(planningUnits solution.PlanningUnitIds) []string {
	headers := make([]string, len(planningUnits))

	for index, value := range planningUnits {
		headers[index] = planningUnitHeading + "-" + value
	}

	return headers
}

func joinAttributes(variable variable.EncodeableDecisionVariable, planningUnits solution.PlanningUnitIds) string {
	planningUnitValues := planningUnitValueList(variable, planningUnits)

	baseVariableValues := join(
		variable.Name,
		toString(variable.Value),
		variable.Measure.String(),
	)

	joinedPlanningUnitValues := join(planningUnitValues...)

	return join(baseVariableValues, joinedPlanningUnitValues)
}

func planningUnitValueList(variable variable.EncodeableDecisionVariable, planningUnits solution.PlanningUnitIds) []string {
	headers := make([]string, len(planningUnits))

	if variable.ValuePerPlanningUnit != nil {
		for index := range planningUnits {
			headers[index] = toString(variable.ValuePerPlanningUnit[index].Value)
		}
	}
	return headers
}

func join(entries ...string) string {
	return strings2.Join(entries, separator)
}

func toString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

type ManagementActionMarshaler struct{}

func (cm *ManagementActionMarshaler) Marshal(solution *solution.Solution) ([]byte, error) {
	return cm.marshalManagementActions(solution)
}

func (cm *ManagementActionMarshaler) marshalManagementActions(solution *solution.Solution) ([]byte, error) {
	csvStringAsBytes := ([]byte)(cm.csvEncodeManagementActions(solution))
	return csvStringAsBytes, nil
}

func (cm *ManagementActionMarshaler) csvEncodeManagementActions(solution *solution.Solution) string {
	builder := new(strings.FluentBuilder)

	headings := csvEncodeActionHeadings(solution)
	builder.Add(join(headings...)).Add(newline)

	for _, planningUnit := range solution.PlanningUnits {
		actionRowValues := cm.buildActionCsvValuesForPlanningUnit(headings, planningUnit, solution)
		builder.Add(join(actionRowValues...)).Add(newline)
	}

	return builder.String()
}

func csvEncodeActionHeadings(solution *solution.Solution) []string {
	headings := make([]string, 1)
	headings[0] = planningUnitHeading

	headings = append(headings, solution.ActiveActionsAsStrings()...)
	headings = append(headings, solution.InactiveActionsAsStrings()...)

	return headings
}

func (cm *ManagementActionMarshaler) buildActionCsvValuesForPlanningUnit(
	actionHeadings []string, planningUnit string, solution *solution.Solution) []string {

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

func actionMatchesColumnNamed(action solution.ManagementActionType, csvHeading string) bool {
	return string(action) == csvHeading
}
