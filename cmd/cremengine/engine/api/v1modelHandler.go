package api

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/pkg/errors"
	"net/http"
)

const v1modelHandler = "v1 model handler"

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
		wrappingError := errors.Wrap(writeError, v1modelHandler)
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
		wrappingError := errors.Wrap(parseError, v1modelHandler)
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
			updateError := m.updateModelWithEncoding(encoding)

			if updateError != nil {
				wrappingError := errors.Wrap(updateError, v1modelHandler)
				m.Logger().Error(wrappingError)
				m.RespondWithError(http.StatusBadRequest, updateError.Error(), w, r)
				return
			}
		}
	}

	restResponse := m.buildModelPatchResponse(w)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, v1modelHandler)
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) updateModelWithEncoding(encoding string) error {
	encodingError := m.reInitialiseModelWithEncoding(encoding)
	if encodingError != nil {
		return encodingError
	}
	m.updateModelSolution()

	return nil
}

func (m *Mux) reInitialiseModelWithEncoding(encoding string) error {
	newModel := m.model.DeepClone()

	compressedModel := modelCompressor.Compress(newModel)
	decodingError := compressedModel.Decode(encoding)

	if decodingError != nil {
		wrappingError := errors.Wrap(decodingError, v1modelHandler)
		m.Logger().Warn(wrappingError)

		wrappingError = errors.Wrap(errors.New("Leaving model in current state"), v1modelHandler)
		m.Logger().Warn(wrappingError)

		return decodingError
	}

	modelCompressor.Decompress(compressedModel, newModel)

	m.model = toCatchmentModel(newModel)
	m.deriveExtraModelAttributes()

	return nil
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
