package common

import (
	"encoding/json"
	"github.com/gophercloud/gophercloud"
)

// ErrorResponse represents failed response
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error represents a Magnum API error message
type Error struct {
	Status    int    `json:"status"`
	Code      string `json:"code"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	RequestID string `json:"request_id,omitempty"`
}

func (r ErrorResponse) Error() string {
	var msg string
	for _, e := range r.Errors {
		msg += e.Detail + "\n"
	}

	return msg
}

type ErrDeleteFailed struct{}

// Error400 extracts the actual error message from the body of the response
func (d ErrDeleteFailed) Error400(r gophercloud.ErrUnexpectedResponseCode) error {
	var s *ErrorResponse
	err := json.Unmarshal(r.Body, &s)
	if err != nil {
		return gophercloud.ErrDefault400{ErrUnexpectedResponseCode: r}
	}

	return s
}

func (d ErrDeleteFailed) Error() string {
	return "Unable to delete bay"
}
