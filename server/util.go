// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LindsayBradford/crem/config"
)

const ContentTypeHeaderKey = "Content-Type"
const CacheControlHeaderKey = "Cache-Control"
const UrlPathSeparator = "/"

const TomlMimeType = "application/toml"
const JsonMimeType = "application/json"

const DefaultResposneContentType = JsonMimeType

func FormattedTimestamp() string {
	return fmt.Sprintf("%v", time.Now().Format(time.RFC3339Nano))
}

func NameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", config.LongApplicationName, config.Version)
}

func SendTextOnResponseBody(text string, w http.ResponseWriter) {
	fmt.Fprintf(w, text)
}
