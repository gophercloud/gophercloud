package databases

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func HandleCreateDBSuccessfully(t *testing.T, instanceID string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/databases", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `
{
  "databases": [
    {
      "character_set": "utf8",
      "collate": "utf8_general_ci",
      "name": "testingdb"
    },
    {
      "name": "sampledb"
    }
  ]
}
`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleListDBsSuccessfully(t *testing.T, instanceID string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/databases", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
	"databases": [
		{
			"name": "anotherexampledb"
		},
		{
			"name": "exampledb"
		},
		{
			"name": "nextround"
		},
		{
			"name": "sampledb"
		},
		{
			"name": "testingdb"
		}
	]
}
`)
	})
}

func HandleDeleteDBSuccessfully(t *testing.T, instanceID, dbName string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/databases/"+dbName, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}
