package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stackevents"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestFindEvents(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFindSuccessfully(t, fakeServer, FindOutput)

	actual, err := stackevents.Find(context.TODO(), client.ServiceClient(fakeServer), "postman_stack").Extract()
	th.AssertNoErr(t, err)

	expected := FindExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer, ListOutput)

	count := 0
	err := stackevents.List(client.ServiceClient(fakeServer), "hello_world", "49181cd6-169a-4130-9455-31185bbfc5bf", nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := stackevents.ExtractEvents(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ListExpected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListResourceEvents(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListResourceEventsSuccessfully(t, fakeServer, ListResourceEventsOutput)

	count := 0
	err := stackevents.ListResourceEvents(client.ServiceClient(fakeServer), "hello_world", "49181cd6-169a-4130-9455-31185bbfc5bf", "my_resource", nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := stackevents.ExtractResourceEvents(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ListResourceEventsExpected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestGetEvent(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer, GetOutput)

	actual, err := stackevents.Get(context.TODO(), client.ServiceClient(fakeServer), "hello_world", "49181cd6-169a-4130-9455-31185bbfc5bf", "my_resource", "93940999-7d40-44ae-8de4-19624e7b8d18").Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, actual)
}
