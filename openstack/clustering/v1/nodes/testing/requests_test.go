package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestUpdateNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": "2016-05-13T07:02:20Z",
				"data": {},
				"dependents": {},
				"domain": null,
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": "2016-05-13T08:02:04Z",
				"metadata": {"foo":"bar"},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": "2016-05-13T09:02:04Z",
				"user": "ab79b9647d074e46ac223a8fa297b846"				}
			}
		}`)
	})

	createdAt, _ := time.Parse(time.RFC3339, "2016-05-13T07:02:20Z")
	initAt, _ := time.Parse(time.RFC3339, "2016-05-13T08:02:04Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2016-05-13T09:02:04Z")
	expected := nodes.Node{
		ClusterID:    "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		CreatedAt:    createdAt,
		Data:         nodes.DataType{},
		Dependents:   map[string]interface{}{},
		Domain:       "",
		ID:           "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
		Index:        2,
		InitAt:       initAt,
		Metadata:     map[string]interface{}{"foo": "bar"},
		Name:         "node-e395be1e-002",
		PhysicalID:   "66a81d68-bf48-4af5-897b-a3bfef7279a8",
		ProfileID:    "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		ProfileName:  "pcirros",
		ProjectID:    "eee0b7c083e84501bdd50fb269d2a10e",
		Role:         "",
		Status:       "ACTIVE",
		StatusReason: "Creation succeeded",
		UpdatedAt:    updatedAt,
		User:         "ab79b9647d074e46ac223a8fa297b846",
	}

	nodeOpts := nodes.UpdateOpts{
		Name: "node-e395be1e-002",
	}
	actual, err := nodes.Update(fake.ServiceClient(), "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", nodeOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, *actual)
}

func TestUpdateNodeEmptyTime(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": null,
				"data": {},
				"dependents": {},
				"domain": null,
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": null,
				"metadata": {"foo":"bar"},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": null,
				"user": "ab79b9647d074e46ac223a8fa297b846"				}
			}
		}`)
	})

	expected := nodes.Node{
		ClusterID:    "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		CreatedAt:    time.Time{},
		Data:         nodes.DataType{},
		Dependents:   map[string]interface{}{},
		Domain:       "",
		ID:           "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
		Index:        2,
		InitAt:       time.Time{},
		Metadata:     map[string]interface{}{"foo": "bar"},
		Name:         "node-e395be1e-002",
		PhysicalID:   "66a81d68-bf48-4af5-897b-a3bfef7279a8",
		ProfileID:    "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		ProfileName:  "pcirros",
		ProjectID:    "eee0b7c083e84501bdd50fb269d2a10e",
		Role:         "",
		Status:       "ACTIVE",
		StatusReason: "Creation succeeded",
		UpdatedAt:    time.Time{},
		User:         "ab79b9647d074e46ac223a8fa297b846",
	}

	nodeOpts := nodes.UpdateOpts{
		Name: "node-e395be1e-002",
	}
	actual, err := nodes.Update(fake.ServiceClient(), "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", nodeOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, *actual)
}

func TestUpdateNodeInvalidTimeFloat(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": 123456789.0,
				"data": {},
				"dependents": {},
				"domain": null,
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": 123456789.0,
				"metadata": {"foo":"bar"},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": 123456789.0,
				"user": "ab79b9647d074e46ac223a8fa297b846"				}
			}
		}`)
	})
	nodeOpts := nodes.UpdateOpts{
		Name: "node-e395be1e-002",
	}
	_, err := nodes.Update(fake.ServiceClient(), "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", nodeOpts).Extract()
	th.AssertEquals(t, false, err == nil)
}

func TestUpdateNodeInvalidTimeString(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": "invalid",
				"data": {},
				"dependents": {},
				"domain": null,
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": "invalid",
				"metadata": {"foo":"bar"},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": "invalid",
				"user": "ab79b9647d074e46ac223a8fa297b846"				}
			}
		}`)
	})
	nodeOpts := nodes.UpdateOpts{
		Name: "node-e395be1e-002",
	}
	_, err := nodes.Update(fake.ServiceClient(), "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1", nodeOpts).Extract()
	th.AssertEquals(t, false, err == nil)
}
