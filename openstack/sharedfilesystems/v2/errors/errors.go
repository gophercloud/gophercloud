package errors

import (
	"encoding/json"
	"errors"

	"github.com/gophercloud/gophercloud/v2"
)

type ManilaError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ErrorDetails map[string]ManilaError

// error types from provider_client.go
func ExtractErrorInto(rawError error, errorDetails *ErrorDetails) (err error) {
	var codeError gophercloud.ErrUnexpectedResponseCode
	if errors.As(rawError, &codeError) {
		return json.Unmarshal(codeError.Body, errorDetails)
	} else {
		return errors.New("Unable to extract detailed error message")
	}
}
