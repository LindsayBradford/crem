package api

import (
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

	// requestContent := requestBodyToString(r)
	// m.writePostContentAsResponse(requestContent,w)

	m.processRequestContentForActiveActions(r, w)
	m.writeActiveActionResponse(w)
}

func (m *Mux) processRequestContentForActiveActions(r *http.Request, w http.ResponseWriter) {
	tmpDataSet := csv.NewDataSet("Content Dataset")
	defer tmpDataSet.Teardown()

	tmpDataSet.ParseCsvTextIntoTable("requestContent", requestBodyToString(r))
	if tmpDataSet.Errors() != nil {
		wrappingError := errors.Wrap(tmpDataSet.Errors(), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		return
	}

	contentTable, tableError := tmpDataSet.Table("requestContent")
	if tableError != nil {
		wrappingError := errors.Wrap(tmpDataSet.Errors(), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		return
	}

	if contentTable == nil {
		wrappingError := errors.Wrap(errors.New("No CSV table content found"), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		m.BadRequestError(w, r)
		return
	}

	headingsTable, hasHeadings := contentTable.(dataset.HeadingsTable)
	if !hasHeadings {
		wrappingError := errors.Wrap(errors.New("CSV table does not have a header row"), "v1 model actions handler")
		m.Logger().Error(wrappingError)
		m.BadRequestError(w, r)
		return
	}

	if headingsTable.Header()[0] != "SubCatchment" {
		wrappingError := errors.Wrap(
			errors.New("CSV table header column misses mandatory 'SubCatchment' first entry"),
			"v1 model actions handler")
		m.Logger().Error(wrappingError)
		m.BadRequestError(w, r)
		return
	}

	modelActions := m.model.ManagementActions()

	var suppliedActionState bool
	colSize, rowSize := headingsTable.ColumnAndRowSize()
	for rowIndex := uint(0); rowIndex < rowSize; rowIndex++ {
		for colIndex := uint(1); colIndex < colSize; colIndex++ {
			if headingsTable.CellFloat64(colIndex, rowIndex) == 0 {
				suppliedActionState = false
			} else {
				suppliedActionState = true
			}

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
	}

	m.updateModelSolution()
}

// TODO: deprecate once CSV content POST processing is in place.
func (m *Mux) writePostContentAsResponse(postContent string, w http.ResponseWriter) {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithCsvContent(postContent)

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with model [" + scenarioName + "] active actions state")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model actions handler")
		m.Logger().Error(wrappingError)
	}
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
	wrappingError := errors.Wrap(contentTypeError, "v1 POST model actions handler")
	m.Logger().Warn(wrappingError)

	m.MethodNotAllowedError(w, r)
}
