package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/apiversions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const NovaAPIVersionResponse_20 = `
{
    "version": {
        "id": "v2.0",
        "status": "SUPPORTED",
        "version": "",
        "min_version": "",
        "updated": "2011-01-21T11:33:21Z",
        "links": [
            {
                "rel": "self",
                "href": "http://10.1.5.216/compute/v2/"
            },
            {
                "rel": "describedby",
                "type": "text/html",
                "href": "http://docs.openstack.org/"
            }
        ],
        "media-types": [
            {
                "base": "application/json",
                "type": "application/vnd.openstack.compute+json;version=2"
            }
        ]
    }
}
`

const NovaAPIVersionResponse_21 = `
{
    "version": {
        "id": "v2.1",
        "status": "CURRENT",
        "version": "2.87",
        "min_version": "2.1",
        "updated": "2013-07-23T11:33:21Z",
        "links": [
            {
                "rel": "self",
                "href": "http://10.1.5.216/compute/v2.1/"
            },
            {
                "rel": "describedby",
                "type": "text/html",
                "href": "http://docs.openstack.org/"
            }
        ],
        "media-types": [
            {
                "base": "application/json",
                "type": "application/vnd.openstack.compute+json;version=2.1"
            }
        ]
    }
}

`

const NovaAPIInvalidVersionResponse = `
{
    "choices": [
        {
            "id": "v2.0",
            "status": "SUPPORTED",
            "links": [
                {
                    "rel": "self",
                    "href": "http://10.1.5.216/compute/v2/compute/v3"
                }
            ],
            "media-types": [
                {
                    "base": "application/json",
                    "type": "application/vnd.openstack.compute+json;version=2"
                }
            ]
        },
        {
            "id": "v2.1",
            "status": "CURRENT",
            "links": [
                {
                    "rel": "self",
                    "href": "http://10.1.5.216/compute/v2.1/compute/v3"
                }
            ],
            "media-types": [
                {
                    "base": "application/json",
                    "type": "application/vnd.openstack.compute+json;version=2.1"
                }
            ]
        }
    ]
}
`

const NovaAllAPIVersionsResponse = `
{
    "versions": [
        {
            "id": "v2.0",
            "status": "SUPPORTED",
            "version": "",
            "min_version": "",
            "updated": "2011-01-21T11:33:21Z",
            "links": [
                {
                    "rel": "self",
                    "href": "http://10.1.5.216/compute/v2/"
                }
            ]
        },
        {
            "id": "v2.1",
            "status": "CURRENT",
            "version": "2.87",
            "min_version": "2.1",
            "updated": "2013-07-23T11:33:21Z",
            "links": [
                {
                    "rel": "self",
                    "href": "http://10.1.5.216/compute/v2.1/"
                }
            ]
        }
    ]
}
`

var NovaAPIVersion20Result = apiversions.APIVersion{
	ID:      "v2.0",
	Status:  "SUPPORTED",
	Updated: time.Date(2011, 1, 21, 11, 33, 21, 0, time.UTC),
}

var NovaAPIVersion21Result = apiversions.APIVersion{
	ID:         "v2.1",
	Status:     "CURRENT",
	Updated:    time.Date(2013, 7, 23, 11, 33, 21, 0, time.UTC),
	MinVersion: "2.1",
	Version:    "2.87",
}

var NovaAllAPIVersionResults = []apiversions.APIVersion{
	NovaAPIVersion20Result,
	NovaAPIVersion21Result,
}

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, NovaAllAPIVersionsResponse)
	})
}

func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc("/v2.1/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, NovaAPIVersionResponse_21)
	})
}

func MockGetMultipleResponses(t *testing.T) {
	th.Mux.HandleFunc("/v3/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, NovaAPIInvalidVersionResponse)
	})
}
