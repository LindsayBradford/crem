package api

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/csv"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/pkg/errors"
	"net/http"
)

func (m *Mux) v1actionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.v1PostActionsHandler(w, r)
	case http.MethodGet:
		m.v1GetActionsHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

type actionsWrapper struct {
	ActiveManagementActions map[planningunit.Id]solution.ManagementActions
}

func (m *Mux) v1GetActionsHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.NotFoundError(w, r)
		return
	}

	m.writeActiveActionResponse(w)
}

func (m *Mux) writeActiveActionResponse(w http.ResponseWriter) {
	activeActions := actionsWrapper{
		ActiveManagementActions: m.modelSolution.ActiveManagementActions,
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(activeActions)

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with model [" + scenarioName + "] active actions state")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model actions handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PostActionsHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.NotFoundError(w, r)
		return
	}

	if m.requestContentTypeWasNotCsv(r, w) {
		return
	}

	processError := m.processRequestContentForActiveActions(r, w)
	if processError != nil {
		return
	}

	restResponse := m.buildPostActionsResponse(w)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model actions handler")
		m.Logger().Error(wrappingError)
	}

}

func (m *Mux) buildPostActionsResponse(w http.ResponseWriter) *rest.Response {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(
			rest.MessageResponse{
				Type:    "SUCCESS",
				Message: "Management actions state change successfully applied",
				Time:    rest.FormattedTimestamp(),
			},
		)

	m.Logger().Info("Responding with acknowledgement of management actions state change ")

	return restResponse
}

func (m *Mux) processRequestContentForActiveActions(r *http.Request, w http.ResponseWriter) error {
	requestTable, requestError := m.deriveRequestTable(r, w)
	if requestError != nil {
		return requestError
	}

	m.processRequestTable(requestTable)
	m.updateModelSolution()

	return nil
}

func (m *Mux) processRequestTable(headingsTable dataset.HeadingsTable) {
	colSize, rowSize := headingsTable.ColumnAndRowSize()
	for rowIndex := uint(0); rowIndex < rowSize; rowIndex++ {
		for colIndex := uint(1); colIndex < colSize; colIndex++ {
			m.processTableCell(headingsTable, colIndex, rowIndex)
		}
	}
}

func (m *Mux) processTableCell(headingsTable dataset.HeadingsTable, colIndex uint, rowIndex uint) {
	suppliedActionState := m.deriveSuppliedActionState(headingsTable, colIndex, rowIndex)
	modelActions := m.model.ManagementActions()

	for actionIndex := 0; actionIndex < len(modelActions); actionIndex++ {
		currentAction := modelActions[actionIndex]

		rawPlanningUnit := headingsTable.CellFloat64(0, rowIndex)
		rawType := headingsTable.Header()[colIndex]

		if currentAction.PlanningUnit() == planningunit.Id(rawPlanningUnit) &&
			string(currentAction.Type()) == rawType {
			m.model.SetManagementAction(actionIndex, suppliedActionState)
		}
	}
}

func (m *Mux) deriveSuppliedActionState(headingsTable dataset.HeadingsTable, colIndex uint, rowIndex uint) bool {
	var suppliedActionState bool

	if headingsTable.CellFloat64(colIndex, rowIndex) == 0 {
		suppliedActionState = false
	} else {
		suppliedActionState = true
	}
	return suppliedActionState
}

func (m *Mux) deriveRequestTable(r *http.Request, w http.ResponseWriter) (dataset.HeadingsTable, error) {
	tmpDataSet := csv.NewDataSet("Content Dataset")
	defer tmpDataSet.Teardown()

	tmpDataSet.ParseCsvTextIntoTable("requestContent", requestBodyToString(r))
	if tmpDataSet.Errors() != nil {
		wrappingError := errors.Wrap(tmpDataSet.Errors(), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	contentTable, tableError := tmpDataSet.Table("requestContent")
	if tableError != nil {
		wrappingError := errors.Wrap(tmpDataSet.Errors(), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	if contentTable == nil {
		wrappingError := errors.Wrap(errors.New("No CSV table content found"), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		m.BadRequestError(w, r)
		return nil, wrappingError
	}

	headingsTable, hasHeadings := contentTable.(dataset.HeadingsTable)
	if !hasHeadings {
		wrappingError := errors.Wrap(errors.New("CSV table does not have a header row"), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		m.BadRequestError(w, r)
		return nil, wrappingError
	}

	if headingsTable.Header()[0] != "SubCatchment" {
		wrappingError := errors.Wrap(
			errors.New("CSV table header column misses mandatory 'SubCatchment' first entry"),
			"v1 model actions handler")
		m.Logger().Error(wrappingError)
		m.BadRequestError(w, r)
		return nil, wrappingError
	}

	colSize, rowSize := headingsTable.ColumnAndRowSize()
	for rowIndex := uint(0); rowIndex < rowSize; rowIndex++ {
		for colIndex := uint(1); colIndex < colSize; colIndex++ {
			cellValue := headingsTable.CellFloat64(colIndex, rowIndex)
			if cellValue != 0 && cellValue != 1 {
				msgText := fmt.Sprintf(
					"Table management action cell [%d,%d] has invalid value [%v]. Must be one of [0,1]",
					colIndex, rowIndex, cellValue)
				wrappingError := errors.Wrap(errors.New(msgText), "v1 model actions handler")
				m.Logger().Error(wrappingError)
				m.RespondWithError(http.StatusBadRequest, msgText, w, r)
				return nil, wrappingError
			}
		}
	}

	return headingsTable, nil
}

func (m *Mux) requestContentTypeWasNotCsv(r *http.Request, w http.ResponseWriter) bool {
	suppliedContentType := r.Header.Get(rest.ContentTypeHeaderKey)
	if suppliedContentType != rest.CsvMimeType {
		m.handleNonCsvContentResponse(r, w, suppliedContentType)
		return true
	}
	return false
}

func (m *Mux) handleNonCsvContentResponse(r *http.Request, w http.ResponseWriter, suppliedContentType string) {
	contentTypeError := errors.New("Request content-type of [" + suppliedContentType + "] was not the expected [" + rest.CsvMimeType + "]")
	wrappingError := errors.Wrap(contentTypeError, "v1 model actions handler")
	m.Logger().Warn(wrappingError)

	m.MethodNotAllowedError(w, r)
}
