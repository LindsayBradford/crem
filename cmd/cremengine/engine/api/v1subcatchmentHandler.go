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

func (m *Mux) v1subcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.v1PostSubcatchmentHandler(w, r)
	case http.MethodGet:
		m.v1GetSubcatchmentHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetSubcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.NotFoundError(w, r)
		return
	}

	pathElements := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	lastElementIndex := len(pathElements) - 1
	subCatchment := pathElements[lastElementIndex]

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with model [" + scenarioName + "] subcatchment [" + subCatchment + "] state")

	subCatchmentAsInteger, convertError := strconv.Atoi(subCatchment)
	if convertError != nil {
		m.NotFoundError(w, r)
		return
	}

	activeActions := m.modelSolution.ActiveManagementActions[planningunit.Id(subCatchmentAsInteger)]
	inactiveActions := m.modelSolution.InactiveManagementActions[planningunit.Id(subCatchmentAsInteger)]

	returnAttributes := attributes.Attributes{}

	for _, action := range activeActions {
		returnAttributes = returnAttributes.Add(string(action), "Active")
	}
	for _, action := range inactiveActions {
		returnAttributes = returnAttributes.Add(string(action), "Inactive")
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(returnAttributes)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model handler")
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
	subCatchment := pathElements[lastElementIndex]

	subCatchmentAsInteger, convertError := strconv.Atoi(subCatchment)
	if convertError != nil {
		m.NotFoundError(w, r)
		return
	}

	var planningUnitFound bool
	for _, value := range m.modelSolution.PlanningUnits {
		if value == planningunit.Id(subCatchmentAsInteger) {
			planningUnitFound = true
		}
	}

	if !planningUnitFound {
		m.NotFoundError(w, r)
		return
	}

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Processing POST of model [" + scenarioName + "] subcatchment [" + subCatchment + "] state")

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

		if entry.Value != "Inactive" && entry.Value != "Active" {
			baseError := errors.New("For named action [" + entry.Name + "], value [" + entry.Value.(string) + "] not one of [Active,Inactive]")
			wrappingError := errors.Wrap(baseError, "v1 POST subcatchment handler")
			m.Logger().Error(wrappingError)
			m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
			return
		}
	}

	for actionIndex, action := range m.model.ManagementActions() {
		if planningunit.Id(subCatchmentAsInteger) != action.PlanningUnit() {
			continue
		}

		for _, entry := range postedAttributes {
			if entry.Name == string(action.Type()) {
				if entry.Value == "Inactive" {
					m.model.SetManagementAction(actionIndex, false)
					m.model.AcceptAll()
					m.Logger().Info("Model subcatchment [" + subCatchment + "], Action [" + entry.Name + "], set to [" + entry.Value.(string) + "]")
				}
				if entry.Value == "Active" {
					m.model.SetManagementAction(actionIndex, true)
					m.model.AcceptAll()
					m.Logger().Info("Model subcatchment [" + subCatchment + "], Action [" + entry.Name + "], set to [" + entry.Value.(string) + "]")
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
