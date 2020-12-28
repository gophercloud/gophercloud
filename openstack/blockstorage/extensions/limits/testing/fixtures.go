package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/limits"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
	"limits": {
	  "rate": [],
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
	Rate: []interface{}{},
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
