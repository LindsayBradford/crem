// Copyright (c) 2019 Australian Rivers Institute.

package excel

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

const (
	DecisionVariablesTableName = "DecisionVariables"
	ManagementActionsTableName = "ManagementActions"
)

const (
	nameHeading = "Name"
	nameColumn  = 0

	valueHeading = "Value"
	valueColumn  = 1

	unitOfMeasureHeading = "UnitOfMeasure"
	unitOfMeasureColumn  = 2
)

var variableHeadings = []string{nameHeading, valueHeading, unitOfMeasureHeading}

const (
	planningUnitHeading = "PlanningUnit"
	planningUnitColumn  = 0

	inactiveActionValue = 0
	activeActionValue   = 1
)

type Marshaler struct{}

func (m *Marshaler) Marshal(solution *solution.Solution, dataSet *excel.DataSet) error {
	if variableErr := m.marshalDecisionVariables(solution, dataSet); variableErr != nil {
		return variableErr
	}

	if variableErr := m.marshalActiveActions(solution, dataSet); variableErr != nil {
		return variableErr
	}

	return nil
}

func (m *Marshaler) marshalDecisionVariables(solution *solution.Solution, dataSet *excel.DataSet) error {
	table := emptyDecisionVariableTable(solution)

	for i, decisionVariable := range solution.DecisionVariables {
		rowIndex := uint(i)
		table.SetCell(nameColumn, rowIndex, decisionVariable.Name)
		table.SetCell(valueColumn, rowIndex, decisionVariable.Value)
		table.SetCell(unitOfMeasureColumn, rowIndex, decisionVariable.Measure.String())
	}

	dataSet.AddTable(table.Name(), table)
	return nil
}

func emptyDecisionVariableTable(solution *solution.Solution) *tables.CsvTableImpl {
	table := new(tables.CsvTableImpl)

	table.SetHeader(variableHeadings)
	table.SetName(DecisionVariablesTableName)
	table.SetColumnAndRowSize(
		uint(len(variableHeadings)),
		uint(len(solution.DecisionVariables)),
	)

	return table
}

func (m *Marshaler) marshalActiveActions(solution *solution.Solution, dataSet *excel.DataSet) error {
	table, actionHeadings := emptyActionTable(solution)

	for y, planningUnit := range solution.PlanningUnits {
		rowIndex := uint(y)
		table.SetCell(planningUnitColumn, rowIndex, string(planningUnit))

		if activeActions, unitHasActiveActions := solution.ActiveManagementActions[planningUnit]; unitHasActiveActions {
			for x, csvHeading := range actionHeadings {
				columnIndex := uint(x)
				if shouldSkipColumnWith(csvHeading) {
					continue
				}

				actionValue := inactiveActionValue
				for _, action := range activeActions {
					if actionMatchesColumnNamed(action, csvHeading) {
						actionValue = activeActionValue
					}
				}

				table.SetCell(columnIndex, rowIndex, actionValue)
			}
		} else {
			for x, csvHeading := range actionHeadings {
				columnIndex := uint(x)
				if shouldSkipColumnWith(csvHeading) {
					continue
				}
				table.SetCell(columnIndex, rowIndex, inactiveActionValue)
			}
		}
	}

	dataSet.AddTable(table.Name(), table)
	return nil
}

func emptyActionTable(solution *solution.Solution) (table *tables.CsvTableImpl, actionHeadings []string) {
	table = new(tables.CsvTableImpl)

	actionHeadings = tableHeadings(solution)
	table.SetHeader(actionHeadings)
	table.SetName(ManagementActionsTableName)
	table.SetColumnAndRowSize(uint(len(actionHeadings)), uint(len(solution.PlanningUnits)))

	return table, actionHeadings
}

func tableHeadings(solution *solution.Solution) []string {
	headings := make([]string, 1)

	headings[planningUnitColumn] = planningUnitHeading

	headings = append(headings, solution.ActiveActionsAsStrings()...)

	return headings
}

func shouldSkipColumnWith(csvHeading string) bool {
	return csvHeading == planningUnitHeading
}

func actionMatchesColumnNamed(action solution.ManagementActionType, csvHeading string) bool {
	return string(action) == csvHeading
}
