package testing

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/policies"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"net/http"
	"testing"
	"time"
)

func TestListPolicies(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "policies": [
        {
					"created_at": "2015-02-15T08:33:13.000000Z",
					"data": {},
					"domain": null,
					"id": "7192d8df-73be-4e98-ab99-1cf6d5066729",
					"name": "test_policy_1",
					"project": "42d9e9663331431f97b75e25136307ff",
					"spec": {
							"description": "A test policy",
							"properties": {
									"criteria": "OLDEST_FIRST",
									"destroy_after_deletion": true,
									"grace_period": 60,
									"reduce_desired_capacity": false
							},
							"type": "senlin.policy.deletion",
							"version": "1.0"
					},
					"type": "senlin.policy.deletion-1.0",
					"updated_at": null,
					"user": "5e5bf8027826429c96af157f68dc9072"
				}
    ]
		}`)
	})

	count := 0

	policies.ListDetail(fake.ServiceClient(), policies.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := policies.ExtractPolicies(page)
		if err != nil {
			t.Errorf("Failed to extract policies: %v", err)
			return false, err
		}

		createdAt, _ := time.Parse(time.RFC3339, "2015-02-15T08:33:13.000000Z")
		updatedAt := time.Time{}

		expected := []policies.Policy{
			{
				CreatedAt:   createdAt,
				Data:        map[string]interface{}{},
				DomainUUID:  "",
				ID:          "7192d8df-73be-4e98-ab99-1cf6d5066729",
				Name:        "test_policy_1",
				ProjectUUID: "42d9e9663331431f97b75e25136307ff",

				Spec: map[string]interface{}{
					"description": "A test policy",
					"properties": map[string]interface{}{
						"criteria":                "OLDEST_FIRST",
						"destroy_after_deletion":  true,
						"grace_period":            float64(60),
						"reduce_desired_capacity": false,
					},
					"type":    "senlin.policy.deletion",
					"version": "1.0",
				},
				Type:      "senlin.policy.deletion-1.0",
				UpdatedAt: updatedAt,
				UserUUID:  "5e5bf8027826429c96af157f68dc9072",
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

	policies.ListDetail(fake.ServiceClient(), policies.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := policies.ExtractPolicies(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
