package api

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func (m *Mux) v1modelsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.v1PostModelsHandler(w, r)
	case http.MethodGet:
		m.v1GetModelsHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetModelsHandler(w http.ResponseWriter, r *http.Request) {
	requestSuppliedModelLabel := deriveModelLabelFrom(r)

	if !m.HasAttribute(scenarioNameKey) {
		m.Logger().Warn("Attempted to request model [" + requestSuppliedModelLabel + "] with no scenario loaded")
		m.NotFoundError(w, r)
		return
	}

	modelLabel := ModelPoolLabel(requestSuppliedModelLabel)

	if !m.modelPool.HasModel(modelLabel) {
		m.Logger().Warn("Attempted to request non-instantiated model [" + requestSuppliedModelLabel + "]")
		m.NotFoundError(w, r)
		return
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.modelPool.Solution(modelLabel))

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with scenario [" + scenarioName + "] model [" + requestSuppliedModelLabel + "] state")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 models handler")
		m.Logger().Error(wrappingError)
	}
}

func deriveModelLabelFrom(r *http.Request) string {
	pathElements := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	lastElementIndex := len(pathElements) - 1
	modelLabelString := pathElements[lastElementIndex]
	return modelLabelString
}

func (m *Mux) v1PostModelsHandler(w http.ResponseWriter, r *http.Request) {
	requestSuppliedModelLabel := deriveModelLabelFrom(r)

	if !m.HasAttribute(scenarioNameKey) {
		m.Logger().Warn("Attempted to instantiate model [" + requestSuppliedModelLabel + "] with no scenario loaded")
		m.NotFoundError(w, r)
		return
	}

	modelLabel := ModelPoolLabel(requestSuppliedModelLabel)

	if modelLabel == AsIs {
		m.Logger().Warn("Attempted to instantiate model [" + requestSuppliedModelLabel + "]")
		m.MethodNotAllowedError(w, r)
		return
	}

	if requestSuppliedModelLabel != "Scratchpad" && m.modelPool.HasModel(modelLabel) {
		m.Logger().Warn("Attempted to re-instantiated model [" + requestSuppliedModelLabel + "]")
		m.MethodNotAllowedError(w, r)
		return
	}

	if m.requestContentTypeWasNotJson(r, w) {
		return
	}

	rawRequestContent := requestBodyToBytes(r)
	requestAttributes := new(attributes.Attributes)
	parseError := json.Unmarshal(rawRequestContent, requestAttributes)

	if parseError != nil {
		m.Logger().Warn("Parsing POST message content for model [" + requestSuppliedModelLabel + "] failed")
		m.RespondWithError(http.StatusBadRequest, parseError.Error(), w, r)
		return
	}

	var encoding string
	var summary string
	for _, entry := range *requestAttributes {
		if entry.Name == "Encoding" {
			encoding = entry.Value.(string)
		}
		if entry.Name == "Summary" {
			summary = entry.Value.(string)
		}
	}

	m.Logger().Info("Instantiating model [" + requestSuppliedModelLabel + "] with encoding [" + encoding + "]")

	m.modelPool.InstantiateModel(modelLabel, encoding, summary)

	restResponse := m.buildModelsPostResponse(requestSuppliedModelLabel, w)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 models handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) buildModelsPostResponse(label string, w http.ResponseWriter) *rest.Response {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(
			rest.MessageResponse{
				Type:    "SUCCESS",
				Message: "Model [" + label + "] successfully instantiated",
				Time:    rest.FormattedTimestamp(),
			},
		)

	m.Logger().Info("Responding with acknowledgement of model instantiation")

	return restResponse
}

func (m *Mux) requestContentTypeWasNotJson(r *http.Request, w http.ResponseWriter) bool {
	suppliedContentType := r.Header.Get(rest.ContentTypeHeaderKey)
	if suppliedContentType != rest.JsonMimeType {
		m.handleNonJsonContentResponse(r, w, suppliedContentType)
		return true
	}
	return false
}

func (m *Mux) handleNonJsonContentResponse(r *http.Request, w http.ResponseWriter, suppliedContentType string) {
	contentTypeError := errors.New("Request content-type of [" + suppliedContentType + "] was not the expected [" + rest.JsonMimeType + "]")
	wrappingError := errors.Wrap(contentTypeError, "v1 models handler")
	m.Logger().Warn(wrappingError)

	m.UnsupportedMediaTypeError(w, r)
}
