// Copyright (c) 2018 Australian Rivers Institute.

package dataset

import (
	"errors"
)

type DataSet interface {
	Initialise(name string) DataSet
	Name() string
	Tables() map[string]Table
	Table(name string) (Table, error)
	AddTable(name string, table Table) error
}

func NewDataSet(name string) DataSet {
	return new(DataSetImpl).Initialise(name)
}

type DataSetImpl struct {
	name   string
	tables map[string]Table
}

func (dsi *DataSetImpl) Initialise(name string) DataSet {
	dsi.name = name
	dsi.tables = make(map[string]Table, 0)
	return dsi
}

func (dsi *DataSetImpl) Name() string {
	return dsi.name
}

func (dsi *DataSetImpl) Tables() map[string]Table {
	return dsi.tables
}

func (dsi *DataSetImpl) Table(name string) (Table, error) {
	table, found := dsi.tables[name]
	if found {
		return table, nil
	}
	return nil, errors.New("no table with name '" + name + "' in DataSet '" + dsi.name + "'")
}

func (dsi *DataSetImpl) AddTable(name string, table Table) error {
	if _, found := dsi.tables[name]; found {
		return errors.New("table with name '" + name + "' already in DataSet '" + dsi.name + "'")
	}

	dsi.tables[name] = table
	return nil
}
