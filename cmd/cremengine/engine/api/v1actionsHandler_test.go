package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"net/http"
	"testing"
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
			TargetUrl:   baseUrl + "api/v1/model/actions",
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

	scenarioTomlText := readTestFileAsText("testdata/ValidTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
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
			TargetUrl:   baseUrl + "api/v1/model/actions",
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
		Name: "POST /model/actions request returns 501 (not implemented) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/model/actions",
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
		Name: "PUT /model/actions request returns 405 (not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseUrl + "api/v1/model/actions",
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

	scenarioTomlText := readTestFileAsText("testdata/ValidTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	context := TestContext{
		Name: "POST /model/actions request returns 501 (not implemented) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/model/actions",
			ContentType: rest.TomlMimeType,
			RequestBody: "This is not the expected CSV mime type",
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelActionsRequest_BadCsvContent_BadContentResponse(t *testing.T) {
	muxUnderTest := buildMuxUnderTest()

	scenarioTomlText := readTestFileAsText("testdata/ValidTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	context := TestContext{
		Name: "POST /model/actions request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/model/actions",
			ContentType: rest.CsvMimeType,
			RequestBody: "This text is pretending to be CSV text.",
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelActionsRequest_GoodCsvContent_OkResponse(t *testing.T) {
	muxUnderTest := buildMuxUnderTest()

	scenarioTomlText := readTestFileAsText("testdata/ValidTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	validRequestBody := readTestFileAsText("testdata/ValidActiveActions.csv")

	context := TestContext{
		Name: "POST /model/actions request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/model/actions",
			ContentType: rest.CsvMimeType,
			RequestBody: validRequestBody,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}
