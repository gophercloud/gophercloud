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
