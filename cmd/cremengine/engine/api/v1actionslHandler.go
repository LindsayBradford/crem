package api

import (
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

func (m *Mux) v1GetActionsHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.NotFoundError(w, r)
		return
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.modelSolution.ActiveManagementActions)

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with model [" + scenarioName + "] state")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PostActionsHandler(w http.ResponseWriter, r *http.Request) {
	m.RespondWithError(http.StatusNotImplemented, "Behaviour not yet implemented", w, r)
}
