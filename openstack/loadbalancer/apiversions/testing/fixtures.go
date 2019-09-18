package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/apiversions"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const OctaviaAllAPIVersionsResponse = `
{
    "versions": [
        {
            "id": "v2.9",
            "links": [
                {
                    "href": "http://192.168.206.8/load-balancer//v2",
                    "rel": "self"
                }
            ],
            "status": "SUPPORTED",
            "updated": "2019-03-04T00:00:00Z"
        },
        {
            "id": "v2.10",
            "links": [
                {
                    "href": "http://192.168.206.8/load-balancer//v2",
                    "rel": "self"
                }
            ],
            "status": "CURRENT",
            "updated": "2019-03-05T00:00:00Z"
        }
    ]
}
`

var APIVersion29 = apiversions.APIVersion{
	ID:      "v2.9",
	Status:  "SUPPORTED",
	Updated: "2019-03-04T00:00:00Z",
}

var APIVersion210 = apiversions.APIVersion{
	ID:      "v2.10",
	Status:  "CURRENT",
	Updated: "2019-03-05T00:00:00Z",
}

var APIVersions = []apiversions.APIVersion{
	APIVersion29,
	APIVersion210,
}

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, OctaviaAllAPIVersionsResponse)
	})
}
