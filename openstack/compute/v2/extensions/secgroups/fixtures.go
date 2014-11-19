package secgroups

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func mockListGroupsResponse(t *testing.T) {
	th.Mux.HandleFunc("/os-security-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_groups": [
    {
      "description": "default",
      "id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
      "name": "default",
      "rules": [],
      "tenant_id": "openstack"
    }
  ]
}
`)
	})
}
