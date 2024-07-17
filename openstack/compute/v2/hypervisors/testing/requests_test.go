package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/hypervisors"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListHypervisorsPre253(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorListPre253Successfully(t, fakeServer)

	pages := 0
	err := hypervisors.List(client.ServiceClient(fakeServer),
		hypervisors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := hypervisors.ExtractHypervisors(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 hypervisors, got %d", len(actual))
		}
		th.CheckDeepEquals(t, HypervisorFakePre253, actual[0])
		th.CheckDeepEquals(t, HypervisorFakePre253, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllHypervisorsPre253(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorListPre253Successfully(t, fakeServer)

	allPages, err := hypervisors.List(client.ServiceClient(fakeServer), hypervisors.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, HypervisorFakePre253, actual[0])
	th.CheckDeepEquals(t, HypervisorFakePre253, actual[1])
}

func TestListHypervisors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorListSuccessfully(t, fakeServer)

	pages := 0
	err := hypervisors.List(client.ServiceClient(fakeServer),
		hypervisors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := hypervisors.ExtractHypervisors(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 1 {
			t.Fatalf("Expected 1 hypervisors on page %d, got %d", pages, len(actual))
		}
		th.CheckDeepEquals(t, HypervisorFake, actual[0])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 2 {
		t.Errorf("Expected 2 pages, saw %d", pages)
	}
}

func TestListAllHypervisors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorListSuccessfully(t, fakeServer)

	allPages, err := hypervisors.List(client.ServiceClient(fakeServer), hypervisors.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, HypervisorFake, actual[0])
	th.CheckDeepEquals(t, HypervisorFake, actual[1])
}

func TestListAllHypervisorsWithParameters(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorListWithParametersSuccessfully(t, fakeServer)

	with_servers := true
	allPages, err := hypervisors.List(client.ServiceClient(fakeServer), hypervisors.ListOpts{WithServers: &with_servers}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, HypervisorFakeWithParameters, actual[0])
	th.CheckDeepEquals(t, HypervisorFakeWithParameters, actual[1])
}

func TestHypervisorsStatistics(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorsStatisticsSuccessfully(t, fakeServer)

	expected := HypervisorsStatisticsExpected

	actual, err := hypervisors.GetStatistics(context.TODO(), client.ServiceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

func TestGetHypervisor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorGetSuccessfully(t, fakeServer)

	expected := HypervisorFake

	actual, err := hypervisors.Get(context.TODO(), client.ServiceClient(fakeServer), expected.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

func TestGetHypervisorEmptyCPUInfo(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorGetEmptyCPUInfoSuccessfully(t, fakeServer)

	expected := HypervisorEmptyCPUInfo

	actual, err := hypervisors.Get(context.TODO(), client.ServiceClient(fakeServer), expected.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

func TestGetHypervisorAfterV287Response(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorAfterV287ResponseSuccessfully(t, fakeServer)

	expected := HypervisorAfterV287Response

	actual, err := hypervisors.Get(context.TODO(), client.ServiceClient(fakeServer), expected.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

func TestHypervisorsUptime(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHypervisorUptimeSuccessfully(t, fakeServer)

	expected := HypervisorUptimeExpected

	actual, err := hypervisors.GetUptime(context.TODO(), client.ServiceClient(fakeServer), HypervisorFake.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}
