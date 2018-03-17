package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/pagination"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/profiles"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListBuildProfiles(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "profiles": [
        {
					"created_at": "2016-01-03T16:22:23Z",
					"domain": null,
					"id": "9e1c6f42-acf5-4688-be2c-8ce954ef0f23",
					"metadata": {},
					"name": "pserver",
					"project": "42d9e9663331431f97b75e25136307ff",
					"spec": {
							"properties": {
									"flavor": 1,
									"image": "cirros-0.3.4-x86_64-uec",
									"key_name": "oskey",
									"name": "cirros_server",
									"networks": [
											{
													"network": "private"
											}
									]
							},
							"type": "os.nova.server",
							"version": 1.0
					},
					"type": "os.nova.server-1.0",
					"updated_at": null,
					"user": "5e5bf8027826429c96af157f68dc9072"
				}
    ]
		}`)
	})

	count := 0

	profiles.ListDetail(fake.ServiceClient(), profiles.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := profiles.ExtractProfiles(page)
		if err != nil {
			t.Errorf("Failed to extract profiles: %v", err)
			return false, err
		}

		createdAt, _ := time.Parse(time.RFC3339, "2016-01-03T16:22:23Z")
		updatedAt := time.Time{}

		expected := []profiles.Profile{
			{
				CreatedAt:   createdAt,
				DomainUUID:  "",
				ID:          "9e1c6f42-acf5-4688-be2c-8ce954ef0f23",
				Metadata:    map[string]interface{}{},
				Name:        "pserver",
				ProjectUUID: "42d9e9663331431f97b75e25136307ff",
				Spec: map[string]interface{}{
					"properties": map[string]interface{}{
						"flavor":   float64(1),
						"image":    "cirros-0.3.4-x86_64-uec",
						"key_name": "oskey",
						"name":     "cirros_server",
						"networks": []interface{}{
							map[string]interface{}{"network": "private"},
						},
					},
					"type":    "os.nova.server",
					"version": 1.0,
				},
				Type:      "os.nova.server-1.0",
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

func TestNonJSONCannotBeExtractedIntoProfiles(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	profiles.ListDetail(fake.ServiceClient(), profiles.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := profiles.ExtractProfiles(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
