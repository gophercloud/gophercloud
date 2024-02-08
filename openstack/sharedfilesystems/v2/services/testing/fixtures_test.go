package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/services"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ServiceListBody is sample response to the List call
const ServiceListBody = `
{
	"services": [
			{
					"status": "enabled",
					"binary": "manila-share",
					"zone": "manila",
					"host": "manila2@generic1",
					"updated_at": "2015-09-07T13:03:57.000000",
					"state": "up",
					"id": 1
			},
			{
					"status": "enabled",
					"binary": "manila-scheduler",
					"zone": "manila",
					"host": "manila2",
					"updated_at": "2015-09-07T13:03:57.000000",
					"state": "up",
					"id": 2
			}
	]
}
`

// First service from the ServiceListBody
var FirstFakeService = services.Service{
	Binary:    "manila-share",
	Host:      "manila2@generic1",
	ID:        1,
	State:     "up",
	Status:    "enabled",
	UpdatedAt: time.Date(2015, 9, 7, 13, 3, 57, 0, time.UTC),
	Zone:      "manila",
}

// Second service from the ServiceListBody
var SecondFakeService = services.Service{
	Binary:    "manila-scheduler",
	Host:      "manila2",
	ID:        2,
	State:     "up",
	Status:    "enabled",
	UpdatedAt: time.Date(2015, 9, 7, 13, 3, 57, 0, time.UTC),
	Zone:      "manila",
}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, ServiceListBody)
	})
}
