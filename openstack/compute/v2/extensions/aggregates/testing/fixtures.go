package testing

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// AggregateListBody is sample response to the List call
const AggregateListBody = `
{
    "aggregates": [
        {
            "name": "test-aggregate1",
            "availability_zone": null,
            "deleted": false,
            "created_at": "2017-12-22T10:12:06.000000",
            "updated_at": null,
            "hosts": [],
            "deleted_at": null,
            "id": 1,
            "metadata": {}
        },
        {
            "name": "test-aggregate2",
            "availability_zone": "test-az",
            "deleted": false,
            "created_at": "2017-12-22T10:16:07.000000",
            "updated_at": null,
            "hosts": [
                "cmp0"
            ],
            "deleted_at": null,
            "id": 4,
            "metadata": {
                "availability_zone": "test-az"
            }
        }
    ]
}
`

const AggregateCreateBody = `
{
    "aggregate": {
        "availability_zone": "london",
        "created_at": "2016-12-27T22:51:32.000000",
        "deleted": false,
        "deleted_at": null,
        "id": 32,
        "name": "name",
        "updated_at": null
    }
}
`

// First aggregate from the AggregateListBody
var FirstFakeAggregate = aggregates.Aggregate{
	AvailabilityZone: "",
	Hosts:            []string{},
	ID:               1,
	Metadata:         map[string]string{},
	Name:             "test-aggregate1",
	CreatedAt:        time.Date(2017, 12, 22, 10, 12, 6, 0, time.UTC),
	UpdatedAt:        time.Time{},
}

// Second aggregate from the AggregateListBody
var SecondFakeAggregate = aggregates.Aggregate{
	AvailabilityZone: "test-az",
	Hosts:            []string{"cmp0"},
	ID:               4,
	Metadata:         map[string]string{"availability_zone": "test-az"},
	Name:             "test-aggregate2",
	CreatedAt:        time.Date(2017, 12, 22, 10, 16, 7, 0, time.UTC),
	UpdatedAt:        time.Time{},
}

var CreatedAggregate = aggregates.Aggregate{
	AvailabilityZone: "london",
	Hosts:            nil,
	ID:               32,
	Metadata:         nil,
	Name:             "name",
	CreatedAt:        time.Date(2016, 12, 27, 22, 51, 32, 0, time.UTC),
	UpdatedAt:        time.Time{},
}

var AggregateIDtoDelete = 1

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-aggregates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, AggregateListBody)
	})
}

func HandleCreateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-aggregates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, AggregateCreateBody)
	})
}

func HandleDeleteSuccessfully(t *testing.T) {
	v := strconv.Itoa(AggregateIDtoDelete)
	th.Mux.HandleFunc("/os-aggregates/"+v, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusOK)
	})
}
