package api

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/csv"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
	"net/http"
)

func (m *Mux) v1actionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		m.v1PutActionsHandler(w, r)
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

func (m *Mux) v1PutActionsHandler(w http.ResponseWriter, r *http.Request) {
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
	m.checkModelInSolutionSummary()
}

func (m *Mux) checkModelInSolutionSummary() {
	compressedModel := modelCompressor.Compress(m.model)
	encodingOfModel := compressedModel.Encoding()
	m.model.ReplaceAttribute("Encoding", encodingOfModel)
	m.checkEncodingInSolutionSummary(encodingOfModel)
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
	rawTableContent := requestBodyToString(r)

	requestTable, parseError := m.deriveSolutionTable(rawTableContent)
	if parseError != nil {
		m.BadRequestError(w, r)
	}
	return requestTable, parseError
}

func (m *Mux) deriveSolutionTable(rawTableContent string) (dataset.HeadingsTable, error) {
	tmpDataSet := csv.NewDataSet("Content Dataset")
	defer tmpDataSet.Teardown()

	tmpDataSet.ParseCsvTextIntoTable("requestContent", rawTableContent)
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
		return nil, wrappingError
	}

	headingsTable, hasHeadings := contentTable.(dataset.HeadingsTable)
	if !hasHeadings {
		wrappingError := errors.Wrap(errors.New("CSV table does not have a header row"), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	if headingsTable.Header()[0] != "SubCatchment" {
		wrappingError := errors.Wrap(
			errors.New("CSV table header column misses mandatory 'SubCatchment' first entry"),
			"v1 model actions handler")
		m.Logger().Error(wrappingError)
		return nil, wrappingError
	}

	updateErrors := compositeErrors.New("v1 POST actions handler")

	colSize, rowSize := headingsTable.ColumnAndRowSize()
	for rowIndex := uint(0); rowIndex < rowSize; rowIndex++ {
		for colIndex := uint(1); colIndex < colSize; colIndex++ {

			cellValue := headingsTable.Cell(colIndex, rowIndex)
			switch cellValue.(type) {
			case float64:
				if cellValue != float64(0) && cellValue != float64(1) {
					msgText := fmt.Sprintf(
						"Table management action cell [%d,%d] has invalid value [%v]. Must be one of [0,1]",
						colIndex, rowIndex, cellValue)
					updateErrors.AddMessage(msgText)
					m.Logger().Error(msgText)
				}
			default:
				msgText := fmt.Sprintf(
					"Table management action cell [%d,%d] has invalid value [%v]. Must be one of [0,1]",
					colIndex, rowIndex, cellValue)
				updateErrors.AddMessage(msgText)
				m.Logger().Error(msgText)
			}
		}
	}

	if updateErrors.Size() > 0 {
		return nil, updateErrors
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

	m.UnsupportedMediaTypeError(w, r)
}

func (m *Mux) SetSolution(solutionFilePath string) {
	rawTableContent := readFileAsText(solutionFilePath)

	requestTable, parseError := m.deriveSolutionTable(rawTableContent)
	if parseError != nil {
		wrappingError := errors.Wrap(parseError, "v1 model actions handler")
		m.Logger().Error(wrappingError)
		return
	}

	m.processRequestTable(requestTable)
	m.updateModelSolution()
}
