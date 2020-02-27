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

// Event represents an event of instance action.
type Event struct {
	Event string `json:"event"`
	// Host is the host of the event.
	// This requires microversion 2.62 or later.
	Host *string `json:"host"`
	// HostID is the host id of the event.
	// This requires microversion 2.62 or later.
	HostID     *string   `json:"hostId"`
	Result     string    `json:"result"`
	Traceback  string    `json:"traceback"`
	StartTime  time.Time `json:"-"`
	FinishTime time.Time `json:"-"`
}

// UnmarshalJSON converts our JSON API response into our instance action struct
func (e *Event) UnmarshalJSON(b []byte) error {
	type tmp Event
	var s struct {
		tmp
		StartTime  gophercloud.JSONRFC3339MilliNoZ `json:"start_time"`
		FinishTime gophercloud.JSONRFC3339MilliNoZ `json:"finish_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*e = Event(s.tmp)

	e.StartTime = time.Time(s.StartTime)
	e.FinishTime = time.Time(s.FinishTime)

	return err
}

// InstanceActionDetail gives more details on instance action.
type InstanceActionDetail struct {
	Action       string `json:"action"`
	InstanceUUID string `json:"instance_uuid"`
	Message      string `json:"message"`
	ProjectID    string `json:"project_id"`
	RequestID    string `json:"request_id"`
	UserID       string `json:"user_id"`
	// Events is the list of events of the action.
	// This requires microversion 2.50 or later.
	Events *[]Event `json:"events"`
	// UpdatedAt last update date of the action.
	// This requires microversion 2.58 or later.
	UpdatedAt *time.Time `json:"-"`
	StartTime time.Time  `json:"-"`
}

// UnmarshalJSON converts our JSON API response into our instance action struct
func (i *InstanceActionDetail) UnmarshalJSON(b []byte) error {
	type tmp InstanceActionDetail
	var s struct {
		tmp
		UpdatedAt *gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
		StartTime gophercloud.JSONRFC3339MilliNoZ  `json:"start_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*i = InstanceActionDetail(s.tmp)

	i.UpdatedAt = (*time.Time)(s.UpdatedAt)
	i.StartTime = time.Time(s.StartTime)
	return err
}

// InstanceActionResult is the result handler of Get.
type InstanceActionResult struct {
	gophercloud.Result
}

// Extract interprets any instanceActionResult as an InstanceActionDetail, if possible.
func (r InstanceActionResult) Extract() (InstanceActionDetail, error) {
	var s InstanceActionDetail
	err := r.ExtractInto(&s)
	return s, err
}

func (r InstanceActionResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "instanceAction")
}

func ExtractInstanceActionsInto(r pagination.Page, v interface{}) error {
	return r.(InstanceActionPage).Result.ExtractIntoSlicePtr(v, "instanceActions")
}
