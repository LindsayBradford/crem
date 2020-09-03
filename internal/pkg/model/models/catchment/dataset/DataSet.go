package dataset

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

const (
	SubcatchmentsTableName = "Subcatchments"
	GulliesTableName       = "Gullies"
	ActionsTableName       = "Actions"
)

type DataSet interface {
	dataset.DataSet
	SubCatchmentsTable()
	ActionsTable()
	GulliesTable()
}

type DataSetImpl struct {
	dataset.DataSet

	SubCatchmentsTable tables.CsvTable
	ActionsTable       tables.CsvTable
	GulliesTable       tables.CsvTable
}

func (c *DataSetImpl) Initialise(wrappedDataSet dataset.DataSet) *DataSetImpl {
	c.DataSet = wrappedDataSet

	c.SubCatchmentsTable = tables.ToCsvTable(c.DataSet, SubcatchmentsTableName)
	c.ActionsTable = tables.ToCsvTable(c.DataSet, ActionsTableName)
	c.GulliesTable = tables.ToCsvTable(c.DataSet, GulliesTableName)

	return c
}
