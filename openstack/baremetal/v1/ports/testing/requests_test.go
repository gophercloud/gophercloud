package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListDetailPorts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePortListDetailSuccessfully(t, fakeServer)

	pages := 0
	err := ports.ListDetail(client.ServiceClient(fakeServer), ports.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := ports.ExtractPorts(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 ports, got %d", len(actual))
		}
		th.CheckDeepEquals(t, PortBar, actual[0])
		th.CheckDeepEquals(t, PortFoo, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListPorts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePortListSuccessfully(t, fakeServer)

	pages := 0
	err := ports.List(client.ServiceClient(fakeServer), ports.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := ports.ExtractPorts(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 ports, got %d", len(actual))
		}
		th.AssertEquals(t, "3abe3f36-9708-4e9f-b07e-0f898061d3a7", actual[0].UUID)
		th.AssertEquals(t, "f2845e11-dbd4-4728-a8c0-30d19f48924a", actual[1].UUID)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListOpts(t *testing.T) {
	// Detail cannot take Fields
	opts := ports.ListOpts{
		Fields: []string{"uuid", "address"},
	}

	_, err := opts.ToPortListDetailQuery()
	th.AssertEquals(t, err.Error(), "fields is not a valid option when getting a detailed listing of ports")

	// Regular ListOpts can
	query, err := opts.ToPortListQuery()
	th.AssertEquals(t, "?fields=uuid%2Caddress", query)
	th.AssertNoErr(t, err)
}

func TestCreatePort(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePortCreationSuccessfully(t, fakeServer, SinglePortBody)

	iTrue := true
	actual, err := ports.Create(context.TODO(), client.ServiceClient(fakeServer), ports.CreateOpts{
		NodeUUID:   "ddd06a60-b91e-4ab4-a6e7-56c0b25b6086",
		Address:    "52:54:00:4d:87:e6",
		PXEEnabled: &iTrue,
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, PortFoo, *actual)
}

func TestDeletePort(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePortDeletionSuccessfully(t, fakeServer)

	res := ports.Delete(context.TODO(), client.ServiceClient(fakeServer), "3abe3f36-9708-4e9f-b07e-0f898061d3a7")
	th.AssertNoErr(t, res.Err)
}

func TestGetPort(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePortGetSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := ports.Get(context.TODO(), c, "f2845e11-dbd4-4728-a8c0-30d19f48924a").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, PortFoo, *actual)
}

func TestUpdatePort(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePortUpdateSuccessfully(t, fakeServer, SinglePortBody)

	c := client.ServiceClient(fakeServer)
	actual, err := ports.Update(context.TODO(), c, "f2845e11-dbd4-4728-a8c0-30d19f48924a", ports.UpdateOpts{
		ports.UpdateOperation{
			Op:    ports.ReplaceOp,
			Path:  "/address",
			Value: "22:22:22:22:22:22",
		},
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, PortFoo, *actual)
}
