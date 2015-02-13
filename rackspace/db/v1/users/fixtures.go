package users

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const singleDB = `{"databases": [{"name": "databaseE"}]}`

func setupHandler(t *testing.T, url, method, requestBody, responseBody string, status int) {
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, method)
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		if requestBody != "" {
			th.TestJSONRequest(t, r, requestBody)
		}

		w.WriteHeader(status)

		if responseBody != "" {
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, responseBody)
		}
	})
}

func HandleChangePasswordSuccessfully(t *testing.T, instanceID string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `
{
  "users": [
    {
      "name": "dbuser1",
      "password": "newpassword"
    },
    {
      "name": "dbuser2",
      "password": "anotherpassword"
    }
  ]
}
`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleUpdateSuccessfully(t *testing.T, instanceID, userName string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/users/"+userName, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `
{
  "user": {
    "name": "new_username",
    "password": "new_password"
  }
}
`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleGetSuccessfully(t *testing.T, instanceID, userName string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/users/"+userName, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
{
  "user": {
    "name": "exampleuser",
    "host": "foo",
    "databases": [
      {
        "name": "databaseA"
      },
      {
        "name": "databaseB"
      }
    ]
  }
}
`)
	})
}

func HandleListUserAccessSuccessfully(t *testing.T, instanceID, userName string) {
	th.Mux.HandleFunc("/instances/"+instanceID+"/users/"+userName+"/databases", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, singleDB)
	})
}

func HandleGrantUserAccessSuccessfully(t *testing.T, instanceID, userName string) {
	url := "/instances/" + instanceID + "/users/" + userName + "/databases"
	setupHandler(t, url, "PUT", singleDB, "", http.StatusAccepted)
}

func HandleRevokeUserAccessSuccessfully(t *testing.T, instanceID, userName, dbName string) {
	url := "/instances/" + instanceID + "/users/" + userName + "/databases/" + dbName
	setupHandler(t, url, "DELETE", "", "", http.StatusAccepted)
}
