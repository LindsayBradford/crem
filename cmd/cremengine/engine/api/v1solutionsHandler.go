package api

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/csv"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
	"net/http"
)

func (m *Mux) v1solutionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.v1PostSolutionsHandler(w, r)
	case http.MethodGet:
		m.v1GetSolutionsHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetSolutionsHandler(w http.ResponseWriter, r *http.Request) {
	if !m.HasAttribute(scenarioTextKey) {
		m.Logger().Warn("Request for solutions dataset received without pre-requisite scenario loaded.")
		m.NotFoundError(w, r)
		return
	}

	if !m.HasAttribute(solutionsTextKey) {
		m.Logger().Warn("Request for solutions dataset received before dataset had been posted.")
		m.NotFoundError(w, r)
		return
	}

	restResponse := m.buildSolutionsGetResponse(w)
	m.logScenarioGetResponse()
	writeError := restResponse.Write()

	m.handleScenarioGetWriteError(writeError)
}

func (m *Mux) buildSolutionsGetResponse(w http.ResponseWriter) *rest.Response {
	responseText := m.Attribute(solutionsTextKey).(string)

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithTomlContent(responseText)
	return restResponse
}

func (m *Mux) logSolutionsGetResponse() {
	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with scenario [" + scenarioName + "] solutions table")
}

func (m *Mux) handleSolutionsGetWriteError(writeError error) {
	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 solutions handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PostSolutionsHandler(w http.ResponseWriter, r *http.Request) {
	if !m.HasAttribute(scenarioTextKey) {
		m.Logger().Warn("Request to POST scenario solutions dataset without scenario loaded.")
		m.MethodNotAllowedError(w, r)
		return
	}

	if m.requestContentTypeWasNotCsv(r, w) {
		return
	}

	processError := m.processRequestContentForSolutions(r, w)
	if processError != nil {
		m.Logger().Warn("Request to POST scenario solutions dataset with invalid solution data detected.")
		m.BadRequestError(w, r)
		return
	}

	m.processSolutionsPostText(r)

	restResponse := m.buildSolutionsPostResponse(w)
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 solutions handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) processRequestContentForSolutions(r *http.Request, w http.ResponseWriter) error {
	rawTableContent := requestBodyToString(r)

	var requestError error
	m.solutionsTable, requestError = m.deriveSolutionsRequestTable(rawTableContent)
	if requestError != nil {
		return requestError
	}

	return nil
}

func (m *Mux) deriveSolutionsRequestTable(rawTableContent string) (dataset.HeadingsTable, error) {
	tmpDataSet := csv.NewDataSet("Content Dataset")
	defer tmpDataSet.Teardown()

	tmpDataSet.ParseCsvTextIntoTable("requestContent", rawTableContent)
	if tmpDataSet.Errors() != nil {
		wrappingError := errors.Wrap(tmpDataSet.Errors(), "v1 solutions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	contentTable, tableError := tmpDataSet.Table("requestContent")
	if tableError != nil {
		wrappingError := errors.Wrap(tmpDataSet.Errors(), "v1 solutions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	if contentTable == nil {
		wrappingError := errors.Wrap(errors.New("No CSV table content found"), "v1 solutions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	contentTableWithHeadings, hasHeadings := contentTable.(dataset.HeadingsTable)
	if !hasHeadings {
		wrappingError := errors.Wrap(errors.New("CSV table does not have a header row"), "v1 solutions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	if contentTableWithHeadings.Header()[0] != "Solution" {
		wrappingError := errors.Wrap(
			errors.New("CSV table header column misses mandatory 'Solution' first entry"),
			"v1 solutions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	updateErrors := compositeErrors.New("v1 POST solutions handler")

	colSize, rowSize := contentTableWithHeadings.ColumnAndRowSize()
	for rowIndex := uint(0); rowIndex < rowSize; rowIndex++ {
		for colIndex := uint(1); colIndex < colSize; colIndex++ {
			if rowIndex > 0 {
				cellValue := contentTableWithHeadings.Cell(colIndex, rowIndex)
				heading := contentTableWithHeadings.Header()[colIndex]
				switch heading {
				case "Solution", "Actions", "Summary":
					switch cellValue.(type) {
					case string:
						break // deliberately do nothing
					default:
						msgText := fmt.Sprintf(
							"Table management action cell [%d,%d] has invalid type. Must be a string",
							colIndex, rowIndex)
						updateErrors.AddMessage(msgText)
						m.Logger().Error(msgText)
					}
				default:
					switch cellValue.(type) {
					case float64:
						break // deliberately does nothing
					default:
						msgText := fmt.Sprintf(
							"Table management action cell [%d,%d] has invalid type. Must be a 64-bit floating point decimal",
							colIndex, rowIndex)
						updateErrors.AddMessage(msgText)
						m.Logger().Error(msgText)
					}
				}
			}
		}
	}

	if updateErrors.Size() > 0 {
		return nil, updateErrors
	}

	return contentTableWithHeadings, nil
}

func (m *Mux) buildSolutionsPostResponse(w http.ResponseWriter) *rest.Response {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(
			rest.MessageResponse{
				Type:    "SUCCESS",
				Message: "Scenario solutions set successfully posted",
				Time:    rest.FormattedTimestamp(),
			},
		)

	m.logSolutionsGetResponse()
	return restResponse
}

func (m *Mux) processSolutionsPostText(r *http.Request) error {
	requestContent := requestBodyToString(r)
	m.rememberSolutionsAttributeState(requestContent)
	return nil
}

func (m *Mux) rememberSolutionsAttributeState(requestContent string) {
	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Scenario [" + scenarioName + "] solutions dataset successfully retrieved")
	m.ReplaceAttribute(solutionsTextKey, requestContent)
}
