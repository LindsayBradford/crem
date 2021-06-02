// Copyright (c) 2019 Australian Rivers Institute.

package excel

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

const (
	idHeading      = "Solution"
	actionsHeading = "Actions"
	summaryHeading = "Summary"

	SummaryTableName = "Summary"
)

type Marshaler struct{}

func (m *Marshaler) Marshal(summary *set.Summary, dataSet *excel.DataSet) error {
	table := emptySummaryTable(summary)

	rowIndex := uint(0)
	for _, value := range summary.AsSortedArray() {
		columnIndex := uint(0)
		table.SetCell(columnIndex, rowIndex, value.Id)

		for index, variable := range value.Variables {
			columnOffset := columnIndex + uint(index+1)
			table.SetCell(columnOffset, rowIndex, variable.Value)
		}

		columnOffset := columnIndex + uint(len(value.Variables)+1)
		table.SetCell(columnOffset, rowIndex, string(value.Actions))
		table.SetCell(columnOffset+1, rowIndex, value.Note)

		rowIndex++
	}

	dataSet.AddTable(table.Name(), table)
	return nil
}

func emptySummaryTable(summary *set.Summary) *tables.CsvTableImpl {
	table := new(tables.CsvTableImpl)

	headings := deriveHeaders(summary)

	table.SetHeader(headings)
	table.SetName(SummaryTableName)
	table.SetColumnAndRowSize(
		uint(len(headings)),
		uint(len(*summary)),
	)

	return table
}

func deriveHeaders(summary *set.Summary) []string {
	exampleVariables := justSomeVariables(summary)

	headingNumber := len(exampleVariables) + 3
	headers := make([]string, headingNumber)

	headers[0] = idHeading
	for index, variable := range exampleVariables {
		headers[index+1] = variable.Name
	}
	headers[headingNumber-2] = actionsHeading
	headers[headingNumber-1] = summaryHeading

	return headers
}

func justSomeVariables(summary *set.Summary) solution.VariableSetSummary {
	for _, currentSummary := range *summary {
		return currentSummary.Variables
	}
	return nil
}
