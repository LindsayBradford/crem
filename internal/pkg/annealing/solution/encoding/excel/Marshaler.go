// Copyright (c) 2019 Australian Rivers Institute.

package excel

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

const (
	DecisionVariablesTableName = "NameMappedVariables"
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

var baseVariableHeadings = []string{nameHeading, valueHeading, unitOfMeasureHeading}

const (
	planningUnitColumn = 0

	inactiveActionValue = 0
	activeActionValue   = 1
)

type Marshaler struct{}

func (m *Marshaler) Marshal(solution *solution.Solution, dataSet *excel.DataSet) error {
	if variableErr := m.marshalDecisionVariables(solution, dataSet); variableErr != nil {
		return variableErr
	}

	if variableErr := m.marshalActionState(solution, dataSet); variableErr != nil {
		return variableErr
	}

	return nil
}

func (m *Marshaler) marshalDecisionVariables(solution *solution.Solution, dataSet *excel.DataSet) error {
	table := emptyDecisionVariableTable(solution)

	var offsetColumn uint = unitOfMeasureColumn + 1

	for i, decisionVariable := range solution.DecisionVariables {
		rowIndex := uint(i)
		table.SetCell(nameColumn, rowIndex, decisionVariable.Name)
		table.SetCell(valueColumn, rowIndex, decisionVariable.Value)
		table.SetCell(unitOfMeasureColumn, rowIndex, decisionVariable.Measure.String())

		if decisionVariable.ValuePerPlanningUnit != nil {
			for planningUnitIndex, planningUnit := range solution.PlanningUnits {
				columnIndex := uint(planningUnitIndex) + offsetColumn
				table.SetCell(columnIndex, rowIndex, 0)

				for _, variableValue := range decisionVariable.ValuePerPlanningUnit {
					if planningUnit == variableValue.PlanningUnit {
						table.SetCell(columnIndex, rowIndex, variableValue.Value)
					}
				}
			}
		}
	}

	dataSet.AddTable(table.Name(), table)
	return nil
}

func emptyDecisionVariableTable(solution *solution.Solution) *tables.CsvTableImpl {
	table := new(tables.CsvTableImpl)

	headings := variableHeadings(solution)

	table.SetHeader(headings)
	table.SetName(DecisionVariablesTableName)
	table.SetColumnAndRowSize(
		uint(len(headings)),
		uint(len(solution.DecisionVariables)),
	)

	return table
}

func variableHeadings(solution *solution.Solution) []string {
	finalisedHeadings := make([]string, len(baseVariableHeadings)+len(solution.PlanningUnits))

	for index, entry := range baseVariableHeadings {
		finalisedHeadings[index] = entry
	}

	baseOffset := len(baseVariableHeadings)

	for index, entry := range solution.PlanningUnits {
		finalisedHeadings[index+baseOffset] = solution.PlanningUnitHeading() + "-" + entry.String()
	}

	return finalisedHeadings
}

func (m *Marshaler) marshalActionState(solution *solution.Solution, dataSet *excel.DataSet) error {
	table, actionHeadings := emptyActionTable(solution)

	for y, planningUnit := range solution.PlanningUnits {
		rowIndex := uint(y)

		table.SetCell(planningUnitColumn, rowIndex, uint64(planningUnit))

		for x, csvHeading := range actionHeadings {
			columnIndex := uint(x)
			if shouldSkipColumnWith(solution, csvHeading) {
				continue
			}
			table.SetCell(columnIndex, rowIndex, inactiveActionValue)
		}

		if activeActions, unitHasActiveActions := solution.ActiveManagementActions[planningUnit]; unitHasActiveActions {
			for x, csvHeading := range actionHeadings {
				columnIndex := uint(x)
				if shouldSkipColumnWith(solution, csvHeading) {
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

	headings[planningUnitColumn] = solution.PlanningUnitHeading()
	headings = append(headings, solution.ActionsAsStrings()...)

	return headings
}

func shouldSkipColumnWith(solution *solution.Solution, csvHeading string) bool {
	return csvHeading == solution.PlanningUnitHeading()
}

func actionMatchesColumnNamed(action solution.ManagementActionType, csvHeading string) bool {
	return string(action) == csvHeading
}
