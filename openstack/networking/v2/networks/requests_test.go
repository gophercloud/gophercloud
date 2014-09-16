package networks

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
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

func TestListAPIVersions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

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

	c := ServiceClient()

	res, err := APIVersions(c)
	th.AssertNoErr(err)

	coll, err := gophercloud.AllPages(res)
	th.AssertNoErr(err)

	actual := ToAPIVersions(coll)

	expected := []APIVersion{
		APIVersion{
			Status: "CURRENT",
			ID:     "v2.0",
		},
	}
	th.AssertDeepEquals(expected, actual)
}

func TestAPIInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

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

	res, err := APIInfo(ServiceClient(), "v2.0")
	th.AssertNoErr(err)

	coll, err := gophercloud.AllPages(res)
	th.AssertNoErr(err)

	actual := ToAPIResource(coll)
	expected := []APIResource{
		APIResource{
			Name:       "subnet",
			Collection: "subnets",
		},
		APIResource{
			Name:       "network",
			Collection: "networks",
		},
		APIResource{
			Name:       "port",
			Collection: "ports",
		},
	}
	th.AssertDeepEquals(expected, actual)
}

func TestListingExtensions(t *testing.T) {

}

func TestGettingExtension(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/extension/agent", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "extension": {
        "updated": "2013-02-03T10:00:00-00:00",
        "name": "agent",
        "links": [],
        "namespace": "http://docs.openstack.org/ext/agent/api/v2.0",
        "alias": "agent",
        "description": "The agent management extension."
    }
}
		`)

		ext, err := GetExtension(ServiceClient(), "agent")
		th.AssertNoErr(t, err)

		th.AssertEquals(t, ext.Updated, "2013-02-03T10:00:00-00:00")
		th.AssertEquals(t, ext.Name, "agent")
		th.AssertEquals(t, ext.Namespace, "http://docs.openstack.org/ext/agent/api/v2.0")
		th.AssertEquals(t, ext.Alias, "agent")
		th.AssertEquals(t, ext.Description, "The agent management extension.")
	})
}

func TestGettingNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

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
        "provider:physical_network": null,
        "admin_state_up": true,
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "provider:network_type": "local",
        "router:external": true,
        "shared": true,
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
        "provider:segmentation_id": null
    }
}
			`)
	})

	n, err := Get(ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertDeepEquals(t, n.Subnets, []interface{}{"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"})
	th.AssertEquals(t, n.Name, "private-network")
	th.AssertEquals(t, n.ProviderPhysicalNetwork, "")
	th.AssertEquals(t, n.ProviderNetworkType, "local")
	th.AssertEquals(t, n.ProviderSegmentationID, 0)
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.TenantID, "4fd44f30292945e481c7b8a0c8908869")
	th.AssertEquals(t, n.RouterExternal, true)
	th.AssertEquals(t, n.Shared, true)
	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestCreateNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)
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
        "segments": [
            {
                "provider:segmentation_id": 2,
                "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
                "provider:network_type": "vlan"
            },
            {
                "provider:segmentation_id": null,
                "provider:physical_network": "8bab8453-1bc9-45af-8c70-f83aa9b50453",
                "provider:network_type": "stt"
            }
        ],
        "shared": false,
        "port_security_enabled": true,
        "id": "4e8e5957-649f-477b-9e5b-f1f75b21c03c"
    }
}
		`)
	})

	options := NetworkOpts{Name: "sample_network", AdminStateUp: true}
	n, err := Create(ServiceClient(), options)
	th.AssertNoErr(err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertDeepEquals(t, n.Subnets, []interface{}{})
	th.AssertEquals(t, n.Name, "net1")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.TenantID, "9bacb3c5d39d41a79512987f338cf177")
	th.AssertDeepEquals(t, n.Segments, []NetworkProvider{
		{ProviderSegmentationID: 2, ProviderPhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453", ProviderNetworkType: "vlan"},
		{ProviderSegmentationID: 0, ProviderPhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453", ProviderNetworkType: "stt"},
	})
	th.AssertEquals(t, n.Shared, false)
	th.AssertEquals(t, n.PortSecurityEnabled, true)
	th.AssertEquals(t, n.ID, "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
}

func TestCreateNetworkWithOptionalFields(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"network": {
			"name": "sample_network",
			"admin_state_up": true,
			"shared": true,
			"tenant_id": "12345"
	}
}
		`)

		w.WriteHeader(http.StatusCreated)
	})

	shared := true
	options := NetworkOpts{Name: "sample_network", AdminStateUp: true, Shared: &shared, TenantID: "12345"}
	_, err := Create(ServiceClient(), options)
	th.AssertNoErr(err)
}

func TestUpdateNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
		"network": {
				"name": "new_network_name",
				"admin_state_up": false,
				"shared": true
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
        "provider:physical_network": null,
        "admin_state_up": false,
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "provider:network_type": "local",
        "router:external": false,
        "shared": true,
        "id": "4e8e5957-649f-477b-9e5b-f1f75b21c03c",
        "provider:segmentation_id": null
    }
}
		`)
	})

	shared := true
	options := NetworkOpts{Name: "new_network_name", AdminStateUp: false, Shared: &shared}
	n, err := Update(ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options)
	th.AssertNoErr(err)

	th.AssertEquals(t, n.Name, "new_network_name")
	th.AssertEquals(t, n.AdminStateUp, false)
	th.AssertEquals(t, n.Shared, true)
	th.AssertEquals(t, n.ID, "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
}

func TestDeleteNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	err := Delete(ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
	th.AssertNoErr(err)
}
