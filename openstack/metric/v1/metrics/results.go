package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
)

// prometheusResponse is the common envelope for all Prometheus HTTP API responses.
type prometheusResponse struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrorType string          `json:"errorType"`
	Error     string          `json:"error"`
}

// checkResponse unmarshals the Prometheus envelope and returns an error if
// the response status is "error".
func checkResponse(body any) (*prometheusResponse, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var resp prometheusResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	if resp.Status == "error" {
		return &resp, fmt.Errorf("prometheus %s: %s", resp.ErrorType, resp.Error)
	}
	return &resp, nil
}

// QueryResult is the result of a Query request.
type QueryResult struct {
	gophercloud.Result
}

// Extract interprets a QueryResult as QueryData.
func (r QueryResult) Extract() (*QueryData, error) {
	resp, err := checkResponse(r.Body)
	if err != nil {
		if r.Err != nil {
			return nil, r.Err
		}
		return nil, err
	}
	var data QueryData
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// LabelsResult is the result of a Labels request.
type LabelsResult struct {
	gophercloud.Result
}

// Extract interprets a LabelsResult as a slice of label name strings.
func (r LabelsResult) Extract() ([]string, error) {
	resp, err := checkResponse(r.Body)
	if err != nil {
		if r.Err != nil {
			return nil, r.Err
		}
		return nil, err
	}
	var data []string
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// LabelValuesResult is the result of a LabelValues request.
type LabelValuesResult struct {
	gophercloud.Result
}

// Extract interprets a LabelValuesResult as a slice of label value strings.
func (r LabelValuesResult) Extract() ([]string, error) {
	resp, err := checkResponse(r.Body)
	if err != nil {
		if r.Err != nil {
			return nil, r.Err
		}
		return nil, err
	}
	var data []string
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// SeriesResult is the result of a Series request.
type SeriesResult struct {
	gophercloud.Result
}

// Extract interprets a SeriesResult as a slice of label-set maps.
func (r SeriesResult) Extract() ([]map[string]string, error) {
	resp, err := checkResponse(r.Body)
	if err != nil {
		if r.Err != nil {
			return nil, r.Err
		}
		return nil, err
	}
	var data []map[string]string
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// TargetsResult is the result of a Targets request.
type TargetsResult struct {
	gophercloud.Result
}

// Extract interprets a TargetsResult as TargetsData.
func (r TargetsResult) Extract() (*TargetsData, error) {
	resp, err := checkResponse(r.Body)
	if err != nil {
		if r.Err != nil {
			return nil, r.Err
		}
		return nil, err
	}
	var data TargetsData
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// RuntimeInfoResult is the result of a RuntimeInfo request.
type RuntimeInfoResult struct {
	gophercloud.Result
}

// Extract interprets a RuntimeInfoResult as RuntimeInfoData.
func (r RuntimeInfoResult) Extract() (*RuntimeInfoData, error) {
	resp, err := checkResponse(r.Body)
	if err != nil {
		if r.Err != nil {
			return nil, r.Err
		}
		return nil, err
	}
	var data RuntimeInfoData
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// CleanTombstonesResult is the result of a CleanTombstones request.
type CleanTombstonesResult struct {
	gophercloud.ErrResult
}

// DeleteSeriesResult is the result of a DeleteSeries request.
type DeleteSeriesResult struct {
	gophercloud.ErrResult
}

// SnapshotResult is the result of a Snapshot request.
type SnapshotResult struct {
	gophercloud.Result
}

// Extract interprets a SnapshotResult as SnapshotData.
func (r SnapshotResult) Extract() (*SnapshotData, error) {
	resp, err := checkResponse(r.Body)
	if err != nil {
		if r.Err != nil {
			return nil, r.Err
		}
		return nil, err
	}
	var data SnapshotData
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// QueryData represents the data field of a Prometheus instant query response.
type QueryData struct {
	// ResultType is the type of result (e.g., "vector", "matrix", "scalar", "string").
	ResultType string `json:"resultType"`

	// Result is the list of metric values returned by the query.
	Result []MetricValue `json:"result"`
}

// MetricValue represents a single metric result with its labels and value.
type MetricValue struct {
	// Metric contains the label set identifying this metric.
	Metric map[string]string `json:"metric"`

	// Value contains a single sample value as [timestamp, "string_value"].
	Value []any `json:"value"`
}

// TargetsData represents the data field of a Prometheus targets response.
type TargetsData struct {
	// ActiveTargets is the list of currently active scrape targets.
	ActiveTargets []ActiveTarget `json:"activeTargets"`

	// DroppedTargets is the list of targets that were dropped during relabeling.
	DroppedTargets []DroppedTarget `json:"droppedTargets"`
}

// ActiveTarget represents an active Prometheus scrape target.
type ActiveTarget struct {
	// DiscoveredLabels are the unmodified labels from service discovery.
	DiscoveredLabels map[string]string `json:"discoveredLabels"`

	// Labels are the labels after relabeling.
	Labels map[string]string `json:"labels"`

	// ScrapePool is the name of the scrape configuration.
	ScrapePool string `json:"scrapePool"`

	// ScrapeURL is the URL being scraped.
	ScrapeURL string `json:"scrapeUrl"`

	// GlobalURL is the URL as available from other Prometheus instances.
	GlobalURL string `json:"globalUrl"`

	// LastError is the last error encountered during scraping.
	LastError string `json:"lastError"`

	// LastScrape is the time of the last scrape.
	LastScrape string `json:"lastScrape"`

	// LastScrapeDuration is the duration of the last scrape in seconds.
	LastScrapeDuration float64 `json:"lastScrapeDuration"`

	// Health is the health status of the target ("up", "down", "unknown").
	Health string `json:"health"`

	// ScrapeInterval is the configured scrape interval.
	ScrapeInterval string `json:"scrapeInterval"`

	// ScrapeTimeout is the configured scrape timeout.
	ScrapeTimeout string `json:"scrapeTimeout"`
}

// DroppedTarget represents a target that was dropped during relabeling.
type DroppedTarget struct {
	// DiscoveredLabels are the unmodified labels from service discovery.
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
}

// RuntimeInfoData represents the data field of a Prometheus runtime info response.
type RuntimeInfoData struct {
	// StartTime is when the Prometheus server started.
	StartTime string `json:"startTime"`

	// CWD is the current working directory of the Prometheus server.
	CWD string `json:"CWD"`

	// ReloadConfigSuccess indicates if the last config reload was successful.
	ReloadConfigSuccess bool `json:"reloadConfigSuccess"`

	// LastConfigTime is the time of the last config reload.
	LastConfigTime string `json:"lastConfigTime"`

	// CorruptionCount is the number of WAL corruptions encountered.
	CorruptionCount int `json:"corruptionCount"`

	// GoroutineCount is the number of goroutines.
	GoroutineCount int `json:"goroutineCount"`

	// GOMAXPROCS is the configured GOMAXPROCS value.
	GOMAXPROCS int `json:"GOMAXPROCS"`

	// GOGC is the configured GOGC value.
	GOGC string `json:"GOGC"`

	// GODEBUG is the configured GODEBUG value.
	GODEBUG string `json:"GODEBUG"`

	// StorageRetention is the configured storage retention.
	StorageRetention string `json:"storageRetention"`
}

// SnapshotData represents the data field of a Prometheus snapshot response.
type SnapshotData struct {
	// Name is the filename of the snapshot created.
	Name string `json:"name"`
}
