package testing

import (
	"fmt"
	"net/http"
	"testing"

	az "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/availabilityzones"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const GetOutput = `
{
    "availabilityZoneInfo": [
        {
            "zoneName": "internal",
            "zoneState": {
                "available": true
            }
        },
        {
            "zoneName": "nova",
            "zoneState": {
                "available": true
            }
        }
    ]
}`

var AZResult = []az.AvailabilityZone{
	{
		ZoneName:  "internal",
		ZoneState: az.ZoneState{Available: true},
	},
	{
		ZoneName:  "nova",
		ZoneState: az.ZoneState{Available: true},
	},
}

// HandleGetSuccessfully configures the test server to respond to a Get request
// for availability zone information.
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-availability-zone", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, GetOutput)
	})
}
