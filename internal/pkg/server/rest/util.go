// Copyright (c) 2018 Australian Rivers Institute.

package rest

import (
	"fmt"
	"net/http"
	"time"
)

const ContentTypeHeaderKey = "Content-Type"
const CacheControlHeaderKey = "Cache-Control"
const UrlPathSeparator = "/"

const TomlMimeType = "application/toml"
const JsonMimeType = "application/json"
const TextMimeType = "text/plain"
const CsvMimeType = "text/csv"

const DefaultResponseContentType = JsonMimeType

type HandlerFunc http.HandlerFunc

func FormattedTimestamp() string {
	return fmt.Sprintf("%v", time.Now().Format(time.RFC3339Nano))
}

func NameAndVersionString() string {
	// TODO:  This should be application specific. Not a general utility function.
	return fmt.Sprintf("%s, version %s", "Some application name", "<someVersion>")
}

func SendTextOnResponseBody(text string, w http.ResponseWriter) {
	fmt.Fprintf(w, text)
}
