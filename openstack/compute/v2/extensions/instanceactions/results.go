package instanceactions

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// InstanceAction represents an instance action.
type InstanceAction struct {
	Action       string    `json:"action"`
	InstanceUUID string    `json:"instance_uuid"`
	Message      string    `json:"message"`
	ProjectID    string    `json:"project_id"`
	RequestID    string    `json:"request_id"`
	StartTime    time.Time `json:"-"`
	UserID       string    `json:"user_id"`
}

// UnmarshalJSON converts our JSON API response into our instance action struct
func (i *InstanceAction) UnmarshalJSON(b []byte) error {
	type tmp InstanceAction
	var s struct {
		tmp
		StartTime gophercloud.JSONRFC3339MilliNoZ `json:"start_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*i = InstanceAction(s.tmp)

	i.StartTime = time.Time(s.StartTime)

	return err
}

// InstanceActionPage abstracts the raw results of making a ListInstanceActiones() request
// against the API. As OpenStack extensions may freely alter the response bodies
// of structures returned to the client, you may only safely access the data
// provided through the ExtractInstanceActions call.
type InstanceActionPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if an InstanceActionPage contains no instance actions.
func (r InstanceActionPage) IsEmpty() (bool, error) {
	instanceactions, err := ExtractInstanceActions(r)
	return len(instanceactions) == 0, err
}

// ListInstanceActiones() call, producing a map of instanceActions.
func ExtractInstanceActions(r pagination.Page) ([]InstanceAction, error) {
	var resp []InstanceAction
	err := ExtractInstanceActionsInto(r, &resp)
	return resp, err
}

type InstanceActionResult struct {
	gophercloud.Result
}

// Extract interprets any instanceActionResult as an InstanceAction, if possible.
func (r InstanceActionResult) Extract() (InstanceAction, error) {
	var s InstanceAction
	err := r.ExtractInto(&s)
	return s, err
}

func (r InstanceActionResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "instanceAction")
}

func ExtractInstanceActionsInto(r pagination.Page, v interface{}) error {
	return r.(InstanceActionPage).Result.ExtractIntoSlicePtr(v, "instanceActions")
}
