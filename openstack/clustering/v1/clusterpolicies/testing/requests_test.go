package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/clusterpolicies"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListClusterPolicies(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "cluster_policies": [
        {
					"cluster_id":   "7d85f602-a948-4a30-afd4-e84f47471c15",
	        "cluster_name": "cluster4",
	        "enabled":      true,
	        "id":           "06be3a1f-b238-4a96-a737-ceec5714087e",
	        "policy_id":    "714fe676-a08f-4196-b7af-61d52eeded15",
	        "policy_name":  "dp01",
	        "policy_type":  "senlin.policy.deletion-1.0"
				}
    ]
		}`)
	})

	count := 0
	clusterpolicies.ListDetail(fake.ServiceClient(), "1", clusterpolicies.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clusterpolicies.ExtractClusterPolicies(page)
		if err != nil {
			t.Errorf("Failed to extract cluster policies: %v", err)
			return false, err
		}

		expected := []clusterpolicies.ClusterPolicy{
			{
				ClusterUUID: "7d85f602-a948-4a30-afd4-e84f47471c15",
				ClusterName: "cluster4",
				Enabled:     true,
				ID:          "06be3a1f-b238-4a96-a737-ceec5714087e",
				PolicyID:    "714fe676-a08f-4196-b7af-61d52eeded15",
				PolicyName:  "dp01",
				PolicyType:  "senlin.policy.deletion-1.0",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoClusterPolicies(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	clusterpolicies.ListDetail(fake.ServiceClient(), "", clusterpolicies.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := clusterpolicies.ExtractClusterPolicies(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
