package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/drivers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListDrivers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListDriversSuccessfully(t)

	pages := 0
	err := drivers.ListDrivers(client.ServiceClient(), drivers.ListDriversOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := drivers.ExtractDrivers(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 3 {
			t.Fatalf("Expected 3 drivers, got %d", len(actual))
		}

		th.CheckDeepEquals(t, DriverAgentIpmitool, actual[0])
		th.CheckDeepEquals(t, DriverFake, actual[1])
		th.AssertEquals(t, "ipmi", actual[2].Name)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestGetDriverDetails(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetDriverDetailsSuccessfully(t)

	c := client.ServiceClient()
	actual, err := drivers.GetDriverDetails(context.TODO(), c, "ipmi").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, DriverIpmi, *actual)
}

func TestGetDriverProperties(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetDriverPropertiesSuccessfully(t)

	c := client.ServiceClient()
	actual, err := drivers.GetDriverProperties(context.TODO(), c, "agent_ipmitool").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, DriverIpmiToolProperties, *actual)
}

func TestGetDriverDiskProperties(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetDriverDiskPropertiesSuccessfully(t)

	c := client.ServiceClient()
	actual, err := drivers.GetDriverDiskProperties(context.TODO(), c, "agent_ipmitool").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, DriverIpmiToolDisk, *actual)
}
