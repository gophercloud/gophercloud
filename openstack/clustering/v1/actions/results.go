package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

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
	Action       string                   `json:"action"`
	Cause        string                   `json:"cause"`
	CreatedAt    time.Time                `json:"-"`
	Data         map[string]interface{}   `json:"data"`
	DependedBy   []map[string]interface{} `json:"data"`
	DependedOn   []map[string]interface{} `json:"data"`
	StartTime    float32                  `json:"start_time"`
	EndTime      float32                  `json:"end_time"`
	ID           string                   `json:"id"`
	Inputs       map[string]interface{}   `json:"inputs"`
	Interval     int                      `json:"interval"`
	Name         string                   `json:"name"`
	Outputs      map[string]interface{}   `json:"outputs"`
	Owner        string                   `json:"owner"`
	ProjectUUID  string                   `json:"project"`
	Status       string                   `json:"status"`
	StatusReason string                   `json:"status_reason"`
	Target       string                   `json:"target"`
	Timeout      int                      `json:"timeout"`
	UpdatedAt    time.Time                `json:"-"`
	UserUUID     string                   `json:"user"`
}

func (r commonResult) ExtractAction() (*Action, error) {
	var s struct {
		Action *Action `json:"action"`
	}
	err := r.ExtractInto(&s)
	return s.Action, err
}

// ExtractActions provides access to the list of actions in a page acquired from the ListDetail operation.
func ExtractActions(r pagination.Page) ([]Action, error) {
	var s struct {
		Actions []Action `json:"actions"`
	}
	err := (r.(ActionPage)).ExtractInto(&s)
	return s.Actions, err
}

// IsEmpty determines if a ClusterPage contains any results.
func (page ActionPage) IsEmpty() (bool, error) {
	actions, err := ExtractActions(page)
	return len(actions) == 0, err
}

type JSONRFC3339Milli time.Time

const RFC3339Milli = "2006-01-02T15:04:05.999999Z"

func (jt *JSONRFC3339Milli) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*jt = JSONRFC3339Milli(time.Time{})
		return nil
	}

	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	t, err := time.Parse(RFC3339Milli, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339Milli(t)
	return nil
}

func (r *Action) UnmarshalJSON(b []byte) error {
	type tmp Action
	var s struct {
		tmp
		CreatedAt JSONRFC3339Milli `json:"created_at"`
		UpdatedAt JSONRFC3339Milli `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal")
		fmt.Println("%v", err)
		return err
	}
	*r = Action(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}
