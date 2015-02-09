package stackevents

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Event represents a stack event.
type Event struct {
	ResourceName         string                 `mapstructure:"resource_name"`
	Time                 time.Time              `mapstructure:"-"`
	Links                []gophercloud.Link     `mapstructure:"links"`
	LogicalResourceID    string                 `mapstructure:"logical_resource_id"`
	ResourceStatusReason string                 `mapstructure:"resource_status_reason"`
	ResourceStatus       string                 `mapstructure:"resource_status"`
	PhysicalResourceID   string                 `mapstructure:"physical_resource_id"`
	ID                   string                 `mapstructure:"id"`
	ResourceProperties   map[string]interface{} `mapstructure:"resource_properties"`
}

// FindResult represents the result of a Find operation.
type FindResult struct {
	gophercloud.Result
}

// Extract returns a slice of Event objects and is called after a
// Find operation.
func (r FindResult) Extract() ([]Event, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Res []Event `mapstructure:"events"`
	}

	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	events := r.Body.(map[string]interface{})["events"].([]interface{})

	for i, eventRaw := range events {
		event := eventRaw.(map[string]interface{})
		if date, ok := event["event_time"]; ok && date != nil {
			t, err := time.Parse(time.RFC3339, date.(string))
			if err != nil {
				return nil, err
			}
			res.Res[i].Time = t
		}
	}

	return res.Res, nil
}

// EventPage abstracts the raw results of making a List() request against the API.
// As OpenStack extensions may freely alter the response bodies of structures returned to the client, you may only safely access the
// data provided through the ExtractResources call.
type EventPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a page contains no Server results.
func (r EventPage) IsEmpty() (bool, error) {
	events, err := ExtractEvents(r)
	if err != nil {
		return true, err
	}
	return len(events) == 0, nil
}

// LastMarker returns the last stack ID in a ListResult.
func (r EventPage) LastMarker() (string, error) {
	events, err := ExtractEvents(r)
	if err != nil {
		return "", err
	}
	if len(events) == 0 {
		return "", nil
	}
	return events[len(events)-1].ID, nil
}

// ExtractEvents interprets the results of a single page from a List() call, producing a slice of Event entities.
func ExtractEvents(page pagination.Page) ([]Event, error) {
	casted := page.(EventPage).Body

	var res struct {
		Res []Event `mapstructure:"events"`
	}

	if err := mapstructure.Decode(casted, &res); err != nil {
		return nil, err
	}

	events := casted.(map[string]interface{})["events"].([]interface{})

	for i, eventRaw := range events {
		event := eventRaw.(map[string]interface{})
		if date, ok := event["event_time"]; ok && date != nil {
			t, err := time.Parse(time.RFC3339, date.(string))
			if err != nil {
				return nil, err
			}
			res.Res[i].Time = t
		}
	}

	return res.Res, nil
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract returns a pointer to an Event object and is called after a
// Get operation.
func (r GetResult) Extract() (*Event, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Res *Event `mapstructure:"event"`
	}

	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	event := r.Body.(map[string]interface{})["event"].(map[string]interface{})

	if date, ok := event["event_time"]; ok && date != nil {
		t, err := time.Parse(time.RFC3339, date.(string))
		if err != nil {
			return nil, err
		}
		res.Res.Time = t
	}

	return res.Res, nil
}
