package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListNodes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "nodes": [
        {
					"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
					"created_at": "2016-05-13T07:02:20Z",
					"data": {},
					"dependents": {},
					"domain": null,
					"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
					"index": 2,
					"init_at": "2016-05-13T07:02:04Z",
					"metadata": {},
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
    ]
		}`)
	})

	count := 0

	nodes.ListDetail(fake.ServiceClient(), nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := nodes.ExtractNodes(page)
		if err != nil {
			t.Errorf("Failed to extract nodes: %v", err)
			return false, err
		}

		createdAt, _ := time.Parse(time.RFC3339, "2016-05-13T07:02:20Z")
		initAt, _ := time.Parse(time.RFC3339, "2016-05-13T07:02:04Z")
		updatedAt := time.Time{}

		expected := []nodes.Node{
			{
				ClusterUUID:  "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				CreatedAt:    createdAt,
				Data:         nodes.DataType{},
				Dependents:   map[string]interface{}{},
				DomainUUID:   "",
				ID:           "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				Index:        2,
				InitAt:       initAt,
				Metadata:     map[string]interface{}{},
				Name:         "node-e395be1e-002",
				PhysicalUUID: "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				ProfileUUID:  "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				ProfileName:  "pcirros",
				ProjectUUID:  "eee0b7c083e84501bdd50fb269d2a10e",
				Role:         "",
				Status:       "ACTIVE",
				StatusReason: "Creation succeeded",
				UpdatedAt:    updatedAt,
				UserUUID:     "ab79b9647d074e46ac223a8fa297b846",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoNodes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	nodes.ListDetail(fake.ServiceClient(), nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := nodes.ExtractNodes(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
