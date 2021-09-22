package api

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	v1ApplicableActionsHandler = "v1 subcatchment applicable actions handler"
)

func (m *Mux) v1ApplicableActionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.v1GetApplicableActionsHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetApplicableActionsHandler(w http.ResponseWriter, r *http.Request) {
	deriveRequestSuppliedSubCatchment := func(r *http.Request) string {
		pathElements := strings.Split(r.URL.Path, rest.UrlPathSeparator)
		lastElementIndex := len(pathElements) - 2
		subCatchmentAsString := pathElements[lastElementIndex]
		return subCatchmentAsString
	}

	requestSuppliedSubCatchment := deriveRequestSuppliedSubCatchment(r)

	if m.modelSolution == nil {
		m.Logger().Warn("Attempted to request subcatchment [" + requestSuppliedSubCatchment + "] state with no model present")
		m.NotFoundError(w, r)
		return
	}

	subCatchment := toPlanningUnitId(requestSuppliedSubCatchment)

	if !m.modelContains(subCatchment) {
		m.Logger().Warn("Attempted to request subcatchment [" + requestSuppliedSubCatchment + "] state not offered by the model")
		m.NotFoundError(w, r)
		return
	}

	m.respondWithApplicableActions(w, subCatchment)
}

func (m *Mux) respondWithApplicableActions(w http.ResponseWriter, subCatchment planningunit.Id) {
	m.logApplicableActionsMessage(subCatchment)
	m.sendSubcatchmentApplicableActionsResponse(w, subCatchment)
}

func (m *Mux) sendSubcatchmentApplicableActionsResponse(w http.ResponseWriter, subCatchment planningunit.Id) {
	responseActions := m.deriveApplicableActionsFor(subCatchment)

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(responseActions)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, v1ApplicableActionsHandler)
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) logApplicableActionsMessage(subCatchment planningunit.Id) {
	scenarioName := m.Attribute(scenarioNameKey).(string)
	responseMessage := fmt.Sprintf("Responding with model [%s] subcatchment [%d] applicable actions ", scenarioName, subCatchment)
	m.Logger().Info(responseMessage)
}

type ActionsArray struct {
	ApplicableActions []string
}

func (m *Mux) deriveApplicableActionsFor(subCatchment planningunit.Id) ActionsArray {
	actionsFound := m.derivedActionsFromModel(subCatchment)
	returnActions := ActionsArray{ApplicableActions: actionsFound}
	return returnActions
}

func (m *Mux) derivedActionsFromModel(subCatchment planningunit.Id) []string {
	activeActions := m.modelSolution.ActiveManagementActions[subCatchment]
	inactiveActions := m.modelSolution.InactiveManagementActions[subCatchment]

	actionsFound := make([]string, 0)

	for _, action := range activeActions {
		actionsFound = append(actionsFound, string(action))
	}
	for _, action := range inactiveActions {
		actionsFound = append(actionsFound, string(action))
	}
	return actionsFound
}
