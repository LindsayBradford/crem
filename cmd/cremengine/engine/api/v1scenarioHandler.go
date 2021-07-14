package api

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
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

	restResponse := m.buildScenarioGetResponse(w)
	m.logScenarioGetResponse()
	writeError := restResponse.Write()

	m.handleScenarioGetWriteError(writeError)
}

func (m *Mux) buildScenarioGetResponse(w http.ResponseWriter) *rest.Response {
	responseText := m.Attribute(scenarioTextKey).(string)

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithTomlContent(responseText)
	return restResponse
}

func (m *Mux) logScenarioGetResponse() {
	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with scenario [" + scenarioName + "] configuration")
}

func (m *Mux) handleScenarioGetWriteError(writeError error) {
	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 scenario handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PostScenarioHandler(w http.ResponseWriter, r *http.Request) {
	if m.requestContentTypeWasNotToml(r, w) {
		return
	}

	scenarioConfig, retrievalError := m.processScenarioPostText(w, r)
	if retrievalError != nil {
		return
	}

	modelErrors := m.deriveDefaultModelForScenario(w, r, scenarioConfig)
	if modelErrors != nil {
		return
	}

	restResponse := m.buildScenarioPostResponse(w)
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 scenario handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) deriveDefaultModelForScenario(w http.ResponseWriter, r *http.Request, scenarioConfig *data.ScenarioConfig) error {
	interpretedModel := m.modelConfigInterpreter.Interpret(&scenarioConfig.Model).Model()
	if modelAsCatchmentModel, isCatchmentModel := interpretedModel.(*catchment.Model); isCatchmentModel {
		m.rememberModelState(modelAsCatchmentModel, scenarioConfig)
	}
	if m.modelConfigInterpreter.Errors() != nil {
		m.handleModelInterpreterErrors(w, r, m.modelConfigInterpreter.Errors())
		return m.modelConfigInterpreter.Errors()
	}
	return nil
}

func (m *Mux) processScenarioPostText(w http.ResponseWriter, r *http.Request) (*data.ScenarioConfig, error) {
	requestContent := requestBodyToString(r)
	config, retrievalError := data.RetrieveScenarioConfigFromString(requestContent)

	if retrievalError != nil {
		m.handleScenarioRetrievalErrors(w, r, retrievalError)
		return config, retrievalError
	}

	m.rememberScenarioAttributeState(config, requestContent)
	return config, nil
}

func (m *Mux) handleScenarioRetrievalErrors(w http.ResponseWriter, r *http.Request, retrieveError error) {
	wrappingError := errors.Wrap(retrieveError, "v1 POST scenario handler")
	m.Logger().Error(wrappingError)
	m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
}

func (m *Mux) rememberScenarioAttributeState(config *data.ScenarioConfig, requestContent string) {
	m.ReplaceAttribute(scenarioNameKey, config.Scenario.Name)
	m.Logger().Info("Scenario configuration [" + config.Scenario.Name + "] successfully retrieved")

	m.ReplaceAttribute(scenarioTextKey, requestContent)
}

func (m *Mux) rememberModelState(modelAsCatchmentModel *catchment.Model, config *data.ScenarioConfig) {
	m.model = modelAsCatchmentModel
	m.model.Initialise(model.AsIs)
	m.model.SetId(config.Scenario.Name)
	m.model.ReplaceAttribute("ParetoFrontMember", "No")

	m.solutionPool = NewSolutionPool(modelAsCatchmentModel)
}

func (m *Mux) handleModelInterpreterErrors(w http.ResponseWriter, r *http.Request, interpreterError error) {
	wrappingError := errors.Wrap(interpreterError, "v1 POST scenario handler")
	m.Logger().Error(wrappingError)
	m.RespondWithError(http.StatusBadRequest, wrappingError.Error(), w, r)
}

func (m *Mux) buildScenarioPostResponse(w http.ResponseWriter) *rest.Response {
	m.modelSolution = new(solution.SolutionBuilder).
		WithId(m.model.Id()).
		ForModel(m.model).
		Build()

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(
			rest.MessageResponse{
				Type:    "SUCCESS",
				Message: "Scenario configuration successfully posted",
				Time:    rest.FormattedTimestamp(),
			},
		)

	m.Logger().Info("Responding with acknowledgement of scenario configuration receipt")
	return restResponse
}

func (m *Mux) requestContentTypeWasNotToml(r *http.Request, w http.ResponseWriter) bool {
	suppliedContentType := r.Header.Get(rest.ContentTypeHeaderKey)
	if suppliedContentType != rest.TomlMimeType {
		m.handleNonTomlContentResponse(r, w, suppliedContentType)
		return true
	}
	return false
}

func (m *Mux) handleNonTomlContentResponse(r *http.Request, w http.ResponseWriter, suppliedContentType string) {
	contentTypeError := errors.New("Request content-type of [" + suppliedContentType + "] was not the expected [" + rest.TomlMimeType + "]")
	wrappingError := errors.Wrap(contentTypeError, "v1 POST scenario handler")
	m.Logger().Warn(wrappingError)

	m.MethodNotAllowedError(w, r)
}

func (m *Mux) SetScenario(scenarioFilePath string) {
	config, retrievalError := data.RetrieveScenarioConfigFromFile(scenarioFilePath)

	if retrievalError != nil {
		m.Logger().Warn(retrievalError)
		return
	}

	m.ReplaceAttribute(scenarioNameKey, config.Scenario.Name)
	m.Logger().Info("Scenario configuration [" + config.Scenario.Name + "] successfully retrieved")

	configFileContent := readFileAsText(scenarioFilePath)
	m.ReplaceAttribute(scenarioTextKey, configFileContent)

	interpretedModel := m.modelConfigInterpreter.Interpret(&config.Model).Model()
	if modelAsCatchmentModel, isCatchmentModel := interpretedModel.(*catchment.Model); isCatchmentModel {
		m.rememberModelState(modelAsCatchmentModel, config)
	}
	if m.modelConfigInterpreter.Errors() != nil {
		m.Logger().Warn(m.modelConfigInterpreter.Errors())
	}

	m.modelSolution = new(solution.SolutionBuilder).
		WithId(m.model.Id()).
		ForModel(m.model).
		Build()
}

func readFileAsText(filePath string) string {
	if b, err := ioutil.ReadFile(filePath); err == nil {
		return string(b)
	}
	return "error reading file"
}
