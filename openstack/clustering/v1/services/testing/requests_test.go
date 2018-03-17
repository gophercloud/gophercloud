package testing

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/services"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"net/http"
	"testing"
)

func TestListReceivers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "services": [
        {
					"binary": "senlin-engine",
					"disabled_reason": null,
					"host": "host1",
					"id": "f93f83f6-762b-41b6-b757-80507834d394",
					"state": "up",
					"status": "enabled",
					"topic": "senlin-engine",
					"updated_at": "2017-04-24T07:43:12"
				}
    ]
		}`)
	})

	count := 0

	services.ListDetail(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := services.ExtractServices(page)
		if err != nil {
			t.Errorf("Failed to extract services: %v", err)
			return false, err
		}

		expected := []services.Service{
			{
				Binary:        "senlin-engine",
				DisableReason: "",
				Host:          "host1",
				ID:            "f93f83f6-762b-41b6-b757-80507834d394",
				State:         "up",
				Status:        "enabled",
				Topic:         "senlin-engine",
				UpdatedAt:     "2017-04-24T07:43:12",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoServices(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	services.ListDetail(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := services.ExtractServices(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
