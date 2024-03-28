package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/services"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListServicesPre253(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	HandleListPre253Successfully(t)

	pages := 0
	err := services.List(client.ServiceClient(), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 4 {
			t.Fatalf("Expected 4 services, got %d", len(actual))
		}
		testhelper.CheckDeepEquals(t, FirstFakeServicePre253, actual[0])
		testhelper.CheckDeepEquals(t, SecondFakeServicePre253, actual[1])
		testhelper.CheckDeepEquals(t, ThirdFakeServicePre253, actual[2])
		testhelper.CheckDeepEquals(t, FourthFakeServicePre253, actual[3])

		return true, nil
	})

	testhelper.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListServices(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	HandleListSuccessfully(t)

	pages := 0
	opts := services.ListOpts{
		Binary: "fake-binary",
		Host:   "host123",
	}
	err := services.List(client.ServiceClient(), opts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 4 {
			t.Fatalf("Expected 4 services, got %d", len(actual))
		}
		testhelper.CheckDeepEquals(t, FirstFakeService, actual[0])
		testhelper.CheckDeepEquals(t, SecondFakeService, actual[1])
		testhelper.CheckDeepEquals(t, ThirdFakeService, actual[2])
		testhelper.CheckDeepEquals(t, FourthFakeService, actual[3])

		return true, nil
	})

	testhelper.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestUpdateService(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	HandleUpdateSuccessfully(t)

	client := client.ServiceClient()
	actual, err := services.Update(context.TODO(), client, "fake-service-id", services.UpdateOpts{Status: services.ServiceDisabled}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	testhelper.CheckDeepEquals(t, FakeServiceUpdateBody, *actual)
}

func TestDeleteService(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	client := client.ServiceClient()
	res := services.Delete(context.TODO(), client, "fake-service-id")

	testhelper.AssertNoErr(t, res.Err)
}
