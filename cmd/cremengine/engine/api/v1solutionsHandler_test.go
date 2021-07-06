package api

import (
	_ "embed"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"net/http"
	"testing"
)

//go:embed testdata/ValidSolutions-Summary.csv
var validSolutions string

//go:embed testdata/InvalidSolutions-Summary.csv
var invalidSolutions string

func TestFirstSolutionsGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "GET /scenario request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    "GET",
			TargetUrl: baseUrl + "api/v1/solutions",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestPostSolutionsResource_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /solutions text request returns 405 (not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: "here is some text that should be TOML",
			ContentType: rest.TextMimeType,
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then

	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestPostSolutionsCsvResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 200 (ok) response",
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
	postContext = TestContext{
		Name: "POST /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: validSolutions,
			ContentType: rest.CsvMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestGetSolutionsCsvResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 200 (ok) response",
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
	postContext = TestContext{
		Name: "POST /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: validSolutions,
			ContentType: rest.CsvMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	postContext = TestContext{
		Name: "GET /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    "GET",
			TargetUrl: baseUrl + "api/v1/solutions",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestGetSolutionsNoScenario_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 200 (ok) response",
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
	postContext = TestContext{
		Name: "GET /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    "GET",
			TargetUrl: baseUrl + "api/v1/solutions",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestGetSolutionsNotCsvResource_UnsupportedMediaTypeResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 200 (ok) response",
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
	postContext = TestContext{
		Name: "POST /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: validSolutions,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusUnsupportedMediaType,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestGetSolutionsInvalidCsv_BadContentResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 200 (ok) response",
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
	postContext = TestContext{
		Name: "POST /solutions text request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: invalidSolutions,
			ContentType: rest.CsvMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

//func TestPostSolutionsTextResource_BadRequestResponse(t *testing.T) {
//	// given
//	muxUnderTest := buildMuxUnderTest()
//
//	// when
//	postContext := TestContext{
//		Name: "POST /scenario text request returns 200 (ok) response",
//		T:    t,
//		Request: httptest.HttpTestRequestContext{
//			Method:      "POST",
//			TargetUrl:   baseUrl + "api/v1/scenario",
//			RequestBody: "This isn't TOML",
//			ContentType: rest.TomlMimeType,
//		},
//		ExpectedResponseStatus: http.StatusBadRequest,
//	}
//
//	// then
//	verifyResponseStatusCode(muxUnderTest, postContext)
//
//	muxUnderTest.Shutdown()
//}
//
//func TestPostValidSolutionsResource_OkResponse(t *testing.T) {
//	// given
//	muxUnderTest := buildMuxUnderTest()
//
//	// when
//	postContext := TestContext{
//		Name: "POST /scenario request returns 202 (accepted) response",
//		T:    t,
//		Request: httptest.HttpTestRequestContext{
//			Method:      "POST",
//			TargetUrl:   baseUrl + "api/v1/scenario",
//			RequestBody: validScenarioTomlText,
//			ContentType: rest.TomlMimeType,
//		},
//		ExpectedResponseStatus: http.StatusOK,
//	}
//
//	// then
//	verifyResponseStatusCode(muxUnderTest, postContext)
//
//	// when
//	getContext := TestContext{
//		Name: "GET /scenario request returns 200 (ok) response",
//		T:    t,
//		Request: httptest.HttpTestRequestContext{
//			Method:      "GET",
//			TargetUrl:   baseUrl + "api/v1/scenario",
//			RequestBody: "here is some text",
//		},
//		ExpectedResponseStatus: http.StatusOK,
//	}
//
//	// then
//	verifyResponseStatusCode(muxUnderTest, getContext)
//
//	muxUnderTest.Shutdown()
//}
//
