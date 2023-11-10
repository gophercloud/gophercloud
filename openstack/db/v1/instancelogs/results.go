package instancelogs

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Log records represent an instance Log API resource.
type Log struct {
	// The name of the log.
	Name string `json:"name"`

	// The type of the log.
	Type string `json:"type"`

	// The log status.
	Status string `json:"status"`

	// Published size of the log.
	Published float64 `json:"published"`

	// Log file size pending to be published.
	Pending float64 `json:"pending"`

	// The object store container where the published log data will be stored.
	// Defaults to None before the log has been published.
	Container string `json:"container"`

	// If the log has been published, this is the prefix location of where the log data are stored.
	// Otherwize the prefix is None.
	Prefix string `json:"prefix"`

	// The log metafile location
	Metafile string `json:"metafile"`
}

type commonResult struct {
	gophercloud.Result
}

// ActionResult represents the result of an Action operation with instance log.
type ActionResult struct {
	commonResult
}

// Extract provides access to the individual Log returned by the Get function.
func (r commonResult) Extract() (*Log, error) {
	var s struct {
		Log *Log `json:"log"`
	}
	err := r.ExtractInto(&s)
	return s.Log, err
}

// LogPage contains a single page of the response from a List call.
type LogPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a page contains any results.
func (page LogPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	logs, err := ExtractLogs(page)
	return len(logs) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (page LogPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"logs_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractLogs provides access to the list of logs in a page acquired from the List operation.
func ExtractLogs(r pagination.Page) ([]Log, error) {
	var s struct {
		Logs []Log `json:"logs"`
	}
	err := (r.(LogPage)).ExtractInto(&s)
	return s.Logs, err
}
