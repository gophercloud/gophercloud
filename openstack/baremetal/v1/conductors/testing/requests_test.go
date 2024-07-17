package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/conductors"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListConductors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConductorListSuccessfully(t, fakeServer)

	pages := 0
	err := conductors.List(client.ServiceClient(fakeServer), conductors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := conductors.ExtractConductors(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 conductors, got %d", len(actual))
		}
		th.AssertEquals(t, "compute1.localdomain", actual[0].Hostname)
		th.AssertEquals(t, "compute2.localdomain", actual[1].Hostname)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListDetailConductors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConductorListDetailSuccessfully(t, fakeServer)

	pages := 0
	err := conductors.List(client.ServiceClient(fakeServer), conductors.ListOpts{Detail: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := conductors.ExtractConductors(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 conductors, got %d", len(actual))
		}
		th.AssertEquals(t, "compute1.localdomain", actual[0].Hostname)
		th.AssertEquals(t, false, actual[0].Alive)
		th.AssertEquals(t, "compute2.localdomain", actual[1].Hostname)
		th.AssertEquals(t, true, actual[1].Alive)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListOpts(t *testing.T) {
	// Detail cannot take Fields
	optsDetail := conductors.ListOpts{
		Fields: []string{"hostname", "alive"},
		Detail: true,
	}

	opts := conductors.ListOpts{
		Fields: []string{"hostname", "alive"},
	}

	_, err := optsDetail.ToConductorListQuery()
	th.AssertEquals(t, err.Error(), "cannot have both fields and detail options for conductors")

	// Regular ListOpts can
	query, err := opts.ToConductorListQuery()
	th.AssertEquals(t, "?fields=hostname%2Calive", query)
	th.AssertNoErr(t, err)
}

func TestGetConductor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConductorGetSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := conductors.Get(context.TODO(), c, "1234asdf").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, ConductorFoo, *actual)
}
