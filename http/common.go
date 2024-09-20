package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// +HttpError+ is the web-only interface for how any api errors are rendered
// This should not be used internally
type HttpError struct {
	HttpCode  int    `json:"http_code"`
	Msg       string `json:"error_message"`
	ErrorCode string `json:"error_code"`
}

// Renders an error with the struct of +HttpError+ and the status code of HttpError's HttpCode
func renderError(httpError HttpError, writer http.ResponseWriter) {
	writer.WriteHeader(httpError.HttpCode)
	json.NewEncoder(writer).Encode(httpError)
	fmt.Println(httpError)
}
