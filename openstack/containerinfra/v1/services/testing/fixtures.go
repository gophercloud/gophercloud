package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/services"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const ServiceListBody = `
{
    "mservices": [
        {
            "host": "magnum-conductor-2",
            "binary": "magnum-conductor",
            "state": "up",
            "id": 25,
            "report_count": 499402,
            "disabled": false,
            "disabled_reason": null,
            "created_at": "2020-11-16T23:35:25+00:00",
            "updated_at": "2021-12-04T09:28:02+00:00"
        },
        {
            "host": "magnum-conductor-0",
            "binary": "magnum-conductor",
            "state": "up",
            "id": 28,
            "report_count": 498163,
            "disabled": false,
            "disabled_reason": null,
            "created_at": "2020-11-16T23:35:32+00:00",
            "updated_at": "2021-12-04T09:28:40+00:00"
        },
        {
            "host": "magnum-conductor-1",
            "binary": "magnum-conductor",
            "state": "up",
            "id": 31,
            "report_count": 499668,
            "disabled": false,
            "disabled_reason": null,
            "created_at": "2020-11-16T23:35:33+00:00",
            "updated_at": "2021-12-04T09:28:00+00:00"
        }
    ]
}
`

var firstCreated, _ = time.Parse(time.RFC3339, "2020-11-16T23:35:25+00:00")
var firstUpdated, _ = time.Parse(time.RFC3339, "2021-12-04T09:28:02+00:00")
var FirstFakeService = services.Service{
	Binary:         "magnum-conductor",
	Host:           "magnum-conductor-2",
	State:          "up",
	ID:             25,
	ReportCount:    499402,
	Disabled:       false,
	DisabledReason: "",
	CreatedAt:      firstCreated,
	UpdatedAt:      firstUpdated,
}

var secondCreated, _ = time.Parse(time.RFC3339, "2020-11-16T23:35:32+00:00")
var secondUpdated, _ = time.Parse(time.RFC3339, "2021-12-04T09:28:40+00:00")
var SecondFakeService = services.Service{
	Host:           "magnum-conductor-0",
	Binary:         "magnum-conductor",
	State:          "up",
	ID:             28,
	ReportCount:    498163,
	Disabled:       false,
	DisabledReason: "",
	CreatedAt:      secondCreated,
	UpdatedAt:      secondUpdated,
}

var thirdCreated, _ = time.Parse(time.RFC3339, "2020-11-16T23:35:33+00:00")
var thirdUpdated, _ = time.Parse(time.RFC3339, "2021-12-04T09:28:00+00:00")
var ThirdFakeService = services.Service{
	Host:           "magnum-conductor-1",
	Binary:         "magnum-conductor",
	State:          "up",
	ID:             31,
	ReportCount:    499668,
	Disabled:       false,
	DisabledReason: "",
	CreatedAt:      thirdCreated,
	UpdatedAt:      thirdUpdated,
}

func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/mservices", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, ServiceListBody)
	})
}
