package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/metric/v1/metrics"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// serviceClient returns a ServiceClient with ResourceBase set to include api/v1/.
func serviceClient(fakeServer th.FakeServer) *gophercloud.ServiceClient {
	sc := client.ServiceClient(fakeServer)
	sc.ResourceBase = sc.Endpoint + "api/v1/"
	return sc
}

// QueryOutput is a sample response to a Query call.
const QueryOutput = `
{
    "status": "success",
    "data": {
        "resultType": "vector",
        "result": [
            {
                "metric": {
                    "__name__": "up",
                    "instance": "localhost:9090",
                    "job": "prometheus"
                },
                "value": [1435781451.781, "1"]
            }
        ]
    }
}
`

// ExpectedQueryData is the expected result from QueryOutput.
var ExpectedQueryData = &metrics.QueryData{
	ResultType: "vector",
	Result: []metrics.MetricValue{
		{
			Metric: map[string]string{
				"__name__": "up",
				"instance": "localhost:9090",
				"job":      "prometheus",
			},
			Value: []any{1435781451.781, "1"},
		},
	},
}

// LabelsOutput is a sample response to a Labels call.
const LabelsOutput = `
{
    "status": "success",
    "data": ["__name__", "instance", "job"]
}
`

// ExpectedLabels is the expected result from LabelsOutput.
var ExpectedLabels = []string{"__name__", "instance", "job"}

// LabelValuesOutput is a sample response to a LabelValues call.
const LabelValuesOutput = `
{
    "status": "success",
    "data": ["node", "prometheus"]
}
`

// ExpectedLabelValues is the expected result from LabelValuesOutput.
var ExpectedLabelValues = []string{"node", "prometheus"}

// SeriesOutput is a sample response to a Series call.
const SeriesOutput = `
{
    "status": "success",
    "data": [
        {
            "__name__": "up",
            "instance": "localhost:9090",
            "job": "prometheus"
        },
        {
            "__name__": "up",
            "instance": "localhost:9100",
            "job": "node"
        }
    ]
}
`

// ExpectedSeries is the expected result from SeriesOutput.
var ExpectedSeries = []map[string]string{
	{
		"__name__": "up",
		"instance": "localhost:9090",
		"job":      "prometheus",
	},
	{
		"__name__": "up",
		"instance": "localhost:9100",
		"job":      "node",
	},
}

// TargetsOutput is a sample response to a Targets call.
const TargetsOutput = `
{
    "status": "success",
    "data": {
        "activeTargets": [
            {
                "discoveredLabels": {
                    "__address__": "localhost:9090",
                    "__scheme__": "http",
                    "job": "prometheus"
                },
                "labels": {
                    "instance": "localhost:9090",
                    "job": "prometheus"
                },
                "scrapePool": "prometheus",
                "scrapeUrl": "http://localhost:9090/metrics",
                "globalUrl": "http://localhost:9090/metrics",
                "lastError": "",
                "lastScrape": "2025-08-20T10:30:00.000Z",
                "lastScrapeDuration": 0.003145,
                "health": "up",
                "scrapeInterval": "15s",
                "scrapeTimeout": "10s"
            }
        ],
        "droppedTargets": [
            {
                "discoveredLabels": {
                    "__address__": "localhost:9091",
                    "job": "dropped"
                }
            }
        ]
    }
}
`

// ExpectedTargetsData is the expected result from TargetsOutput.
var ExpectedTargetsData = &metrics.TargetsData{
	ActiveTargets: []metrics.ActiveTarget{
		{
			DiscoveredLabels: map[string]string{
				"__address__": "localhost:9090",
				"__scheme__":  "http",
				"job":         "prometheus",
			},
			Labels: map[string]string{
				"instance": "localhost:9090",
				"job":      "prometheus",
			},
			ScrapePool:         "prometheus",
			ScrapeURL:          "http://localhost:9090/metrics",
			GlobalURL:          "http://localhost:9090/metrics",
			LastError:          "",
			LastScrape:         "2025-08-20T10:30:00.000Z",
			LastScrapeDuration: 0.003145,
			Health:             "up",
			ScrapeInterval:     "15s",
			ScrapeTimeout:      "10s",
		},
	},
	DroppedTargets: []metrics.DroppedTarget{
		{
			DiscoveredLabels: map[string]string{
				"__address__": "localhost:9091",
				"job":         "dropped",
			},
		},
	},
}

// RuntimeInfoOutput is a sample response to a RuntimeInfo call.
const RuntimeInfoOutput = `
{
    "status": "success",
    "data": {
        "startTime": "2025-08-20T10:00:00.000Z",
        "CWD": "/prometheus",
        "reloadConfigSuccess": true,
        "lastConfigTime": "2025-08-20T10:00:00.000Z",
        "corruptionCount": 0,
        "goroutineCount": 42,
        "GOMAXPROCS": 4,
        "GOGC": "",
        "GODEBUG": "",
        "storageRetention": "15d"
    }
}
`

// ExpectedRuntimeInfoData is the expected result from RuntimeInfoOutput.
var ExpectedRuntimeInfoData = &metrics.RuntimeInfoData{
	StartTime:           "2025-08-20T10:00:00.000Z",
	CWD:                 "/prometheus",
	ReloadConfigSuccess: true,
	LastConfigTime:      "2025-08-20T10:00:00.000Z",
	CorruptionCount:     0,
	GoroutineCount:      42,
	GOMAXPROCS:          4,
	GOGC:                "",
	GODEBUG:             "",
	StorageRetention:    "15d",
}

// SnapshotOutput is a sample response to a Snapshot call.
const SnapshotOutput = `
{
    "status": "success",
    "data": {
        "name": "20250820T100000Z-abcdef1234567890"
    }
}
`

// ExpectedSnapshotData is the expected result from SnapshotOutput.
var ExpectedSnapshotData = &metrics.SnapshotData{
	Name: "20250820T100000Z-abcdef1234567890",
}

// HandleQuerySuccessfully configures the test server to respond to a Query request.
func HandleQuerySuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/query", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, QueryOutput)
	})
}

// HandleLabelsSuccessfully configures the test server to respond to a Labels request.
func HandleLabelsSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/labels", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, LabelsOutput)
	})
}

// HandleLabelValuesSuccessfully configures the test server to respond to a LabelValues request.
func HandleLabelValuesSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/label/job/values", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, LabelValuesOutput)
	})
}

// HandleSeriesSuccessfully configures the test server to respond to a Series request.
func HandleSeriesSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/series", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, SeriesOutput)
	})
}

// HandleTargetsSuccessfully configures the test server to respond to a Targets request.
func HandleTargetsSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/targets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, TargetsOutput)
	})
}

// HandleRuntimeInfoSuccessfully configures the test server to respond to a RuntimeInfo request.
func HandleRuntimeInfoSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/status/runtimeinfo", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, RuntimeInfoOutput)
	})
}

// HandleCleanTombstonesSuccessfully configures the test server to respond to a CleanTombstones request.
func HandleCleanTombstonesSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/admin/tsdb/clean_tombstones", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleDeleteSeriesSuccessfully configures the test server to respond to a DeleteSeries request.
func HandleDeleteSeriesSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/admin/tsdb/delete_series", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleSnapshotSuccessfully configures the test server to respond to a Snapshot request.
func HandleSnapshotSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/api/v1/admin/tsdb/snapshot", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, SnapshotOutput)
	})
}
