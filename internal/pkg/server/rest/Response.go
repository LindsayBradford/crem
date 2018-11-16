// Copyright (c) 2018 Australian Rivers Institute.

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crem/pkg/errors"
	errors2 "github.com/pkg/errors"
)

type ErrorResponse struct {
	ErrorMessage string
	Time         string
}

type Response struct {
	ResponseCode int
	ContentType  string
	CacheControl string
	Content      string
	Writer       http.ResponseWriter

	errors *errors.CompositeError
}

func (r *Response) Initialise() *Response {
	r.errors = new(errors.CompositeError)
	r.WithContentType(DefaultResponseContentType)
	return r
}

func (r *Response) WithWriter(writer http.ResponseWriter) *Response {
	r.Writer = writer
	return r
}

func (r *Response) WithResponseCode(responseCode int) *Response {
	r.ResponseCode = responseCode
	return r
}

func (r *Response) WithCacheControlMaxAge(cacheMaxAgeInSeconds uint64) *Response {
	r.CacheControl = fmt.Sprintf("max-age=%d", cacheMaxAgeInSeconds)
	return r
}

func (r *Response) WithCacheControlPublic() *Response {
	r.CacheControl = "public"
	return r
}

func (r *Response) WithContentType(contentType string) *Response {
	r.ContentType = contentType
	return r
}

func (r *Response) WithJsonContent(content interface{}) *Response {
	r.WithContentType(JsonMimeType)

	contentAsJsonBytes, encodeError := json.MarshalIndent(content, "", "  ")
	if encodeError != nil {
		wrappingError := errors2.Wrap(encodeError, "json content encoding")
		r.errors.Add(wrappingError)
	} else {
		r.Content = string(contentAsJsonBytes)
	}
	return r
}

func (r *Response) WithTomlContent(content interface{}) *Response {
	r.WithContentType(TomlMimeType)
	contentAsString, ok := content.(string)
	if ok {
		r.Content = contentAsString
	}
	return r
}

func (r *Response) Write() error {
	r.writeHeader()
	r.writeBody()

	if r.errors.Size() > 0 {
		return r.errors
	}
	return nil
}

func (r *Response) writeHeader() {
	r.setHeaderEntries()
	r.writeHeaderContent()
}

func (r *Response) setHeaderEntries() {
	r.setCacheControlHeaderEntry()
	r.setContentTypeHeaderEntry()
}

func (r *Response) writeHeaderContent() {
	r.Writer.WriteHeader(r.ResponseCode)
}

func (r *Response) setCacheControlHeaderEntry() {
	if r.CacheControl != "" {
		r.Writer.Header().Set(CacheControlHeaderKey, r.CacheControl)
	}
}

func (r *Response) setContentTypeHeaderEntry() {
	r.Writer.Header().Set(ContentTypeHeaderKey, r.ContentType)
}

func (r *Response) writeBody() {
	fmt.Fprintf(r.Writer, r.Content)
}
