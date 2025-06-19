package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetalintrospection/v1/introspection"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListIntrospections(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListIntrospectionsSuccessfully(t, fakeServer)

	pages := 0
	err := introspection.ListIntrospections(client.ServiceClient(fakeServer), introspection.ListIntrospectionsOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := introspection.ExtractIntrospections(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 introspections, got %d", len(actual))
		}
		th.CheckDeepEquals(t, IntrospectionFoo, actual[0])
		th.CheckDeepEquals(t, IntrospectionBar, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestGetIntrospectionStatus(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetIntrospectionStatusSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := introspection.GetIntrospectionStatus(context.TODO(), c, "c244557e-899f-46fa-a1ff-5b2c6718616b").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, IntrospectionBar, *actual)
}

func TestStartIntrospection(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleStartIntrospectionSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := introspection.StartIntrospection(context.TODO(), c, "c244557e-899f-46fa-a1ff-5b2c6718616b", introspection.StartOpts{}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAbortIntrospection(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAbortIntrospectionSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := introspection.AbortIntrospection(context.TODO(), c, "c244557e-899f-46fa-a1ff-5b2c6718616b").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetIntrospectionData(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetIntrospectionDataSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := introspection.GetIntrospectionData(context.TODO(), c, "c244557e-899f-46fa-a1ff-5b2c6718616b").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, IntrospectionDataRes, *actual)
}

func TestReApplyIntrospection(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleReApplyIntrospectionSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := introspection.ReApplyIntrospection(context.TODO(), c, "c244557e-899f-46fa-a1ff-5b2c6718616b").ExtractErr()
	th.AssertNoErr(t, err)
}
