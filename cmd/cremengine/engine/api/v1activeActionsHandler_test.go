package api

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	. "github.com/onsi/gomega"
)

//go:embed testdata/InvalidActiveActions.csv
var invalidActionsCsvContent string

//go:embed testdata/ValidActiveActions.csv
var validRActionsCsvContent string

const (
	activeActionsPath = "active"
)

func TestFirstActionsGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "GET /api/v1/model/actions request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestGetValidModelActionsResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	getContext := TestContext{
		Name: "GET /api/v1/model request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestModelActionsRequestNoScenario_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "PUT /model/actions request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelActionsPutRequest_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "POST /model/actions request returns 405 (not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelActionsRequestNotCsv_NotFoundResponse(t *testing.T) {
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	context := TestContext{
		Name: "PUT /model/actions request returns 415 (unsupported media type) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			ContentType: rest.TomlMimeType,
			RequestBody: "This is not the expected CSV mime type",
		},
		ExpectedResponseStatus: http.StatusUnsupportedMediaType,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelActionsRequest_BadCsvContent_BadContentResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	context := TestContext{
		Name: "PUT /model/actions request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			ContentType: rest.CsvMimeType,
			RequestBody: "This text is pretending to be CSV text.",
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelActionsRequest_BadCsvCells_BadContentResponse(t *testing.T) {
	g := NewGomegaWithT(t)
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when

	context := TestContext{
		Name: "PUT /model/actions request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			ContentType: rest.CsvMimeType,
			RequestBody: invalidActionsCsvContent,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	responseBody := verifyResponseStatusCode(muxUnderTest, context)
	jsonResponseBody := responseBody.JsonMap
	g.Expect(jsonResponseBody["Type"]).To(Equal("ERROR"))

	muxUnderTest.Shutdown()
}

func TestModelActionsRequest_GoodCsvContent_OkResponse(t *testing.T) {
	g := NewGomegaWithT(t)
	muxUnderTest := buildMuxUnderTest()

	// when
	scenarioPostContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, scenarioPostContext)

	// when
	actionsPostContext := TestContext{
		Name: "PUT /model/actions request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
			ContentType: rest.CsvMimeType,
			RequestBody: validRActionsCsvContent,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	actionsPostResponseContainer := verifyResponseStatusCode(muxUnderTest, actionsPostContext)
	jsonPostResponse := actionsPostResponseContainer.JsonMap
	g.Expect(jsonPostResponse["Type"]).To(Equal("SUCCESS"))

	//when

	actionsGetContext := TestContext{
		Name: "GET /model/actions request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    "GET",
			TargetUrl: baseActionsUrl + rest.UrlPathSeparator + activeActionsPath,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	actionsGetResponseContainer := verifyResponseStatusCode(muxUnderTest, actionsGetContext)

	g.Expect(actionsGetResponseContainer.JsonMap["ActiveManagementActions"]).To(Not(BeNil()),
		actionsPostContext.Name+" should return an ActiveManagementActions json map")

	rawActionsMap := actionsGetResponseContainer.JsonMap["ActiveManagementActions"]
	actionsMap, isStringMap := rawActionsMap.(map[string]interface{})
	if !isStringMap {
		g.Expect("").ToNot(BeEmpty(), "ActiveManagementActions map didn't match expected type")
	}

	rawSc17Array := actionsMap["17"]
	sc17Array, sc17ValueIsArray := rawSc17Array.([]interface{})
	if !sc17ValueIsArray {
		g.Expect("").ToNot(BeEmpty(), "ActiveManagementActions[17] map didn't match expected type")
	}

	g.Expect(len(sc17Array)).To(BeNumerically("==", 1))
	g.Expect(sc17Array[0]).To(Equal("GullyRestoration"), actionsPostContext.Name+" Subcatchment 17 expected to have gully restoration")

	rawSc18Array := actionsMap["18"]
	sc18Array, sc18ValueIsArray := rawSc18Array.([]interface{})
	if !sc18ValueIsArray {
		g.Expect("").ToNot(BeEmpty(), "ActiveManagementActions[18] map didn't match expected type")
	}

	g.Expect(len(sc18Array)).To(BeNumerically("==", 1))
	g.Expect(sc18Array[0]).To(Equal("RiverBankRestoration"), actionsPostContext.Name+" Subcatchment 18 expected to have river bank restoration")

	rawSc19Array := actionsMap["19"]
	sc19Array, sc19ValueIsArray := rawSc19Array.([]interface{})
	if !sc19ValueIsArray {
		g.Expect("").ToNot(BeEmpty(), "ActiveManagementActions[19] map didn't match expected type")
	}

	g.Expect(len(sc19Array)).To(BeNumerically("==", 1))
	g.Expect(sc19Array[0]).To(Equal("HillSlopeRestoration"), actionsPostContext.Name+" Subcatchment 19 expected to have hillslope restoration")

	rawSc20Array := actionsMap["20"]
	g.Expect(rawSc20Array).To(BeNil())

	rawSc21Array := actionsMap["21"]
	sc21Array, sc21ValueIsArray := rawSc21Array.([]interface{})
	if !sc21ValueIsArray {
		g.Expect("").ToNot(BeEmpty(), "ActiveManagementActions[21] map didn't match expected type")
	}

	g.Expect(len(sc21Array)).To(BeNumerically("==", 1))
	g.Expect(sc21Array[0]).To(Equal("WetlandsEstablishment"), actionsPostContext.Name+" Subcatchment 21 expected to have wetlands establishment")

	muxUnderTest.Shutdown()
}
