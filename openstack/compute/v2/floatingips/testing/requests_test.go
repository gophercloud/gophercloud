package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/floatingips"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	count := 0
	err := floatingips.List(client.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := floatingips.ExtractFloatingIPs(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedFloatingIPsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	actual, err := floatingips.Create(context.TODO(), client.ServiceClient(), floatingips.CreateOpts{
		Pool: "nova",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedFloatingIP, actual)
}

func TestCreateWithNumericID(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateWithNumericIDSuccessfully(t)

	actual, err := floatingips.Create(context.TODO(), client.ServiceClient(), floatingips.CreateOpts{
		Pool: "nova",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedFloatingIP, actual)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := floatingips.Get(context.TODO(), client.ServiceClient(), "2").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &SecondFloatingIP, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	err := floatingips.Delete(context.TODO(), client.ServiceClient(), "1").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAssociate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAssociateSuccessfully(t)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: "10.10.10.2",
	}

	err := floatingips.AssociateInstance(context.TODO(), client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", associateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAssociateFixed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAssociateFixedSuccessfully(t)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: "10.10.10.2",
		FixedIP:    "166.78.185.201",
	}

	err := floatingips.AssociateInstance(context.TODO(), client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", associateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDisassociateInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDisassociateSuccessfully(t)

	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: "10.10.10.2",
	}

	err := floatingips.DisassociateInstance(context.TODO(), client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", disassociateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
