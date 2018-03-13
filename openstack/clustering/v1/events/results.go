package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"time"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// Extract provides access to the individual event returned by Get and extracts Event
func (r commonResult) Extract() (*Event, error) {
	var s struct {
		Event *Event `json:"event"`
	}
	err := r.ExtractInto(&s)
	return s.Event, err
}

// ExtractEvents provides access to the list of events in a page acquired from the ListDetail operation.
func ExtractEvents(r pagination.Page) ([]Event, error) {
	var s struct {
		Events []Event `json:"events"`
	}
	err := (r.(EventPage)).ExtractInto(&s)
	return s.Events, err
}

// EventPage contains a single page of all events from a ListDetails call.
type EventPage struct {
	pagination.LinkedPageBase
}

// Event represents a Detailed Event
type Event struct {
	Action       string    `json:"action,omitempty"`
	ClusterUUID  string    `json:"cluster_id"`
	ID           string    `json:"id"`
	Level        string    `json:"level"`
	OidUUID      string    `json:"oid"`
	OName        string    `json:"oname"`
	OType        string    `json:"otype"`
	ProjectUUID  string    `json:"project"`
	Status       string    `json:"status"`
	StatusReason string    `json:"status_reason"`
	Timestamp    time.Time `json:"-"`
	UserUUID     string    `json:"user"`
}

// IsEmpty determines if a EventPage contains any results.
func (page EventPage) IsEmpty() (bool, error) {
	events, err := ExtractEvents(page)
	return len(events) == 0, err
}

// TODO: Should consolidate this. Should time parsing be strict
type JSONRFC3339Milli time.Time

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
	t, err := time.Parse(gophercloud.RFC3339Milli, s)
	if err != nil {
		t, err = time.Parse(gophercloud.RFC3339MilliNoZ, s)
		if err != nil {
			t, err = time.Parse(gophercloud.RFC3339NoZ, s)
			if err != nil {
				t, err = time.Parse(time.RFC3339, s)
				if err != nil {
					return err
				}
			}
		}
	}
	*jt = JSONRFC3339Milli(t)
	return nil
}

func (r *Event) UnmarshalJSON(b []byte) error {
	type tmp Event
	var s struct {
		tmp
		Timestamp JSONRFC3339Milli `json:"timestamp"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	*r = Event(s.tmp)

	r.Timestamp = time.Time(s.Timestamp)

	return nil
}
