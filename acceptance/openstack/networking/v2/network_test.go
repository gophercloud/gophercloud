// +build acceptance networking

package v2

import (
	"os"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/utils"
	th "github.com/rackspace/gophercloud/testhelper"
)

var (
	Client *gophercloud.ServiceClient
)

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
	th.AssertNoErr(t, err)
	Client = client
}

func Teardown() {
	Client = nil
}

func TestListAPIVersions(t *testing.T) {
	Setup(t)
	defer Teardown()

	res, err := networks.APIVersions(Client)
	th.AssertNoErr(t, err)

	err = gophercloud.EachPage(res, func(page gophercloud.Collection) bool {
		t.Logf("--- Page ---")
		for _, v := range networks.ToAPIVersions(page) {
			t.Logf("API version: ID [%s] Status [%s]", v.ID, v.Status)
		}
		return true
	})
	th.AssertNoErr(t, err)
}

func TestGetApiInfo(t *testing.T) {
	Setup(t)
	defer Teardown()

	res, err := networks.APIInfo(Client, "v2.0")
	th.AssertNoErr(t, err)

	err = gophercloud.EachPage(res, func(page gophercloud.Collection) bool {
		t.Logf("--- Page ---")
		for _, r := range networks.ToAPIResource(page) {
			t.Logf("API resource: Name [%s] Collection [%s]", r.Name, r.Collection)
		}
		return true
	})
	th.AssertNoErr(t, err)
}

func TestListExts(t *testing.T) {
	//networks.Extensions()
}

func TestGetExt(t *testing.T) {
	Setup(t)
	defer Teardown()

	ext, err := networks.GetExtension(Client, "service-type")
	th.AssertNoErr(t, err)

	th.AssertEquals(t, ext.Updated, "2013-01-20T00:00:00-00:00")
	th.AssertEquals(t, ext.Name, "Neutron Service Type Management")
	th.AssertEquals(t, ext.Namespace, "http://docs.openstack.org/ext/neutron/service-type/api/v1.0")
	th.AssertEquals(t, ext.Alias, "service-type")
	th.AssertEquals(t, ext.Description, "API for retrieving service providers for Neutron advanced services")
}

func TestListNetworks(t *testing.T) {
	//networks.List()
}

func TestNetworkCRUDOperations(t *testing.T) {
	Setup(t)
	defer Teardown()

	// Create a network
	res, err := networks.Create(Client, networks.NetworkOpts{Name: "sample_network", AdminStateUp: true})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, res.Name, "sample_network")
	th.AssertEquals(t, res.AdminStateUp, true)
	networkID := res.ID

	// Get a network
	if networkID == "" {
		t.Fatalf("In order to retrieve a network, the NetworkID must be set")
	}
	n, err := networks.Get(Client, networkID)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertDeepEquals(t, n.Subnets, []interface{}{})
	th.AssertEquals(t, n.Name, "sample_network")
	th.AssertEquals(t, n.ProviderPhysicalNetwork, "")
	th.AssertEquals(t, n.ProviderNetworkType, "local")
	th.AssertEquals(t, n.ProviderSegmentationID, 0)
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.RouterExternal, false)
	th.AssertEquals(t, n.Shared, false)
	th.AssertEquals(t, n.ID, networkID)

	// Update network
	n, err = networks.Update(Client, networkID, networks.NetworkOpts{Name: "new_network_name"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Name, "new_network_name")

	// Delete network
	err := networks.Delete(Client, networkID)
	th.AssertNoErr(t, err)
}

func TestCreateMultipleNetworks(t *testing.T) {
	//networks.CreateMany()
}
