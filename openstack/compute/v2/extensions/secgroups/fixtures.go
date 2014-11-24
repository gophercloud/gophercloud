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
      "id": 1,
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
    "id": 1,
    "name": "test",
    "rules": [],
    "tenant_id": "openstack"
  }
}
`)
	})
}

func mockUpdateGroupResponse(t *testing.T, groupID int) {
	url := fmt.Sprintf("%s/%d", rootPath, groupID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "security_group": {
    "name": "new_name",
		"description": "new_desc"
  }
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group": {
    "description": "something",
    "id": 1,
    "name": "new_name",
    "rules": [],
    "tenant_id": "openstack"
  }
}
`)
	})
}

func mockGetGroupsResponse(t *testing.T, groupID int) {
	url := fmt.Sprintf("%s/%d", rootPath, groupID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "security_group": {
    "description": "default",
    "id": 1,
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
        "parent_group_id": 1,
        "ip_range": {
						"cidr": "0.0.0.0"
				},
        "id": 2
      }
    ],
    "tenant_id": "openstack"
  }
}
			`)
	})
}

func mockDeleteGroupResponse(t *testing.T, groupID int) {
	url := fmt.Sprintf("%s/%d", rootPath, groupID)
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
    "parent_group_id": 1,
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
    "parent_group_id": 1,
    "ip_range": {
      "cidr": "0.0.0.0/0"
    },
    "id": 2
  }
}`)
	})
}

func mockDeleteRuleResponse(t *testing.T, ruleID int) {
	url := fmt.Sprintf("/os-security-group-rules/%d", ruleID)
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
