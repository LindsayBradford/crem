package api

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/attributes"
	. "github.com/onsi/gomega"
)

//go:embed testdata/ValidTestScenario.toml
var validScenarioTomlConfig string

const (
	scenarioUrl = baseUrl + "api/v1/scenario"

	baseSubcatchmentUrl  = baseUrl + "api/v1/model/subcatchment"
	validSubcatchment    = "18"
	validSubcatchmentUrl = baseSubcatchmentUrl + rest.UrlPathSeparator + validSubcatchment
)

func TestFirstSubcatchmentGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + rest.UrlPathSeparator + validSubcatchment

	// when
	context := TestContext{
		Name: http.MethodGet + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestMissingSubcatchmentGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/1"
	getContext := TestContext{
		Name: http.MethodGet + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestInvalidSubcatchmentGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/nope"
	getContext := TestContext{
		Name: http.MethodGet + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestDeleteValidSubcathmentResource_BadMethodResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	getContext := TestContext{
		Name: http.MethodDelete + " " + validSubcatchmentUrl + " request returns 405 (method not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodDelete,
			TargetUrl: validSubcatchmentUrl,
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetValidSubcathmentResource_OkResponse(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	getContext := TestContext{
		Name: http.MethodGet + " " + validSubcatchmentUrl + " request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: validSubcatchmentUrl,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	response := verifyResponseStatusCode(muxUnderTest, getContext)
	jsonResponse := response.JsonMap
	g.Expect(len(jsonResponse)).To(BeNumerically("==", 0))

	muxUnderTest.Shutdown()
}

func TestFirstSubcatchmentPutRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + rest.UrlPathSeparator + validSubcatchment

	// when
	context := TestContext{
		Name: http.MethodPut + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodPut,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestMissingSubcatchmentPutRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/1"
	getContext := TestContext{
		Name: http.MethodPut + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodPut,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestInvalidSubcatchmentPutRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/nope"
	getContext := TestContext{
		Name: http.MethodPut + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodPut,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPostValidSubcatchmentResource_OkResponse(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionAttributes := attributes.Attributes{}.
		Add("GullyRestoration", ActiveAction).
		Add("RiverBankRestoration", ActiveAction)

	actionStatusBytes, _ := json.Marshal(actionAttributes)

	// when
	postContext := TestContext{
		Name: http.MethodPut + " " + validSubcatchmentUrl + " request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPut,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(actionStatusBytes),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	postResponse := verifyResponseStatusCode(muxUnderTest, postContext)
	jsonPostResponse := postResponse.JsonMap
	g.Expect(jsonPostResponse["Type"]).To(Equal("SUCCESS"))

	// when
	getContext := TestContext{
		Name: http.MethodGet + " " + validSubcatchmentUrl + " request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: validSubcatchmentUrl,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)

	muxUnderTest.Shutdown()
}

func TestPutInvalidSubcatchmentJson_ErrorResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	invalidJson := "This is not valid\" Json}"

	// when
	getContext := TestContext{
		Name: http.MethodPut + " " + validSubcatchmentUrl + " request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPut,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(invalidJson),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPutInvalidSubcatchmentActionResource_ErrorResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionAttributes := attributes.Attributes{}.
		Add("NonExistentAction", ActiveAction)

	actionStatusBytes, _ := json.Marshal(actionAttributes)

	// when
	getContext := TestContext{
		Name: http.MethodPut + " " + validSubcatchmentUrl + " request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPut,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(actionStatusBytes),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPutInvalidSubcatchmentActionStateResource_ErrorResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionAttributes := attributes.Attributes{}.
		Add("GullyRestoration", "ThisIsNotAValidState")

	actionStatusBytes, _ := json.Marshal(actionAttributes)

	// when
	getContext := TestContext{
		Name: http.MethodPut + " " + validSubcatchmentUrl + " request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPut,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(actionStatusBytes),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPutInvalidSubcatchmentActionSemantics_ErrorResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionAttributes := attributes.Attributes{}.
		Add("HillSlopeRestoration", "Active"). // SC 18 doesn't have a WetlandsEstablishment action
		Add("WetlandsEstablishment", "Active") // SC 18 doesn't have a WetlandsEstablishment action

	actionStatusBytes, _ := json.Marshal(actionAttributes)

	// when
	getContext := TestContext{
		Name: http.MethodPut + " " + validSubcatchmentUrl + " request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPut,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(actionStatusBytes),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func buildValidScenario(t *testing.T, muxUnderTest *Mux) {

	// when
	postContext := TestContext{
		Name: http.MethodPost + scenarioUrl + " request returns 200 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPost,
			TargetUrl:   scenarioUrl,
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)
}
