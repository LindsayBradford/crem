package api

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	ActiveAction   = "Active"
	InactiveAction = "Inactive"
)

func (m *Mux) v1subcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.v1GetSubcatchmentHandler(w, r)
	case http.MethodPost:
		m.v1PostSubcatchmentHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetSubcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	pathElements := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	lastElementIndex := len(pathElements) - 1
	subCatchmentAsString := pathElements[lastElementIndex]

	if m.modelSolution == nil {
		m.Logger().Warn("Attempted to request subcatchment [" + subCatchmentAsString + "] state with no model present")
		m.NotFoundError(w, r)
		return
	}

	subCatchmentAsInteger, convertError := strconv.Atoi(subCatchmentAsString)
	if convertError != nil {
		panic("Should not reach here -- regular expression map should stop non-integers from being passed to handler")
	}

	subCatchment := planningunit.Id(subCatchmentAsInteger)

	var subCatchmentFound bool
	for _, action := range m.model.ManagementActions() {
		if action.PlanningUnit() == subCatchment {
			subCatchmentFound = true
		}
	}
	if !subCatchmentFound {
		m.Logger().Warn("Attempted to request subcatchment [" + subCatchmentAsString + "] state not offered by the model")
		m.NotFoundError(w, r)
		return
	}

	activeActions := m.modelSolution.ActiveManagementActions[subCatchment]
	inactiveActions := m.modelSolution.InactiveManagementActions[subCatchment]

	returnAttributes := attributes.Attributes{}

	for _, action := range activeActions {
		returnAttributes = returnAttributes.Add(string(action), ActiveAction)
	}
	for _, action := range inactiveActions {
		returnAttributes = returnAttributes.Add(string(action), InactiveAction)
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(returnAttributes)

	writeError := restResponse.Write()

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with model [" + scenarioName + "] subcatchment [" + subCatchmentAsString + "] state")

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 subcatchment handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PostSubcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.NotFoundError(w, r)
		return
	}

	pathElements := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	lastElementIndex := len(pathElements) - 1
	subCatchmentAsString := pathElements[lastElementIndex]

	subCatchmentAsInteger, convertError := strconv.Atoi(subCatchmentAsString)
	if convertError != nil {
		panic("Should not reach here -- regular expression map should stop non-integers from being passed to handler")
	}

	subCatchment := planningunit.Id(subCatchmentAsInteger)
	var planningUnitFound bool
	for _, value := range m.modelSolution.PlanningUnits {
		if value == subCatchment {
			planningUnitFound = true
		}
	}

	if !planningUnitFound {
		m.NotFoundError(w, r)
		return
	}

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Processing POST of model [" + scenarioName + "] subcatchment [" + subCatchmentAsString + "] state")

	requestContent := requestBodyToBytes(r)
	postedAttributes := attributes.Attributes{}

	unmnarshalError := json.Unmarshal(requestContent, &postedAttributes)
	if unmnarshalError != nil {
		wrappingError := errors.Wrap(unmnarshalError, "v1 POST subcatchment handler")
		m.Logger().Error(wrappingError)
		m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
	}

	for _, entry := range postedAttributes {
		if entry.Name != "RiverBankRestoration" && entry.Name != "HillSlopeRestoration" && entry.Name != "GullyRestoration" {
			baseError := errors.New("Name [" + entry.Name + "] not one of [RiverBankRestoration, HillSlopeRestoration, GullyRestoration]")
			wrappingError := errors.Wrap(baseError, "v1 POST subcatchment handler")
			m.Logger().Error(wrappingError)
			m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
			return
		}

		if entry.Value != InactiveAction && entry.Value != ActiveAction {
			baseError := errors.New("For named action [" + entry.Name + "], value [" + entry.Value.(string) + "] not one of [Active,Inactive]")
			wrappingError := errors.Wrap(baseError, "v1 POST subcatchment handler")
			m.Logger().Error(wrappingError)
			m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
			return
		}
	}

	for actionIndex, action := range m.model.ManagementActions() {
		if subCatchment != action.PlanningUnit() {
			continue
		}

		for _, entry := range postedAttributes {
			if entry.Name == string(action.Type()) {
				if entry.Value == InactiveAction {
					m.model.SetManagementAction(actionIndex, false)
					m.model.AcceptAll()
					m.Logger().Info("Model subcatchment [" + subCatchmentAsString + "], Action [" + entry.Name + "], set to [" + entry.Value.(string) + "]")
				}
				if entry.Value == ActiveAction {
					m.model.SetManagementAction(actionIndex, true)
					m.model.AcceptAll()
					m.Logger().Info("Model subcatchment [" + subCatchmentAsString + "], Action [" + entry.Name + "], set to [" + entry.Value.(string) + "]")
				}
			}
		}
	}

	m.modelSolution = new(solution.SolutionBuilder).
		WithId(m.model.Id()).
		ForModel(m.model).
		Build()

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.modelSolution)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 POST subcatchment handler")
		m.Logger().Error(wrappingError)
	}
}
