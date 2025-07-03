package testing

import (
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

var (
	projectID = "{projectID}"
	resURL    = "/mgmt/" + "quotas/" + projectID
)

// getQuotasResp is a sample response to a Get call.
var getQuotasResp = `
{
    "quotas": [
        {
            "in_use": 5,
            "limit": 15,
            "reserved": 0,
            "resource": "instances"
        },
        {
            "in_use": 2,
            "limit": 50,
            "reserved": 0,
            "resource": "backups"
        },
        {
            "in_use": 1,
            "limit": 40,
            "reserved": 0,
            "resource": "volumes"
        }
    ]
}
`

func HandleGet(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL, "GET", "", getQuotasResp, 200)
}
