package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	count := 0
	err := floatingips.List(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
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

	actual, err := floatingips.Create(client.ServiceClient(), floatingips.CreateOpts{
		Pool: "nova",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedFloatingIP, actual)
}

func TestBadCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	_, err := floatingips.Create(client.ServiceClient(), floatingips.CreateOpts{}).Extract()
	if err == nil {
		t.Fatalf("Expected an error due to missing Pool")
	}
}

func TestCreateWithNumericID(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateWithNumericIDSuccessfully(t)

	actual, err := floatingips.Create(client.ServiceClient(), floatingips.CreateOpts{
		Pool: "nova",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedFloatingIP, actual)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := floatingips.Get(client.ServiceClient(), "2").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &SecondFloatingIP, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	err := floatingips.Delete(client.ServiceClient(), "1").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAssociate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAssociateSuccessfully(t)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: "10.10.10.2",
	}

	err := floatingips.AssociateInstance(client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", associateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestBadAssociate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAssociateSuccessfully(t)

	associateOpts := floatingips.AssociateOpts{}
	err := floatingips.AssociateInstance(client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", associateOpts).ExtractErr()
	if err == nil {
		t.Fatalf("Expected an error due to missing FloatingIP")
	}
}

func TestAssociateFixed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAssociateFixedSuccessfully(t)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: "10.10.10.2",
		FixedIP:    "166.78.185.201",
	}

	err := floatingips.AssociateInstance(client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", associateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDisassociateInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDisassociateSuccessfully(t)

	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: "10.10.10.2",
	}

	err := floatingips.DisassociateInstance(client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", disassociateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestBadDisassociateInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDisassociateSuccessfully(t)

	disassociateOpts := floatingips.DisassociateOpts{}
	err := floatingips.DisassociateInstance(client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0", disassociateOpts).ExtractErr()
	if err == nil {
		t.Fatalf("Expected an error due to missing Floating IP")
	}
}

func TestBadGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleBadGetSuccessfully(t)

	_, err := floatingips.Get(client.ServiceClient(), "2").Extract()
	if err == nil {
		t.Fatalf("Expected an unmarshal error")
	}
}
