package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/aggregates"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAggregates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	pages := 0
	err := aggregates.List(client.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	expected := CreatedAggregate

	opts := aggregates.CreateOpts{
		Name:             "name",
		AvailabilityZone: "london",
	}

	actual, err := aggregates.Create(context.TODO(), client.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestDeleteAggregates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	err := aggregates.Delete(context.TODO(), client.ServiceClient(), AggregateIDtoDelete).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetAggregates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	expected := SecondFakeAggregate

	actual, err := aggregates.Get(context.TODO(), client.ServiceClient(), AggregateIDtoGet).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestUpdateAggregate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateSuccessfully(t)

	expected := UpdatedAggregate

	opts := aggregates.UpdateOpts{
		Name:             "test-aggregates2",
		AvailabilityZone: "nova2",
	}

	actual, err := aggregates.Update(context.TODO(), client.ServiceClient(), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestAddHostAggregate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAddHostSuccessfully(t)

	expected := AggregateWithAddedHost

	opts := aggregates.AddHostOpts{
		Host: "cmp1",
	}

	actual, err := aggregates.AddHost(context.TODO(), client.ServiceClient(), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestRemoveHostAggregate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleRemoveHostSuccessfully(t)

	expected := AggregateWithRemovedHost

	opts := aggregates.RemoveHostOpts{
		Host: "cmp1",
	}

	actual, err := aggregates.RemoveHost(context.TODO(), client.ServiceClient(), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestSetMetadataAggregate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSetMetadataSuccessfully(t)

	expected := AggregateWithUpdatedMetadata

	opts := aggregates.SetMetadataOpts{
		Metadata: map[string]any{"key": "value"},
	}

	actual, err := aggregates.SetMetadata(context.TODO(), client.ServiceClient(), expected.ID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}
