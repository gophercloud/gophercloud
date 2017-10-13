package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListDomains(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListDomainsSuccessfully(t)

	count := 0
	err := domains.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
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
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListDomainsSuccessfully(t)

	allPages, err := domains.List(client.ServiceClient(), nil).AllPages()
	th.AssertNoErr(t, err)
	actual, err := domains.ExtractDomains(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedDomainsSlice, actual)
}

func TestGetDomain(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetDomainSuccessfully(t)

	actual, err := domains.Get(client.ServiceClient(), "9fe1d3").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondDomain, *actual)
}
