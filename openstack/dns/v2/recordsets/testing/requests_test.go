package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListByZone(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListByZoneSuccessfully(t)

	count := 0
	err := recordsets.ListByZone(client.ServiceClient(), "2150b1bf-dee2-4221-9d85-11f7886fb15f", nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := recordsets.ExtractRecordSets(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedRecordSetSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListByZoneAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListByZoneSuccessfully(t)

	allPages, err := recordsets.ListByZone(client.ServiceClient(), "2150b1bf-dee2-4221-9d85-11f7886fb15f", nil).AllPages()
	th.AssertNoErr(t, err)
	allRecordSets, err := recordsets.ExtractRecordSets(allPages)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 2, len(allRecordSets))
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := recordsets.Get(client.ServiceClient(), "2150b1bf-dee2-4221-9d85-11f7886fb15f", "f7b10e9b-0cae-4a91-b162-562bc6096648").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstRecordSet, actual)
}
