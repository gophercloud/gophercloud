package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/apiversions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const MagnumAPIVersionResponse = `
{
    "versions":[
       {
          "status":"CURRENT",
          "min_version":"1.1",
          "max_version":"1.7",
          "id":"v1",
          "links":[
             {
                "href":"http://10.164.180.104:9511/v1/",
                "rel":"self"
             }
          ]
       }
    ],
    "name":"OpenStack Magnum API",
    "description":"Magnum is an OpenStack project which aims to provide container management."
 }
`

const MagnumAllAPIVersionsResponse = `
{
    "versions":[
       {
          "status":"CURRENT",
          "min_version":"1.1",
          "max_version":"1.7",
          "id":"v1",
          "links":[
             {
                "href":"http://10.164.180.104:9511/v1/",
                "rel":"self"
             }
          ]
       }
    ],
    "name":"OpenStack Magnum API",
    "description":"Magnum is an OpenStack project which aims to provide container management."
 }
`

var MagnumAPIVersion1Result = apiversions.APIVersion{
	ID:         "v1",
	Status:     "CURRENT",
	MinVersion: "1.1",
	Version:    "1.7",
}

var MagnumAllAPIVersionResults = []apiversions.APIVersion{
	MagnumAPIVersion1Result,
}

func MockListResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, MagnumAllAPIVersionsResponse)
	})
}

func MockGetResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, MagnumAPIVersionResponse)
	})
}
