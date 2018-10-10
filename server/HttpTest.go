// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type HttpTestRequestContext struct {
	Method      string
	TargetUrl   string
	RequestBody string
	ContentType string
	Handler     http.HandlerFunc
}

type JsonResponseContainer struct {
	StatusCode int
	JsonMap    map[string]interface{}
}

func (context *HttpTestRequestContext) BuildJsonResponse() JsonResponseContainer {
	response := context.buildResponse()

	return JsonResponseContainer{
		StatusCode: response.StatusCode,
		JsonMap:    context.buildJsonMapFrom(response.Body),
	}
}

func (context *HttpTestRequestContext) buildResponse() *http.Response {
	w := httptest.NewRecorder()
	context.Handler(w, context.newRequest())
	return w.Result()
}

func (context *HttpTestRequestContext) newRequest() *http.Request {
	bodyReader := strings.NewReader(context.RequestBody)

	request := httptest.NewRequest(context.Method, context.TargetUrl, bodyReader)

	if context.ContentType != "" {
		request.Header.Add(ContentTypeHeaderKey, context.ContentType)
	}

	return request

}

func (context *HttpTestRequestContext) buildJsonMapFrom(responseBody io.ReadCloser) map[string]interface{} {
	jsonMap := make(map[string]interface{})

	responseBodyBytes, _ := ioutil.ReadAll(responseBody)
	json.Unmarshal(responseBodyBytes, &jsonMap)

	return jsonMap
}
