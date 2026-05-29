package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAvailableDomains(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAvailableDomainsSuccessfully(t, fakeServer)

	count := 0
	err := domains.ListAvailable(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := domains.ExtractDomains(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedAvailableDomainsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListDomains(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListDomainsSuccessfully(t, fakeServer)

	count := 0
	err := domains.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := domains.ExtractDomains(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedDomainsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListDomainsAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListDomainsSuccessfully(t, fakeServer)

	allPages, err := domains.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := domains.ExtractDomains(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedDomainsSlice, actual)
}

func TestGetDomain(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetDomainSuccessfully(t, fakeServer)

	actual, err := domains.Get(context.TODO(), client.ServiceClient(fakeServer), "9fe1d3").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondDomain, *actual)
}

func TestCreateDomain(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateDomainSuccessfully(t, fakeServer)

	createOpts := domains.CreateOpts{
		Name: "domain two",
	}

	actual, err := domains.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondDomain, *actual)
}

func TestDeleteDomain(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteDomainSuccessfully(t, fakeServer)

	res := domains.Delete(context.TODO(), client.ServiceClient(fakeServer), "9fe1d3")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateDomain(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateDomainSuccessfully(t, fakeServer)

	var description = "Staging Domain"
	updateOpts := domains.UpdateOpts{
		Description: &description,
	}

	actual, err := domains.Update(context.TODO(), client.ServiceClient(fakeServer), "9fe1d3", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondDomainUpdated, *actual)
}
