package snapshots

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

const TokenID = "123"

func ServiceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{
			TokenID: TokenID,
		},
		Endpoint: th.Endpoint(),
	}
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/snapshots", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"snapshots": [
				{
					"id": "289da7f8-6440-407c-9fb4-7db01ec49164",
					"display_name": "snapshot-001"
				},
				{
					"id": "96c3bda7-c82a-4f50-be73-ca7621794835",
					"display_name": "snapshot-002"
				}
			]
		}
		`)
	})

	client := ServiceClient()
	count := 0

	List(client, &ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractSnapshots(page)
		if err != nil {
			t.Errorf("Failed to extract snapshots: %v", err)
			return false, err
		}

		expected := []Snapshot{
			Snapshot{
				ID:   "289da7f8-6440-407c-9fb4-7db01ec49164",
				Name: "snapshot-001",
			},
			Snapshot{
				ID:   "96c3bda7-c82a-4f50-be73-ca7621794835",
				Name: "snapshot-002",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/snapshots/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
    "snapshot": {
        "display_name": "snapshot-001",
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22"
    }
}
			`)
	})

	v, err := Get(ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, v.Name, "snapshot-001")
	th.AssertEquals(t, v.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/snapshots", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "snapshot": {
        "display_name": "snapshot-001"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "snapshot": {
        "display_name": "snapshot-001",
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22"
    }
}
		`)
	})

	options := &CreateOpts{Name: "snapshot-001"}
	n, err := Create(ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Name, "snapshot-001")
	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestUpdateMetadata(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/snapshots/123/metadata", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `
		{
			"metadata": {
				"key": "v1"
			}
		}
		`)

		fmt.Fprintf(w, `
			{
				"metadata": {
					"key": "v1"
				}
			}
		`)
	})

	expected := map[string]interface{}{"key": "v1"}

	options := &UpdateMetadataOpts{
		Metadata: map[string]interface{}{
			"key": "v1",
		},
	}
	actual, err := UpdateMetadata(ServiceClient(), "123", options).ExtractMetadata()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/snapshots/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	err := Delete(ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, err)
}
