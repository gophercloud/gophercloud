// +build acceptance networking

package v2

import (
	"os"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/utils"
)

var Client *gophercloud.ServiceClient

func NewClient() (*gophercloud.ServiceClient, error) {
	opts, err := utils.AuthOptions()
	if err != nil {
		return nil, err
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	return openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

func Setup(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	Client = client
}

func Teardown() {
	Client = nil
}

func Equals(t *testing.T, actual interface{}, expected interface{}) {
	if expected != actual {
		t.Fatalf("Expected %#v but got %#v", expected, actual)
	}
}

func TestListAPIVersions(t *testing.T) {
	Setup(t)
	defer Teardown()

	res, err := networks.APIVersions(Client)
	if err != nil {
		t.Fatalf("Failed to list API versions: %v", err)
	}

	err = gophercloud.EachPage(res, func(page gophercloud.Collection) bool {
		t.Logf("--- Page ---")
		for _, v := range networks.ToAPIVersions(page) {
			t.Logf("API version: ID [%s] Status [%s]", v.ID, v.Status)
		}
		return true
	})
	if err != nil {
		t.Fatalf("Unexpected error while iterating API versions: %v", err)
	}
}

func TestGetApiInfo(t *testing.T) {
	Setup(t)
	defer Teardown()

	res, err := networks.APIInfo(Client, "v2.0")
	if err != nil {
		t.Fatalf("Failed to list API info for v2: %v", err)
	}

	err = gophercloud.EachPage(res, func(page gophercloud.Collection) bool {
		t.Logf("--- Page ---")
		for _, r := range networks.ToAPIResource(page) {
			t.Logf("API resource: Name [%s] Collection [%s]", r.Name, r.Collection)
		}
		return true
	})
	if err != nil {
		t.Fatalf("Unexpected error while iteratoring API resources: %v", err)
	}
}

func TestListExts(t *testing.T) {
	//networks.Extensions()
}

func TestGetExt(t *testing.T) {
	Setup(t)
	defer Teardown()

	ext, err := networks.GetExtension(Client, "service-type")
	if err != nil {
		t.Fatalf("Unexpected error when getting extension: %#v", err)
	}

	Equals(t, ext.Updated, "2013-01-20T00:00:00-00:00")
	Equals(t, ext.Name, "Neutron Service Type Management")
	Equals(t, ext.Namespace, "http://docs.openstack.org/ext/neutron/service-type/api/v1.0")
	Equals(t, ext.Alias, "service-type")
	Equals(t, ext.Description, "API for retrieving service providers for Neutron advanced services")
}

func TestListNetworks(t *testing.T) {
	//networks.List()
}

func TestGetNetwork(t *testing.T) {
	//networks.Get()
}

func TestCreateNetwork(t *testing.T) {
	//networks.Create()
}

func TestCreateMultipleNetworks(t *testing.T) {
	//networks.CreateMany()
}

func TestUpdateNetwork(t *testing.T) {
	//networks.Update()
}

func TestDeleteNetwork(t *testing.T) {
	//networks.Delete()
}
