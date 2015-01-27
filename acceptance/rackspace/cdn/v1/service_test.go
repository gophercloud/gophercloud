// +build acceptance

package v1

import (
	"testing"

	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/cdn/v1/services"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/cdn/v1/services"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestService(t *testing.T) {
	client := newClient(t)

	t.Log("Creating Service")
	loc := testServiceCreate(t, client)
	t.Logf("Created service at location: %s", loc)

	defer testServiceDelete(t, client, loc)

	t.Log("Updating Service")
	testServiceUpdate(t, client, loc)

	t.Log("Retrieving Service")
	testServiceGet(t, client, loc)

	t.Log("Listing Services")
	testServiceList(t, client)
}

func testServiceCreate(t *testing.T, client *gophercloud.ServiceClient) string {
	createOpts := os.CreateOpts{
		Name: "gophercloud-test-service",
		Domains: []os.Domain{
			os.Domain{
				Domain: "www.gophercloud-test-service.com",
			},
		},
		Origins: []os.Origin{
			os.Origin{
				Origin: "gophercloud-test-service.com",
				Port:   80,
				SSL:    false,
			},
		},
		FlavorID: "cdn",
	}
	l, err := services.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)
	return l
}

func testServiceGet(t *testing.T, client *gophercloud.ServiceClient, id string) {
	s, err := services.Get(client, id).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Retrieved service: %+v", *s)
}

func testServiceUpdate(t *testing.T, client *gophercloud.ServiceClient, id string) {
	updateOpts := os.UpdateOpts{
		os.UpdateOpt{
			Op:   os.Add,
			Path: "/domains/-",
			Value: map[string]interface{}{
				"domain":   "newDomain.com",
				"protocol": "http",
			},
		},
	}
	loc, err := services.Update(client, id, updateOpts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Successfully updated service at location: %s", loc)
}

func testServiceList(t *testing.T, client *gophercloud.ServiceClient) {
	err := services.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		serviceList, err := os.ExtractServices(page)
		th.AssertNoErr(t, err)

		for _, service := range serviceList {
			t.Logf("Listing service: %+v", service)
		}

		return true, nil
	})

	th.AssertNoErr(t, err)
}

func testServiceDelete(t *testing.T, client *gophercloud.ServiceClient, id string) {
	err := services.Delete(client, id).ExtractErr()
	th.AssertNoErr(t, err)
	t.Logf("Successfully deleted service (%s)", id)
}
