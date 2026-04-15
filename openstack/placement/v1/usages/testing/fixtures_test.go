package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/usages"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const ProjectID = "42a2b0fa980d4f7f873e8f0d8b4e1b0e"
const UserID = "b29d880b2c114d5d9a55748b26b2e41e"

var GetUsagesBody = `
{
    "usages": {
        "INSTANCE": {
            "VCPU": 2,
            "MEMORY_MB": 2048,
            "consumer_count": 1
        }
    }
}
`

var ExpectedUsages = usages.Usages{
	Usages: map[string]usages.ConsumerTypeUsage{
		"INSTANCE": {
			"VCPU":           2,
			"MEMORY_MB":      2048,
			"consumer_count": 1,
		},
	},
}

var GetUsagesWithUserBody = `
{
    "usages": {
        "INSTANCE": {
            "VCPU": 2,
            "MEMORY_MB": 2048,
            "consumer_count": 1
        }
    }
}
`

var GetEmptyUsagesBody = `
{
    "usages": {}
}
`

var ExpectedEmptyUsages = usages.Usages{
	Usages: map[string]usages.ConsumerTypeUsage{},
}

var GetUsagesPre138Body = `
{
    "usages": {
        "VCPU": 2,
        "MEMORY_MB": 2048
    }
}
`

var ExpectedUsagesPre138 = usages.UsagesPre138{
	Usages: map[string]int{
		"VCPU":      2,
		"MEMORY_MB": 2048,
	},
}

var GetEmptyUsagesPre138Body = `
{
    "usages": {}
}
`

var ExpectedEmptyUsagesPre138 = usages.UsagesPre138{
	Usages: map[string]int{},
}

func HandleGetUsagesSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/usages", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		th.AssertEquals(t, ProjectID, r.URL.Query().Get("project_id"))
		th.AssertEquals(t, "", r.URL.Query().Get("user_id"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetUsagesBody)
	})
}

func HandleGetUsagesWithUserSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/usages", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		th.AssertEquals(t, ProjectID, r.URL.Query().Get("project_id"))
		th.AssertEquals(t, UserID, r.URL.Query().Get("user_id"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetUsagesWithUserBody)
	})
}

func HandleGetEmptyUsagesSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/usages", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		th.AssertEquals(t, ProjectID, r.URL.Query().Get("project_id"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetEmptyUsagesBody)
	})
}

func HandleGetUsagesPre138Success(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/usages", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		th.AssertEquals(t, ProjectID, r.URL.Query().Get("project_id"))
		th.AssertEquals(t, "", r.URL.Query().Get("user_id"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetUsagesPre138Body)
	})
}

func HandleGetUsagesPre138WithUserSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/usages", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		th.AssertEquals(t, ProjectID, r.URL.Query().Get("project_id"))
		th.AssertEquals(t, UserID, r.URL.Query().Get("user_id"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetUsagesPre138Body)
	})
}

func HandleGetEmptyUsagesPre138Success(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/usages", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		th.AssertEquals(t, ProjectID, r.URL.Query().Get("project_id"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetEmptyUsagesPre138Body)
	})
}
