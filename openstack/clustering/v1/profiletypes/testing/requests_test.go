package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/profiletypes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListProfileTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "profile_types": [
        {
					"name": "container.dockerinc.docker-1.0"
				}
    ]
		}`)
	})

	count := 0

	profiletypes.ListDetail(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := profiletypes.ExtractProfileTypes(page)
		if err != nil {
			t.Errorf("Failed to extract profile types: %v", err)
			return false, err
		}

		expected := []profiletypes.ProfileType{
			{
				Name: "container.dockerinc.docker-1.0",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoProfileTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	profiletypes.ListDetail(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := profiletypes.ExtractProfileTypes(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
