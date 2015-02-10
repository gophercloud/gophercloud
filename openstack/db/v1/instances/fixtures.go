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
					"href": "https://my-openstack.com/v1.0/1234/flavors/1",
					"rel": "self"
				},
				{
					"href": "https://my-openstack.com/v1.0/1234/flavors/1",
					"rel": "bookmark"
				}
			]
		},
		"links": [
			{
				"href": "https://my-openstack.com/v1.0/1234/instances/1",
				"rel": "self"
			}
		],
		"hostname": "e09ad9a3f73309469cf1f43d11e79549caf9acf2.my-openstack.com",
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
    }
  }
}
`)

		fmt.Fprintf(w, singleInstanceJson)
	})
}

func HandleListInstanceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
{
  "instances": [
    {
      "name": "xml_rack_instance",
      "status": "ACTIVE",
      "volume": {
        "size": 2
      },
      "flavor": {
        "id": "1",
        "links": [
          {
            "href": "https://openstack.example.com/v1.0/1234/flavors/1",
            "rel": "self"
          },
          {
            "href": "https://openstack.example.com/flavors/1",
            "rel": "bookmark"
          }
        ]
      },
      "id": "8fb081af-f237-44f5-80cc-b46be1840ca9",
      "links": [
        {
          "href": "https://openstack.example.com/v1.0/1234/instances/8fb081af-f237-44f5-80cc-b46be1840ca9",
          "rel": "self"
        }
      ]
    }
  ]
}
`)
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

func HandleDeleteInstanceSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}
