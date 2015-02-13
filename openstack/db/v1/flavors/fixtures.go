package flavors

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func HandleListFlavorsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/flavors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "flavors": [
    {
      "id": 1,
      "links": [
        {
          "href": "https://openstack.example.com/v1.0/1234/flavors/1",
          "rel": "self"
        },
        {
          "href": "https://openstack.example.com/flavors/1",
          "rel": "bookmark"
        }
      ],
      "name": "m1.tiny",
      "ram": 512
    },
    {
      "id": 2,
      "links": [
        {
          "href": "https://openstack.example.com/v1.0/1234/flavors/2",
          "rel": "self"
        },
        {
          "href": "https://openstack.example.com/flavors/2",
          "rel": "bookmark"
        }
      ],
      "name": "m1.small",
      "ram": 1024
    },
    {
      "id": 3,
      "links": [
        {
          "href": "https://openstack.example.com/v1.0/1234/flavors/3",
          "rel": "self"
        },
        {
          "href": "https://openstack.example.com/flavors/3",
          "rel": "bookmark"
        }
      ],
      "name": "m1.medium",
      "ram": 2048
    },
    {
      "id": 4,
      "links": [
        {
          "href": "https://openstack.example.com/v1.0/1234/flavors/4",
          "rel": "self"
        },
        {
          "href": "https://openstack.example.com/flavors/4",
          "rel": "bookmark"
        }
      ],
      "name": "m1.large",
      "ram": 4096
    }
  ]
}
`)
	})
}

func HandleGetFlavorSuccessfully(t *testing.T, flavorID string) {
	th.Mux.HandleFunc("/flavors/"+flavorID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "flavor": {
    "id": 1,
    "links": [
      {
        "href": "https://openstack.example.com/v1.0/1234/flavors/1",
        "rel": "self"
      }
    ],
    "name": "m1.tiny",
    "ram": 512
  }
}
`)
	})
}
