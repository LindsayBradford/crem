package api

import (
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
	m.RespondWithError(http.StatusNotFound, "Behaviour not yet implemented", w, r)
}
