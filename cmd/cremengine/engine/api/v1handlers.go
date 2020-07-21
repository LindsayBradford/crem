package api

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func (m *Mux) v1scenarioHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.v1PostScenarioHandler(w, r)
	case http.MethodGet:
		m.v1GetScenarioHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetScenarioHandler(w http.ResponseWriter, r *http.Request) {
	if !m.HasAttribute(scenarioTextKey) {
		m.NotFoundError(w, r)
		return
	}

	responseText := m.Attribute(scenarioTextKey).(string)

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithTomlContent(responseText)

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with scenario [" + scenarioName + "] configuration")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 scenario handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PostScenarioHandler(w http.ResponseWriter, r *http.Request) {
	if m.requestContentTypeWasNotToml(r, w) {
		return
	}

	requestContent := requestBodyToString(r)
	config, retrieveError := data.RetrieveScenarioConfigFromString(requestContent)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "v1 POST scenario handler")
		m.Logger().Error(wrappingError)
		m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
		return
	}

	m.ReplaceAttribute(scenarioNameKey, config.Scenario.Name)
	m.Logger().Info("Scenario configuration [" + config.Scenario.Name + "] successfully retrieved")

	m.ReplaceAttribute(scenarioTextKey, requestContent)

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(
			rest.MessageResponse{
				Message: "Scenario configuration successfully posted",
				Time:    rest.FormattedTimestamp(),
			},
		)

	m.Logger().Info("Responding with acknowledgement of scenario configuration receipt")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 scenario handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) requestContentTypeWasNotToml(r *http.Request, w http.ResponseWriter) bool {
	if r.Header.Get(rest.ContentTypeHeaderKey) != rest.TomlMimeType {
		m.MethodNotAllowedError(w, r)
		return true
	}
	return false
}

func requestBodyToString(r *http.Request) string {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	return string(responseBodyBytes)
}
