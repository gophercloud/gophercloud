package secgroups

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const rootPath = "/os-security-groups"

const listGroupsJSON = `
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
`

func mockListGroupsResponse(t *testing.T) {
	th.Mux.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, listGroupsJSON)
	})
}

func mockListGroupsByServerResponse(t *testing.T, serverID string) {
	url := fmt.Sprintf("%s/servers/%s%s", rootPath, serverID, rootPath)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, listGroupsJSON)
	})
}

func mockCreateGroupResponse(t *testing.T) {
	th.Mux.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "security_group": {
    "name": "test",
    "description": "something"
  }
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group": {
    "description": "something",
    "id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
    "name": "test",
    "rules": [],
    "tenant_id": "openstack"
  }
}
`)
	})
}

func mockUpdateGroupResponse(t *testing.T, groupID string) {
	url := fmt.Sprintf("%s/%s", rootPath, groupID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "security_group": {
    "name": "new_name"
  }
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group": {
    "description": "something",
    "id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
    "name": "new_name",
    "rules": [],
    "tenant_id": "openstack"
  }
}
`)
	})
}

func mockGetGroupsResponse(t *testing.T, groupID string) {
	url := fmt.Sprintf("%s/%s", rootPath, groupID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group": {
    "description": "default",
    "id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
    "name": "default",
    "rules": [
      {
        "from_port": 80,
        "group": {
          "tenant_id": "openstack",
          "name": "default"
        },
        "ip_protocol": "TCP",
        "to_port": 85,
        "parent_group_id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
        "ip_range": {
						"cidr": "0.0.0.0"
				},
        "id": "ebe599e2-6b8c-457c-b1ff-a75e48f10923"
      }
    ],
    "tenant_id": "openstack"
  }
}
			`)
	})
}

func mockDeleteGroupResponse(t *testing.T, groupID string) {
	url := fmt.Sprintf("%s/%s", rootPath, groupID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

func mockAddRuleResponse(t *testing.T) {
	th.Mux.HandleFunc("/os-security-group-rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "security_group_rule": {
    "from_port": 22,
    "ip_protocol": "TCP",
    "to_port": 22,
    "parent_group_id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
    "cidr": "0.0.0.0/0"
  }
}	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group_rule": {
    "from_port": 22,
    "group": {},
    "ip_protocol": "TCP",
    "to_port": 22,
    "parent_group_id": "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
    "ip_range": {
      "cidr": "0.0.0.0/0"
    },
    "id": "f9a97fcf-3a97-47b0-b76f-919136afb7ed"
  }
}`)
	})
}

func mockDeleteRuleResponse(t *testing.T, ruleID string) {
	url := fmt.Sprintf("/os-security-group-rules/%s", ruleID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

func mockAddServerToGroupResponse(t *testing.T, serverID string) {
	url := fmt.Sprintf("/servers/%s/action", serverID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "addSecurityGroup": {
    "name": "test"
  }
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

func mockRemoveServerFromGroupResponse(t *testing.T, serverID string) {
	url := fmt.Sprintf("/servers/%s/action", serverID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "removeSecurityGroup": {
    "name": "test"
  }
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}
