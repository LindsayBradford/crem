// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type ModelAttribute string

func (ma ModelAttribute) String() string {
	return string(ma)
}

const (
	Encoding             ModelAttribute = "Encoding"
	ParetoFrontMember    ModelAttribute = "ParetoFrontMember"
	ValidAgainstScenario ModelAttribute = "ValidAgainstScenario"
	ValidationErrors     ModelAttribute = "ValidationErrors"
)

func (m *Mux) deriveExtraModelAttributes() {
	encodingOfModel := m.deriveModelActionEncoding()
	m.checkEncodingInSolutionSummary(encodingOfModel)
	m.deriveModelValidityAgainstScenario()
}

func (m *Mux) deriveModelActionEncoding() string {
	compressedModel := modelCompressor.Compress(m.model)
	encodingOfModel := compressedModel.Encoding()
	m.model.ReplaceAttribute(Encoding.String(), encodingOfModel)
	return encodingOfModel
}

func (m *Mux) deriveModelValidityAgainstScenario() bool {
	isValid, validationErrors := m.model.StateIsValid()
	m.model.ReplaceAttribute(ValidAgainstScenario.String(), isValid)
	if !isValid {
		m.handleInvalidModel(validationErrors)
	} else {
		m.handleValidModel()
	}
	return isValid
}

func (m *Mux) handleValidModel() {
	msgText := fmt.Sprintf("New model is valid against supplied scenario")
	m.Logger().Info(msgText)

	m.model.RemoveAttribute(ValidationErrors.String())
}

func (m *Mux) handleInvalidModel(validationErrors *compositeErrors.CompositeError) {
	msgText := fmt.Sprintf("New model is invalid against supplied scenario")
	m.Logger().Info(msgText)
	m.Logger().Info("Validation errors:" + validationErrors.Error())

	m.model.ReplaceAttribute(ValidationErrors.String(), validationErrors.Error())
}

func (m *Mux) updateModelSolution() {
	m.modelSolution = new(solution.SolutionBuilder).
		WithId(m.model.Id()).
		ForModel(m.model).
		Build()
}

func (m *Mux) checkEncodingInSolutionSummary(encoding string) {
	if m.solutionSetTable == nil {
		return
	}

	encodingFound := m.encodingPresentInSolutionSummaryParetoFront(encoding)
	m.attributeModelWithParetoFrontPresence(encoding, encodingFound)
}

func (m *Mux) attributeModelWithParetoFrontPresence(encoding string, encodingFound bool) {
	if encodingFound {
		m.model.ReplaceAttribute(ParetoFrontMember.String(), true)
	} else {
		m.model.ReplaceAttribute(ParetoFrontMember.String(), false)
		msgText := fmt.Sprintf(
			"New model encoding [%s] matches no pareto front member", encoding)
		m.Logger().Info(msgText)
	}
}

func (m *Mux) encodingPresentInSolutionSummaryParetoFront(encoding string) bool {
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
	return encodingFound
}

func (m *Mux) requestContentTypeWasNotCsv(r *http.Request, w http.ResponseWriter) bool {
	suppliedContentType := r.Header.Get(rest.ContentTypeHeaderKey)
	if suppliedContentType != rest.CsvMimeType {
		m.handleNonCsvContentResponse(r, w, suppliedContentType)
		return true
	}
	return false
}

func (m *Mux) handleNonCsvContentResponse(r *http.Request, w http.ResponseWriter, suppliedContentType string) {
	contentTypeError := errors.New("Request content-type of [" + suppliedContentType + "] was not the expected [" + rest.CsvMimeType + "]")
	wrappingError := errors.Wrap(contentTypeError, "v1 model actions handler")
	m.Logger().Warn(wrappingError)

	m.UnsupportedMediaTypeError(w, r)
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

func toCatchmentModel(thisModel model.Model) *catchment.Model {
	catchmentModel, isCatchmentModel := thisModel.(*catchment.Model)
	if isCatchmentModel {
		return catchmentModel
	}
	assert.That(false).WithFailureMessage("Should not get here").Holds()
	return nil
}

func readFileAsText(filePath string) string {
	if b, err := ioutil.ReadFile(filePath); err == nil {
		return string(b)
	}
	return "error reading file"
}
