package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func (m *Mux) v1solutionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.v1GetSolutionHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetSolutionHandler(w http.ResponseWriter, r *http.Request) {
	requestSuppliedModelLabel := deriveModelLabelFrom(r)

	if !m.HasAttribute(scenarioNameKey) {
		m.Logger().Warn("Attempted to request model [" + requestSuppliedModelLabel + "] with no scenario loaded")
		m.NotFoundError(w, r)
		return
	}

	if m.solutionSetTable == nil {
		m.Logger().Warn("Attempted to request solution [" + requestSuppliedModelLabel + "] with no solution set loaded")
		m.NotFoundError(w, r)
		return
	}

	if !m.solutionSetTableContainsEntry(requestSuppliedModelLabel) {
		m.Logger().Warn("Attempted to request solution [" + requestSuppliedModelLabel + "] which is not in supplied solution set")
		m.NotFoundError(w, r)
		return
	}

	modelLabel := SolutionPoolLabel(requestSuppliedModelLabel)

	if !m.solutionPool.HasSolution(modelLabel) {
		m.Logger().Info("Loading solution [" + requestSuppliedModelLabel + "] into solution pool ")
		detail := m.getSolutionDetail(requestSuppliedModelLabel)
		m.solutionPool.AddSolution(modelLabel, detail.encoding, detail.summary)
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.solutionPool.Solution(modelLabel))

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with scenario [" + scenarioName + "] model [" + requestSuppliedModelLabel + "] state")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 models handler")
		m.Logger().Error(wrappingError)
	}
}

type solutionDetail struct {
	label    string
	encoding string
	summary  string
}

func (m *Mux) solutionSetTableContainsEntry(solutionLabel string) bool {
	const labelIndex = 0
	_, rowSize := m.solutionSetTable.ColumnAndRowSize()
	for rowIndex := uint(0); rowIndex < rowSize; rowIndex++ {
		if m.solutionSetTable.CellString(labelIndex, rowIndex) == solutionLabel {
			return true
		}
	}
	return false
}

func (m *Mux) getSolutionDetail(solutionLabel string) *solutionDetail {
	const (
		labelIndex    = 0
		encodingIndex = 6
		summaryIndex  = 7
	)
	_, rowSize := m.solutionSetTable.ColumnAndRowSize()
	for rowIndex := uint(1); rowIndex < rowSize; rowIndex++ {
		if m.solutionSetTable.CellString(labelIndex, rowIndex) == solutionLabel {
			return &solutionDetail{
				label:    solutionLabel,
				encoding: m.solutionSetTable.CellString(encodingIndex, rowIndex),
				summary:  m.solutionSetTable.CellString(summaryIndex, rowIndex),
			}
		}
	}
	return nil
}

func deriveModelLabelFrom(r *http.Request) string {
	pathElements := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	lastElementIndex := len(pathElements) - 1
	modelLabelString := pathElements[lastElementIndex]
	return modelLabelString
}
