package webhooks

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations. Call its Extract method to
// interpret it as a Version.
type GetResult struct {
	commonResult
}

type TriggerResult struct {
	commonResult
}

// Webhook represents a detailed webhook
type Webhook struct {
	Action string `json:"action"`
}

// Extract retrieves the response action
func (r commonResult) Extract() (string, error) {
	var s *Webhook
	err := r.ExtractInto(&s)
	if err != nil {
		fmt.Println("Error WebHooks ExtractAction")
	}
	return s.Action, err
}

// ExtractAction retrieves the response action
func ExtractAction(r commonResult) (string, error) {
	var s *Webhook
	err := r.ExtractInto(&s)
	if err != nil {
		fmt.Println("Error WebHooks ExtractAction")
	}
	return s.Action, err
}
