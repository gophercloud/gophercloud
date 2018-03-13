package testing

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/clusters"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"net/http"
	"testing"
	"time"
)

func TestListClusters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "clusters": [
        {
					"config": {},
					"created_at": "2015-02-10T14:26:14Z",
					"data": {},
					"dependents": {},
					"desired_capacity": 4,
					"domain": null,
					"id": "7d85f602-a948-4a30-afd4-e84f47471c15",
					"init_at": "2015-02-10T14:26:14Z",
					"max_size": -1,
					"metadata": {},
					"min_size": 0,
					"name": "cluster1",
					"nodes": [
							"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
							"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
							"da1e9c87-e584-4626-a120-022da5062dac"
					],
					"policies": [],
					"profile_id": "edc63d0a-2ca4-48fa-9854-27926da76a4a",
					"profile_name": "mystack",
					"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
					"status": "ACTIVE",
					"status_reason": "Cluster scale-in succeeded",
					"timeout": 3600,
					"updated_at": null,
					"user": "5e5bf8027826429c96af157f68dc9072"
				}
    ]
		}`)
	})

	count := 0

	clusters.ListDetail(fake.ServiceClient(), clusters.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractClusters(page)
		if err != nil {
			t.Errorf("Failed to extract clusters: %v", err)
			return false, err
		}

		createdAt, _ := time.Parse(clusters.RFC3339Milli, "2015-02-10T14:26:14Z")
		initAt, _ := time.Parse(clusters.RFC3339Milli, "2015-02-10T14:26:14Z")
		expected := []clusters.Cluster{
			{
				Config:          map[string]interface{}{},
				CreatedAt:       createdAt,
				Data:            map[string]interface{}{},
				Dependents:      map[string]interface{}{},
				DesiredCapacity: 4,
				DomainUUID:      "",
				ID:              "7d85f602-a948-4a30-afd4-e84f47471c15",
				InitAt:          initAt,
				MaxSize:         -1,
				Metadata:        map[string]interface{}{},
				MinSize:         0,
				Name:            "cluster1",
				Nodes: []string{
					"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
					"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
					"da1e9c87-e584-4626-a120-022da5062dac",
				},
				Policies:     []string{},
				ProfileID:    "edc63d0a-2ca4-48fa-9854-27926da76a4a",
				ProfileName:  "mystack",
				ProjectUUID:  "6e18cc2bdbeb48a5b3cad2dc499f6804",
				Status:       "ACTIVE",
				StatusReason: "Cluster scale-in succeeded",
				Timeout:      3600,
				UpdatedAt:    time.Time{},
				UserUUID:     "5e5bf8027826429c96af157f68dc9072"},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoClusters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	clusters.ListDetail(fake.ServiceClient(), clusters.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := clusters.ExtractClusters(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
