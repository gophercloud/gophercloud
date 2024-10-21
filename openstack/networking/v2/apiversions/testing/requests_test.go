package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/apiversions"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListVersions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "versions": [
        {
            "status": "CURRENT",
            "id": "v2.0",
            "links": [
                {
                    "href": "http://23.253.228.211:9696/v2.0",
                    "rel": "self"
                }
            ]
        }
    ]
}`)
	})

	count := 0

	err := apiversions.ListVersions(fake.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := apiversions.ExtractAPIVersions(page)
		if err != nil {
			t.Errorf("Failed to extract API versions: %v", err)
			return false, err
		}

		expected := []apiversions.APIVersion{
			{
				Status: "CURRENT",
				ID:     "v2.0",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoAPIVersions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := apiversions.ListVersions(fake.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		if _, err := apiversions.ExtractAPIVersions(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
	th.AssertErr(t, err)
}

func TestAPIInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "resources": [
        {
            "links": [
                {
                    "href": "http://23.253.228.211:9696/v2.0/subnets",
                    "rel": "self"
                }
            ],
            "name": "subnet",
            "collection": "subnets"
        },
        {
            "links": [
                {
                    "href": "http://23.253.228.211:9696/v2.0/networks",
                    "rel": "self"
                }
            ],
            "name": "network",
            "collection": "networks"
        },
        {
            "links": [
                {
                    "href": "http://23.253.228.211:9696/v2.0/ports",
                    "rel": "self"
                }
            ],
            "name": "port",
            "collection": "ports"
        }
    ]
}
			`)
	})

	count := 0

	err := apiversions.ListVersionResources(fake.ServiceClient(), "v2.0").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := apiversions.ExtractVersionResources(page)
		if err != nil {
			t.Errorf("Failed to extract version resources: %v", err)
			return false, err
		}

		expected := []apiversions.APIVersionResource{
			{
				Name:       "subnet",
				Collection: "subnets",
			},
			{
				Name:       "network",
				Collection: "networks",
			},
			{
				Name:       "port",
				Collection: "ports",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoAPIVersionResources(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := apiversions.ListVersionResources(fake.ServiceClient(), "v2.0").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		if _, err := apiversions.ExtractVersionResources(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
	th.AssertErr(t, err)
}
