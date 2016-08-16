package common

import (
	"encoding/json"
	"github.com/gophercloud/gophercloud"
	"strings"
)

// ErrorResponse represents failed response
type ErrorResponse struct {
	gophercloud.ErrUnexpectedResponseCode
	Errors []Error `json:"errors" required:"true"`
}

// Error represents a Magnum API error message
type Error struct {
	Status    int    `json:"status"`
	Code      string `json:"code"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	RequestID string `json:"request_id,omitempty"`
}

// Error returns the error message details from the server
func (e *ErrorResponse) Error() string {
	errors := []string{}
	for _, e := range e.Errors {
		errors = append(errors, e.Detail)
	}
	return strings.Join(errors, "\n")
}

func (e *ErrorResponse) unwrapError(r gophercloud.ErrUnexpectedResponseCode) bool {
	e.ErrUnexpectedResponseCode = r
	err := json.Unmarshal(r.Body, &e)
	return err == nil
}

// Error400 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error400(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}

	return gophercloud.ErrDefault400{ErrUnexpectedResponseCode: r}
}

// Error401 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error401(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}
	return gophercloud.ErrDefault401{ErrUnexpectedResponseCode: r}
}

// Error404 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error404(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}
	return gophercloud.ErrDefault404{ErrUnexpectedResponseCode: r}
}

// Error405 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error405(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}

	return gophercloud.ErrDefault405{ErrUnexpectedResponseCode: r}
}

// Error408 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error408(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}

	return gophercloud.ErrDefault408{ErrUnexpectedResponseCode: r}
}

// Error409 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error409(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}

	return gophercloud.ErrDefault409{ErrUnexpectedResponseCode: r}
}

// Error429 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error429(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}

	return gophercloud.ErrDefault429{ErrUnexpectedResponseCode: r}
}

// Error500 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error500(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}

	return gophercloud.ErrDefault500{ErrUnexpectedResponseCode: r}
}

// Error503 extracts the actual error message from the body of the response
func (e *ErrorResponse) Error503(r gophercloud.ErrUnexpectedResponseCode) error {
	if e.unwrapError(r) {
		return e
	}

	return gophercloud.ErrDefault503{ErrUnexpectedResponseCode: r}
}
