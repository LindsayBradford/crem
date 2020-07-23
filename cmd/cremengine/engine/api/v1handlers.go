package api

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
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

	interpretedModel := m.modelConfigInterpreter.Interpret(&config.Model).Model()
	if modelAsCatchmentModel, isCatchmentModel := interpretedModel.(*catchment.Model); isCatchmentModel {
		m.model = modelAsCatchmentModel
		m.model.InitialiseToAsIsState()
		m.model.SetId(config.Scenario.Name)
	} else {
		interpreterError := m.modelConfigInterpreter.Errors()
		wrappingError := errors.Wrap(interpreterError, "v1 POST scenario handler")
		m.Logger().Error(wrappingError)
		m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
		return
	}

	if m.model.ParameterErrors() != nil {
		parameterErrors := m.model.ParameterErrors()
		wrappingError := errors.Wrap(parameterErrors, "v1 POST scenario handler")
		m.Logger().Error(wrappingError)
		m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
		return
	}

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

func (m *Mux) v1modelHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.v1PostModelHandler(w, r)
	case http.MethodGet:
		m.v1GetModelHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetModelHandler(w http.ResponseWriter, r *http.Request) {
	if m.model == nil {
		m.NotFoundError(w, r)
		return
	}

	solution := new(solution.SolutionBuilder).
		WithId(m.model.Id()).
		ForModel(m.model).
		Build()

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(solution)

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with model [" + scenarioName + "] status")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PostModelHandler(w http.ResponseWriter, r *http.Request) {
	m.RespondWithError(http.StatusNotFound, "Behaviour not yet implemented", w, r)
}

func requestBodyToString(r *http.Request) string {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	return string(responseBodyBytes)
}
