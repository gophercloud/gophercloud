package users

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func HandleCreateUserSuccessfully(t *testing.T, instanceID string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `
{
  "users": [
    {
      "databases": [
        {
          "name": "databaseA"
        }
      ],
      "name": "dbuser3",
      "password": "secretsecret"
    },
    {
      "databases": [
        {
          "name": "databaseB"
        },
        {
          "name": "databaseC"
        }
      ],
      "name": "dbuser4",
      "password": "secretsecret"
    }
  ]
}
`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleListUsersSuccessfully(t *testing.T, instanceID string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
	"users": [
		{
			"databases": [
				{
					"name": "databaseA"
				}
			],
			"name": "dbuser3"
		},
		{
			"databases": [
				{
					"name": "databaseB"
				},
				{
					"name": "databaseC"
				}
			],
			"name": "dbuser4"
		}
	]
}
`)
	})
}

func HandleDeleteUserSuccessfully(t *testing.T, instanceID, dbName string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/users/"+dbName, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}
