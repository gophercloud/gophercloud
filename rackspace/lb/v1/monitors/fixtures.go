package nodes

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func _rootURL(lbID int) string {
	return "/loadbalancers/" + strconv.Itoa(lbID) + "/healthmonitor"
}

func mockGetResponse(t *testing.T, lbID, nodeID int) {
	th.Mux.HandleFunc(_nodeURL(lbID, nodeID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "healthMonitor": {
    "type": "CONNECT",
    "delay": 10,
    "timeout": 10,
    "attemptsBeforeDeactivation": 3
  }
}
  `)
	})
}

func mockUpdateResponse(t *testing.T, lbID, nodeID int) {
	th.Mux.HandleFunc(_nodeURL(lbID, nodeID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "healthMonitor": {
    "type": "CONNECT",
    "delay": 10,
    "timeout": 10,
    "attemptsBeforeDeactivation": 3
  }
}
    `)

		w.WriteHeader(http.StatusOK)
	})
}

func mockDeleteResponse(t *testing.T, lbID, nodeID int) {
	th.Mux.HandleFunc(_nodeURL(lbID, nodeID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
	})
}
