package api

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/pkg/errors"
)

const (
	v1ApplicableActionsHandler = "v1 applicable actions handler"
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
	if m.modelSolution == nil {
		m.NotFoundError(w, r)
		return
	}

	m.logApplicableActionsMessage()
	m.sendApplicableActionsResponse(w)
}

type applicableActionsWrapper struct {
	ApplicableActions map[planningunit.Id]solution.ManagementActions
}

func (m *Mux) sendApplicableActionsResponse(w http.ResponseWriter) {
	applicableActions := applicableActionsWrapper{
		ApplicableActions: m.deriveApplicableActions(),
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(applicableActions)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, v1ApplicableActionsHandler)
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) logApplicableActionsMessage() {
	scenarioName := m.Attribute(scenarioNameKey).(string)
	responseMessage := fmt.Sprintf("Responding with model [%s] applicable actions ", scenarioName)
	m.Logger().Info(responseMessage)
}

type ActionsArray struct {
	ApplicableActions []string
}

func (m *Mux) deriveApplicableActions() map[planningunit.Id]solution.ManagementActions {
	actionsMap := map[planningunit.Id]solution.ManagementActions{}

	for _, planningUnit := range m.modelSolution.PlanningUnits {
		actionsMap[planningUnit] = m.deriveApplicableActionsFor(planningUnit)
	}

	return actionsMap
}

func (m *Mux) deriveApplicableActionsFor(subCatchment planningunit.Id) solution.ManagementActions {
	activeActions := m.modelSolution.ActiveManagementActions[subCatchment]
	inactiveActions := m.modelSolution.InactiveManagementActions[subCatchment]

	actionsFound := make(solution.ManagementActions, 0)

	for _, action := range activeActions {
		actionsFound = append(actionsFound, action)
	}
	for _, action := range inactiveActions {
		actionsFound = append(actionsFound, action)
	}
	return actionsFound
}
