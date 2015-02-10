package instances

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const singleInstanceJson = `
{
  "instance": {
    "created": "2014-02-13T21:47:13",
    "datastore": {
      "type": "mysql",
      "version": "5.6"
    },
    "flavor": {
      "id": "1",
      "links": [
        {
          "href": "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1",
          "rel": "self"
        },
        {
          "href": "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1",
          "rel": "bookmark"
        }
      ]
    },
    "links": [
      {
        "href": "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1",
        "rel": "self"
      }
    ],
    "hostname": "e09ad9a3f73309469cf1f43d11e79549caf9acf2.rackspaceclouddb.com",
    "id": "d4603f69-ec7e-4e9b-803f-600b9205576f",
    "name": "json_rack_instance",
    "status": "BUILD",
    "updated": "2014-02-13T21:47:13",
    "volume": {
      "size": 2
    }
  }
}
`

func HandleCreateInstanceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		th.TestJSONRequest(t, r, `
{
  "instance": {
    "databases": [
      {
        "character_set": "utf8",
        "collate": "utf8_general_ci",
        "name": "sampledb"
      },
      {
        "name": "nextround"
      }
    ],
    "flavorRef": "1",
    "name": "json_rack_instance",
    "users": [
      {
        "databases": [
          {
            "name": "sampledb"
          }
        ],
        "name": "demouser",
        "password": "demopassword"
      }
    ],
    "volume": {
      "size": 2
    },
		"restorePoint": "1234567890"
  }
}
`)

		fmt.Fprintf(w, singleInstanceJson)
	})
}

func HandleGetInstanceSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, singleInstanceJson)
	})
}
