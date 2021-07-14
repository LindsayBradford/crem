package api

import (
	"encoding/json"
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/pkg/errors"
	"net/http"
)

func (m *Mux) v1modelHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.v1GetModelHandler(w, r)
	case http.MethodPatch:
		m.v1PatchModelHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetModelHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.Logger().Warn("Attempted to get model resource with no scenario loaded")
		m.NotFoundError(w, r)
		return
	}

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.modelSolution)

	scenarioName := m.Attribute(scenarioNameKey).(string)
	m.Logger().Info("Responding with model [" + scenarioName + "] state")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) v1PatchModelHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.Logger().Warn("Attempted to patch model resource attributes with no scenario loaded.")
		m.NotFoundError(w, r)
		return
	}

	if m.requestContentTypeWasNotJson(r, w) {
		return
	}

	rawRequestContent := requestBodyToBytes(r)
	requestAttributes := new(attributes.Attributes)
	parseError := json.Unmarshal(rawRequestContent, requestAttributes)

	if parseError != nil {
		wrappingError := errors.Wrap(parseError, "v1 model handler")
		m.Logger().Error(wrappingError)
		m.Logger().Error("Parsing PATCH message content for model failed")

		m.RespondWithError(http.StatusBadRequest, parseError.Error(), w, r)
		return
	}

	m.Logger().Info("Joining newly supplied attributes to current model")
	m.model.JoiningAttributes(*requestAttributes)

	for _, entry := range *requestAttributes {
		if entry.Name == "Encoding" {
			encoding := entry.Value.(string)
			m.Logger().Info("Re-initialising model with attribute-supp[ied alternate encoding [" + encoding + "]")
			m.updateModelWithEncoding(encoding)
		}
	}

	restResponse := m.buildModelPatchResponse(w)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 model handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) updateModelWithEncoding(encoding string) {
	m.reInitialiseModelWithEncoding(encoding)
	m.updateModelSolution()
}

func (m *Mux) reInitialiseModelWithEncoding(encoding string) {
	newModel := m.model.DeepClone()

	compressedModel := modelCompressor.Compress(newModel)
	compressedModel.Decode(encoding)
	modelCompressor.Decompress(compressedModel, newModel)

	m.model = toCatchmentModel(newModel)
	m.deriveExtraModelAttributes()
}

func (m *Mux) checkEncodingInSolutionSummary(encoding string) {
	if m.solutionSetTable == nil {
		return
	}

	colSize, rowSize := m.solutionSetTable.ColumnAndRowSize()
	var (
		labelIndex    = uint(0)
		encodingIndex = colSize - 2
		encodingFound = false
	)

	for rowIndex := uint(1); rowIndex < rowSize; rowIndex++ {
		if encoding == m.solutionSetTable.CellString(encodingIndex, rowIndex) {
			encodingFound = true
			label := m.solutionSetTable.CellString(labelIndex, rowIndex)
			msgText := fmt.Sprintf(
				"New model's encoding [%s] matches pareto front solution set member [%s]", encoding, label)
			m.Logger().Info(msgText)
		}
	}

	if encodingFound {
		m.model.ReplaceAttribute("ParetoFrontMember", "Yes")
	} else {
		m.model.ReplaceAttribute("ParetoFrontMember", "No")
		msgText := fmt.Sprintf(
			"New model encoding [%s] matches no solution set member", encoding)
		m.Logger().Info(msgText)
	}
}

func (m *Mux) buildModelPatchResponse(w http.ResponseWriter) *rest.Response {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(
			rest.MessageResponse{
				Type:    "SUCCESS",
				Message: "Model resource successfully patched",
				Time:    rest.FormattedTimestamp(),
			},
		)

	m.Logger().Info("Responding with acknowledgement of model patch")

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
	wrappingError := errors.Wrap(contentTypeError, "v1 model handler")
	m.Logger().Error(wrappingError)

	m.UnsupportedMediaTypeError(w, r)
}
