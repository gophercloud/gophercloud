package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/aggregates"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAggregates(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	pages := 0
	err := aggregates.List(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := aggregates.ExtractAggregates(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 aggregates, got %d", len(actual))
		}
		th.CheckDeepEquals(t, FirstFakeAggregate, actual[0])
		th.CheckDeepEquals(t, SecondFakeAggregate, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestCreateAggregates(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer)

	expected := CreatedAggregate

	opts := aggregates.CreateOpts{
		Name:             "name",
		AvailabilityZone: "london",
	}

	actual, err := aggregates.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestDeleteAggregates(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	err := aggregates.Delete(context.TODO(), client.ServiceClient(fakeServer), AggregateIDtoDelete).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetAggregates(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)

	expected := SecondFakeAggregate

	actual, err := aggregates.Get(context.TODO(), client.ServiceClient(fakeServer), AggregateIDtoGet).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestUpdateAggregate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	expected := UpdatedAggregate
	opts := aggregates.UpdateOpts{
		Name:             ptr.To("test-aggregates2"),
		AvailabilityZone: ptr.To("nova2"),
	}

	actual, err := aggregates.Update(context.TODO(), client.ServiceClient(fakeServer), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestAddHostAggregate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAddHostSuccessfully(t, fakeServer)

	expected := AggregateWithAddedHost

	opts := aggregates.AddHostOpts{
		Host: "cmp1",
	}

	actual, err := aggregates.AddHost(context.TODO(), client.ServiceClient(fakeServer), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestRemoveHostAggregate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRemoveHostSuccessfully(t, fakeServer)

	expected := AggregateWithRemovedHost

	opts := aggregates.RemoveHostOpts{
		Host: "cmp1",
	}

	actual, err := aggregates.RemoveHost(context.TODO(), client.ServiceClient(fakeServer), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestSetMetadataAggregate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSetMetadataSuccessfully(t, fakeServer)

	expected := AggregateWithUpdatedMetadata

	opts := aggregates.SetMetadataOpts{
		Metadata: map[string]any{"key": "value"},
	}

	actual, err := aggregates.SetMetadata(context.TODO(), client.ServiceClient(fakeServer), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}
