package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListHosts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListHosts(t, fakeServer)

	allPages, err := hosts.List(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedHostsList, actual)

	// Blazar reports a null status and updated_at for an unmodified host.
	th.AssertEquals(t, "", actual[0].Status)
	th.AssertEquals(t, true, actual[0].UpdatedAt == nil)
}

// Blazar flattens extra capabilities into the host object.
func TestListHostsWithCapabilities(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListHostsWithCapabilities(t, fakeServer)

	allPages, err := hosts.List(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedHostsListWithCapabilities, actual)

	// The timestamps are tagged "-"
	_, ok := actual[0].ExtraCapabilities["created_at"]
	th.AssertEquals(t, false, ok)
	_, ok = actual[0].ExtraCapabilities["updated_at"]
	th.AssertEquals(t, false, ok)
}
