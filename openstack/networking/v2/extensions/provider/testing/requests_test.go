package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/provider"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "networks": [
        {
            "status": "ACTIVE",
            "subnets": [
                "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
            ],
            "name": "private-network",
            "admin_state_up": true,
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "shared": true,
            "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
            "provider:segmentation_id": "2",
            "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
            "provider:network_type": "vlan"
        },
        {
            "status": "ACTIVE",
            "subnets": [
                "08eae331-0402-425a-923c-34f7cfe39c1b"
            ],
            "name": "private",
            "admin_state_up": true,
            "tenant_id": "26a7980765d0414dbc1fc1f88cdb7e6e",
            "shared": true,
            "id": "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
            "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
            "provider:network_type": "stt"
        }
    ]
}
			`)
	})

	client := fake.ServiceClient()
	count := 0

	provider.List(client, provider.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := provider.ExtractNetworkExtAttrs(page)
		if err != nil {
			t.Errorf("Failed to extract networks: %v", err)
			return false, err
		}

		expected := []provider.NetworkExtAttrs{
			{
				Status:          "ACTIVE",
				Subnets:         []string{"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"},
				Name:            "private-network",
				AdminStateUp:    true,
				TenantID:        "4fd44f30292945e481c7b8a0c8908869",
				Shared:          true,
				ID:              "d32019d3-bc6e-4319-9c1d-6722fc136a22",
				SegmentationID:  "2",
				PhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453",
				NetworkType:     "vlan",
			},
			{
				Status:          "ACTIVE",
				Subnets:         []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
				Name:            "private",
				AdminStateUp:    true,
				TenantID:        "26a7980765d0414dbc1fc1f88cdb7e6e",
				Shared:          true,
				ID:              "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
				PhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453",
				NetworkType:     "stt",
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

	th.Mux.HandleFunc("/v2.0/networks/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "network": {
        "status": "ACTIVE",
        "subnets": [
            "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
        ],
        "name": "private-network",
        "admin_state_up": true,
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "shared": true,
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
        "provider:segmentation_id": 2,
        "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
        "provider:network_type": "vlan"
    }
}
			`)
	})

	n, err := provider.Get(fake.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertDeepEquals(t, n.Subnets, []string{"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"})
	th.AssertEquals(t, n.Name, "private-network")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.TenantID, "4fd44f30292945e481c7b8a0c8908869")
	th.AssertEquals(t, n.Shared, true)
	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertEquals(t, n.PhysicalNetwork, "8bab8453-1bc9-45af-8c70-f83aa9b50453")
	th.AssertEquals(t, n.NetworkType, "vlan")
	th.AssertEquals(t, n.SegmentationID, "2")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "network": {
        "name": "sample_network",
        "admin_state_up": true
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "network": {
        "status": "ACTIVE",
        "subnets": [],
        "name": "net1",
        "admin_state_up": true,
        "tenant_id": "9bacb3c5d39d41a79512987f338cf177",
        "shared": false,
        "id": "4e8e5957-649f-477b-9e5b-f1f75b21c03c"
    }
}
		`)
	})

	iTrue := true
	options := provider.CreateOpts{Name: "sample_network", AdminStateUp: &iTrue}
	n, err := provider.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertDeepEquals(t, n.Subnets, []string{})
	th.AssertEquals(t, n.Name, "net1")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.TenantID, "9bacb3c5d39d41a79512987f338cf177")
	th.AssertEquals(t, n.Shared, false)
	th.AssertEquals(t, n.ID, "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
}

func TestCreateWithOptionalFields(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"network": {
			"name": "sample_network",
			"admin_state_up": true,
			"shared": true,
			"tenant_id": "12345",
                        "provider:segmentation_id": "2",
                        "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
                        "provider:network_type": "vlan"
	}
}
		`)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{}`)
	})

	iTrue := true
	options := provider.CreateOpts{Name: "sample_network", AdminStateUp: &iTrue, Shared: &iTrue, TenantID: "12345", SegmentationID: "2", PhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453", NetworkType: "vlan"}
	_, err := provider.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
		"network": {
				"name": "new_network_name",
				"admin_state_up": false,
				"shared": true,
                                "provider:segmentation_id": "2",
                                "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
                                "provider:network_type": "vlan"
		}
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "network": {
        "status": "ACTIVE",
        "subnets": [],
        "name": "new_network_name",
        "admin_state_up": false,
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "shared": true,
        "id": "4e8e5957-649f-477b-9e5b-f1f75b21c03c",
        "provider:segmentation_id": 2,
        "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
        "provider:network_type": "vlan"
    }
}
		`)
	})

	iTrue, iFalse := true, false
	options := provider.UpdateOpts{
		Name:            "new_network_name",
		AdminStateUp:    &iFalse,
		Shared:          &iTrue,
		SegmentationID:  "2",
		PhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453",
		NetworkType:     "vlan",
	}
	n, err := provider.Update(fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Name, "new_network_name")
	th.AssertEquals(t, n.AdminStateUp, false)
	th.AssertEquals(t, n.Shared, true)
	th.AssertEquals(t, n.ID, "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
	th.AssertEquals(t, n.TenantID, "4fd44f30292945e481c7b8a0c8908869")
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := networks.Delete(fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
	th.AssertNoErr(t, res.Err)
}

func TestIDFromName(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "networks": [
        {
            "status": "ACTIVE",
            "subnets": [
                "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
            ],
            "name": "private-network",
            "admin_state_up": true,
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "shared": true,
            "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
            "provider:segmentation_id": "2",
            "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
            "provider:network_type": "vlan"
        },
        {
            "status": "ACTIVE",
            "subnets": [
                "08eae331-0402-425a-923c-34f7cfe39c1b"
            ],
            "name": "private",
            "admin_state_up": true,
            "tenant_id": "26a7980765d0414dbc1fc1f88cdb7e6e",
            "shared": true,
            "id": "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
            "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
            "provider:network_type": "stt"
        }
    ]
}
			`)
	})

	client := fake.ServiceClient()
	expectID := "d32019d3-bc6e-4319-9c1d-6722fc136a22"

	id, err := provider.IDFromName(client, "provider-network")
	if err != nil {
		t.Errorf("Network %s not found.", id)
		return "", err
	}

	th.AssertEquals(t, id, expectID)
}