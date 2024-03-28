package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/limits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
	"limits": {
		"rate": [
			{
				"regex": ".*",
				"uri": "*",
				"limit": [
					{
						"verb": "GET",
						"next-available": "1970-01-01T00:00:00",
						"unit": "MINUTE",
						"value": 10,
						"remaining": 10
					},
					{
						"verb": "POST",
						"next-available": "1970-01-01T00:00:00",
						"unit": "HOUR",
						"value": 5,
						"remaining": 5
					}
				]
			},
			{
				"regex": "changes-since",
				"uri": "changes-since*",
				"limit": [
					{
						"verb": "GET",
						"next-available": "1970-01-01T00:00:00",
						"unit": "MINUTE",
						"value": 5,
						"remaining": 5
					}
				]
			}
		],
		"absolute": {
			"maxTotalVolumes": 40,
			"maxTotalSnapshots": 40,
			"maxTotalVolumeGigabytes": 1000,
			"maxTotalBackups": 10,
			"maxTotalBackupGigabytes": 1000,
			"totalVolumesUsed": 1,
			"totalGigabytesUsed": 100,
			"totalSnapshotsUsed": 1,
			"totalBackupsUsed": 1,
			"totalBackupGigabytesUsed": 50
		}
	}
}
`

// LimitsResult is the result of the limits in GetOutput.
var LimitsResult = limits.Limits{
	Rate: []limits.Rate{
		{
			Regex: ".*",
			URI:   "*",
			Limit: []limits.Limit{
				{
					Verb:          "GET",
					NextAvailable: "1970-01-01T00:00:00",
					Unit:          "MINUTE",
					Value:         10,
					Remaining:     10,
				},
				{
					Verb:          "POST",
					NextAvailable: "1970-01-01T00:00:00",
					Unit:          "HOUR",
					Value:         5,
					Remaining:     5,
				},
			},
		},
		{
			Regex: "changes-since",
			URI:   "changes-since*",
			Limit: []limits.Limit{
				{
					Verb:          "GET",
					NextAvailable: "1970-01-01T00:00:00",
					Unit:          "MINUTE",
					Value:         5,
					Remaining:     5,
				},
			},
		},
	},
	Absolute: limits.Absolute{
		MaxTotalVolumes:          40,
		MaxTotalSnapshots:        40,
		MaxTotalVolumeGigabytes:  1000,
		MaxTotalBackups:          10,
		MaxTotalBackupGigabytes:  1000,
		TotalVolumesUsed:         1,
		TotalGigabytesUsed:       100,
		TotalSnapshotsUsed:       1,
		TotalBackupsUsed:         1,
		TotalBackupGigabytesUsed: 50,
	},
}

// HandleGetSuccessfully configures the test server to respond to a Get request
// for a limit.
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/limits", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, GetOutput)
	})
}
