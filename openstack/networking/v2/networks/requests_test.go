package networks

import (
	"fmt"
	"net/http"
	"reflect"
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

func Equals(t *testing.T, actual interface{}, expected interface{}) {
	if expected != actual {
		t.Fatalf("Expected %#v but got %#v", expected, actual)
	}
}

func DeepEquals(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %#v but got %#v", expected, actual)
	}
}

func CheckErr(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("An error occurred: %#v", e)
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
	if err != nil {
		t.Fatalf("Error listing API versions: %v", err)
	}

	coll, err := gophercloud.AllPages(res)

	actual := ToAPIVersions(coll)

	expected := []APIVersion{
		APIVersion{
			Status: "CURRENT",
			ID:     "v2.0",
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
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

	c := ServiceClient()

	res, err := APIInfo(c, "v2.0")
	if err != nil {
		t.Fatalf("Error getting API info: %v", err)
	}

	coll, err := gophercloud.AllPages(res)

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

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
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

		c := ServiceClient()

		ext, err := GetExtension(c, "agent")
		CheckErr(t, err)

		Equals(t, ext.Updated, "2013-02-03T10:00:00-00:00")
		Equals(t, ext.Name, "agent")
		Equals(t, ext.Namespace, "http://docs.openstack.org/ext/agent/api/v2.0")
		Equals(t, ext.Alias, "agent")
		Equals(t, ext.Description, "The agent management extension.")
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

	c := ServiceClient()

	n, err := Get(c, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	Equals(t, n.Status, "ACTIVE")
	DeepEquals(t, n.Subnets, []string{"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"})
	Equals(t, n.Name, "private-network")
	Equals(t, n.ProviderPhysicalNetwork, "")
	Equals(t, n.ProviderNetworkType, "local")
	Equals(t, n.ProviderSegmentationID, "")
	Equals(t, n.AdminStateUp, true)
	Equals(t, n.TenantID, "4fd44f30292945e481c7b8a0c8908869")
	Equals(t, n.RouterExternal, true)
	Equals(t, n.Shared, true)
	Equals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
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
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	Equals(t, n.Status, "ACTIVE")
	DeepEquals(t, n.Subnets, []interface{}{})
	Equals(t, n.Name, "net1")
	Equals(t, n.AdminStateUp, true)
	Equals(t, n.TenantID, "9bacb3c5d39d41a79512987f338cf177")
	DeepEquals(t, n.Segments, []NetworkProvider{
		{ProviderSegmentationID: 2, ProviderPhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453", ProviderNetworkType: "vlan"},
		{ProviderSegmentationID: 0, ProviderPhysicalNetwork: "8bab8453-1bc9-45af-8c70-f83aa9b50453", ProviderNetworkType: "stt"},
	})
	Equals(t, n.Shared, false)
	Equals(t, n.PortSecurityEnabled, true)
	Equals(t, n.ID, "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
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
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
}
