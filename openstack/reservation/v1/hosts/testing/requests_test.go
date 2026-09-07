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

	allPages, err := hosts.List(client.ServiceClient(fakeServer), hosts.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedHostsList, actual)

	// Blazar reports a null status and updated_at for an unmodified host.
	th.AssertEquals(t, "", actual[0].Status)
	th.AssertEquals(t, true, actual[0].UpdatedAt == nil)
}

func TestCreateHost(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateHost(t, fakeServer)

	createOpts := hosts.CreateOpts{
		Name: "compute-1.example.com",
		ExtraCapabilities: map[string]any{
			"gpu":  "a100",
			"rack": "b12",
		},
	}

	actual, err := hosts.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedHostWithCapabilities, actual)
}

// Blazar requires a name, so the request should not reach the server without one.
func TestCreateHostWithoutName(t *testing.T) {
	_, err := hosts.CreateOpts{}.ToHostCreateMap()
	th.AssertEquals(t, true, err != nil)
}

// An extra capability must not be able to overwrite a host field.
func TestCreateHostCapabilityCollision(t *testing.T) {
	createOpts := hosts.CreateOpts{
		Name:              "compute-1.example.com",
		ExtraCapabilities: map[string]any{"name": "somethingelse"},
	}

	_, err := createOpts.ToHostCreateMap()
	th.AssertEquals(t, true, err != nil)
}

func TestGetHost(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetHost(t, fakeServer)

	actual, err := hosts.Get(context.TODO(), client.ServiceClient(fakeServer), "18").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedHostWithCapabilities, actual)
}

// Blazar flattens extra capabilities into the host object.
func TestListHostsWithCapabilities(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListHostsWithCapabilities(t, fakeServer)

	allPages, err := hosts.List(client.ServiceClient(fakeServer), hosts.ListOpts{}).AllPages(context.TODO())
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
