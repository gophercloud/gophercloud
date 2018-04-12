package actions

import (
	"encoding/json"
	"time"

	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// ActionPage contains a single page of all actions from a ListDetails call.
type ActionPage struct {
	pagination.LinkedPageBase
}

// Action represents a Detailed Action
type Action struct {
	Action       string                 `json:"action"`
	Cause        string                 `json:"cause"`
	CreatedAt    time.Time              `json:"-"`
	Data         map[string]interface{} `json:"data"`
	DependedBy   []string               `json:"depended_by"`
	DependsOn    []string               `json:"depends_on"`
	StartTime    float32                `json:"start_time"`
	EndTime      float32                `json:"end_time"`
	ID           string                 `json:"id"`
	Inputs       map[string]interface{} `json:"inputs"`
	Interval     int                    `json:"interval"`
	Name         string                 `json:"name"`
	Outputs      map[string]interface{} `json:"outputs"`
	Owner        string                 `json:"owner"`
	Project      string                 `json:"project"`
	Status       string                 `json:"status"`
	StatusReason string                 `json:"status_reason"`
	Target       string                 `json:"target"`
	Timeout      int                    `json:"timeout"`
	UpdatedAt    time.Time              `json:"-"`
	User         string                 `json:"user"`
}

// ExtractActions provides access to the list of actions in a page acquired from the List operation.
func ExtractActions(r pagination.Page) ([]Action, error) {
	var s struct {
		Actions []Action `json:"actions"`
	}
	err := (r.(ActionPage)).ExtractInto(&s)
	return s.Actions, err
}

// IsEmpty determines if a ActionPage contains any results.
func (r ActionPage) IsEmpty() (bool, error) {
	actions, err := ExtractActions(r)
	return len(actions) == 0, err
}

type JSONRFC3339MilliAllowEmpty gophercloud.JSONRFC3339Milli

func (jt *JSONRFC3339MilliAllowEmpty) UnmarshalJSON(data []byte) error {
	d := string(data)
	d = strings.Trim(d, "\"' ")
	if len(d) == 0 || d == "null" {
		// Allows empty or "null" field
		*jt = JSONRFC3339MilliAllowEmpty(time.Time{})
		return nil
	}

	gt := gophercloud.JSONRFC3339Milli{}
	err := gt.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339MilliAllowEmpty(gt)
	return err
}
func (r *Action) UnmarshalJSON(b []byte) error {
	type tmp Action
	var s struct {
		tmp
		CreatedAt JSONRFC3339MilliAllowEmpty `json:"created_at"`
		UpdatedAt JSONRFC3339MilliAllowEmpty `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Action(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}
