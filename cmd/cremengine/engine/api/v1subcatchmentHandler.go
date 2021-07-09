package api

import (
	"encoding/json"
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	ActiveAction   = "Active"
	InactiveAction = "Inactive"
)

func (m *Mux) v1subcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.v1GetSubcatchmentHandler(w, r)
	case http.MethodPut:
		m.v1PutSubcatchmentHandler(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1GetSubcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	requestSuppliedSubCatchment := deriveSubCatchmentFrom(r)

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

	m.respondWithSubcatchmentState(w, subCatchment)
}

func (m *Mux) respondWithSubcatchmentState(w http.ResponseWriter, subCatchment planningunit.Id) {
	m.logSubcatchmentStateMessage(subCatchment)
	m.sendSubcatchmentStateResponse(w, subCatchment)
}

func (m *Mux) sendSubcatchmentStateResponse(w http.ResponseWriter, subCatchment planningunit.Id) {
	responseAttributes := m.deriveResponseAttributesFor(subCatchment)

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(responseAttributes)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 subcatchment handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) logSubcatchmentStateMessage(subCatchment planningunit.Id) {
	scenarioName := m.Attribute(scenarioNameKey).(string)
	responseMessage := fmt.Sprintf("Responding with model [%s] subcatchment [%d] state", scenarioName, subCatchment)
	m.Logger().Info(responseMessage)
}

func (m *Mux) deriveResponseAttributesFor(subCatchment planningunit.Id) attributes.Attributes {
	activeActions := m.modelSolution.ActiveManagementActions[subCatchment]
	inactiveActions := m.modelSolution.InactiveManagementActions[subCatchment]

	returnAttributes := attributes.Attributes{}

	for _, action := range activeActions {
		returnAttributes = returnAttributes.Add(string(action), ActiveAction)
	}
	for _, action := range inactiveActions {
		returnAttributes = returnAttributes.Add(string(action), InactiveAction)
	}
	return returnAttributes
}

func (m *Mux) v1PutSubcatchmentHandler(w http.ResponseWriter, r *http.Request) {
	if m.modelSolution == nil {
		m.NotFoundError(w, r)
		return
	}

	requestSuppliedSubCatchment := deriveSubCatchmentFrom(r)
	subCatchment := toPlanningUnitId(requestSuppliedSubCatchment)

	if !m.modelContains(subCatchment) {
		m.NotFoundError(w, r)
		return
	}

	processingError := m.processSubcatchmentPost(w, r, subCatchment)
	if processingError != nil {
		m.reportProcessingError(w, r, processingError)
		return
	}

	restResponse := m.buildSubcatchmentResponse(w)
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "v1 POST subcatchment handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) reportProcessingError(w http.ResponseWriter, r *http.Request, processingError error) {
	m.Logger().Error(processingError)
	m.RespondWithError(http.StatusBadRequest, processingError.Error(), w, r)
}

func (m *Mux) processSubcatchmentPost(w http.ResponseWriter, r *http.Request, subCatchment planningunit.Id) error {
	scenarioName := m.Attribute(scenarioNameKey).(string)
	responseMessage := fmt.Sprintf("Processing POST of model [%s] subcatchment [%d] state", scenarioName, subCatchment)
	m.Logger().Info(responseMessage)

	requestContent := requestBodyToBytes(r)
	postedAttributes := attributes.Attributes{}

	unmnarshalError := json.Unmarshal(requestContent, &postedAttributes)
	if unmnarshalError != nil {
		return errors.Wrap(unmnarshalError, "v1 subcatchment handler")
	}

	syntaxCheckError := m.syntaxCheckPostedAttributes(postedAttributes)
	if syntaxCheckError != nil {
		return syntaxCheckError
	}

	updateModelError := m.updateModel(subCatchment, postedAttributes)
	if updateModelError != nil {
		return updateModelError
	}

	return nil
}

func (m *Mux) updateModelSolution() {
	m.modelSolution = new(solution.SolutionBuilder).
		WithId(m.model.Id()).
		ForModel(m.model).
		Build()
}

func (m *Mux) syntaxCheckPostedAttributes(postedAttributes attributes.Attributes) error {
	// TODO:  Hardcoding these is a bad code smell.
	for _, entry := range postedAttributes {
		if entry.Name != "RiverBankRestoration" && entry.Name != "HillSlopeRestoration" &&
			entry.Name != "GullyRestoration" && entry.Name != "WetlandsEstablishment" {
			baseError := errors.New("Name [" + entry.Name + "] not one of [RiverBankRestoration, HillSlopeRestoration, GullyRestoration, WetlandsEstablishment]")
			return errors.Wrap(baseError, "v1 POST subcatchment handler")
		}

		if entry.Value != InactiveAction && entry.Value != ActiveAction {
			baseError := errors.New("For named action [" + entry.Name + "], value [" + entry.Value.(string) + "] not one of [Active,Inactive]")
			return errors.Wrap(baseError, "v1 POST subcatchment handler")
		}
	}
	return nil
}

func (m *Mux) updateModel(subCatchment planningunit.Id, postedAttributes attributes.Attributes) error {
	updateErrors := compositeErrors.New("Model Update failure")
	for _, entry := range postedAttributes {
		postedEntryFound := false
		for _, action := range m.model.ManagementActions() {
			if entry.Name == string(action.Type()) && subCatchment == action.PlanningUnit() {
				postedEntryFound = true
			}
		}
		if !postedEntryFound {
			notFoundWarning := fmt.Sprintf("Subcatchment [%d], Action [%s] is not supported", subCatchment, entry.Name)
			updateErrors.AddMessage(notFoundWarning)
		}
	}

	if updateErrors.Size() > 0 {
		return updateErrors
	}

	for actionIndex, action := range m.model.ManagementActions() {
		if subCatchment != action.PlanningUnit() {
			continue
		}

		for _, entry := range postedAttributes {
			if entry.Name == string(action.Type()) {
				if entry.Value == InactiveAction {
					m.model.SetManagementAction(actionIndex, false)
				}
				if entry.Value == ActiveAction {
					m.model.SetManagementAction(actionIndex, true)
				}
				m.model.AcceptAll()
				infoMessage := fmt.Sprintf("Model subcatchment [%d], Action [%s] set to [%s]", subCatchment, entry.Name, entry.Value)
				m.Logger().Info(infoMessage)
			}
		}
	}
	m.updateModelSolution()
	return nil
}

func deriveSubCatchmentFrom(r *http.Request) string {
	pathElements := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	lastElementIndex := len(pathElements) - 1
	subCatchmentAsString := pathElements[lastElementIndex]
	return subCatchmentAsString
}

func toPlanningUnitId(subCatchmentAsString string) planningunit.Id {
	subCatchmentAsInteger, convertError := strconv.Atoi(subCatchmentAsString)
	if convertError != nil {
		panic("Should not reach here -- regular expression map should stop non-integers from being passed to handler")
	}

	subCatchment := planningunit.Id(subCatchmentAsInteger)
	return subCatchment
}

func (m *Mux) modelContains(subCatchment planningunit.Id) bool {
	for _, value := range m.modelSolution.PlanningUnits {
		if value == subCatchment {
			return true
		}
	}
	return false
}

func (m *Mux) buildSubcatchmentResponse(w http.ResponseWriter) *rest.Response {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(
			rest.MessageResponse{
				Type:    "SUCCESS",
				Message: "Subcatchment state change successfully applied",
				Time:    rest.FormattedTimestamp(),
			},
		)

	m.Logger().Info("Responding with acknowledgement of subcatchment state change ")

	return restResponse
}
