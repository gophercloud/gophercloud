package defsecrules

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const rootPath = "/os-security-group-default-rules"

func mockListRulesResponse(t *testing.T) {
	th.Mux.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group_default_rules": [
    {
      "from_port": 80,
      "id": "f9a97fcf-3a97-47b0-b76f-919136afb7ed",
      "ip_protocol": "TCP",
      "ip_range": {
        "cidr": "10.10.10.0/24"
      },
      "to_port": 80
    }
  ]
}
      `)
	})
}

func mockCreateRuleResponse(t *testing.T) {
	th.Mux.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "security_group_default_rule": {
    "ip_protocol": "TCP",
    "from_port": 80,
    "to_port": 80,
    "cidr": "10.10.12.0/24"
  }
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group_default_rule": {
    "from_port": 80,
    "id": 1,
    "ip_protocol": "TCP",
    "ip_range": {
      "cidr": "10.10.12.0/24"
    },
    "to_port": 80,
		"id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5"
  }
}
`)
	})
}
