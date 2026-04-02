package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"

	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListResourceProviders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderList(t, fakeServer)

	count := 0
	err := resourceproviders.List(client.ServiceClient(fakeServer), resourceproviders.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := resourceproviders.ExtractResourceProviders(page)
		if err != nil {
			t.Errorf("Failed to extract resource providers: %v", err)
			return false, err
		}
		th.AssertDeepEquals(t, ExpectedResourceProviders, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreateResourceProvider(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderCreate(t, fakeServer)

	expected := ExpectedResourceProvider1

	opts := resourceproviders.CreateOpts{
		Name:               ExpectedResourceProvider1.Name,
		UUID:               ExpectedResourceProvider1.UUID,
		ParentProviderUUID: ExpectedResourceProvider1.ParentProviderUUID,
	}

	actual, err := resourceproviders.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestGetResourceProvider(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGet(t, fakeServer)

	expected := ExpectedResourceProvider1

	actual, err := resourceproviders.Get(context.TODO(), client.ServiceClient(fakeServer), ExpectedResourceProvider1.UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &expected, actual)
}

func TestDeleteResourceProvider(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderDelete(t, fakeServer)

	res := resourceproviders.Delete(context.TODO(), client.ServiceClient(fakeServer), "b99b3ab4-3aa6-4fba-b827-69b88b9c544a")
	th.AssertNoErr(t, res.Err)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderUpdate(t, fakeServer)

	name := "new_name"
	parentProviderUUID := "b99b3ab4-3aa6-4fba-b827-69b88b9c544a"

	options := resourceproviders.UpdateOpts{
		Name:               &name,
		ParentProviderUUID: &parentProviderUUID,
	}
	rp, err := resourceproviders.Update(context.TODO(), client.ServiceClient(fakeServer), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, rp.Name, name)
	th.AssertEquals(t, rp.ParentProviderUUID, parentProviderUUID)
}

func TestGetResourceProvidersUsages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetUsages(t, fakeServer)

	actual, err := resourceproviders.GetUsages(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedUsages, *actual)
}

func TestGetResourceProvidersInventories(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetInventories(t, fakeServer)

	actual, err := resourceproviders.GetInventories(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedInventories, *actual)
}

func TestGetResourceProviderInventory(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetInventory(t, fakeServer)

	actual, err := resourceproviders.GetInventory(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, PresentInventoryResourceClass).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedInventory, *actual)
}

func TestGetResourceProviderInventoryNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetInventoryNotFound(t, fakeServer)

	_, err := resourceproviders.GetInventory(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, MissingInventoryResourceClass).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestUpdateResourceProvidersInventories(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderPutInventories(t, fakeServer)

	opts := resourceproviders.UpdateInventoriesOpts(ExpectedInventories)
	actual, err := resourceproviders.UpdateInventories(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedInventories, *actual)
}

func TestUpdateResourceProvidersInventoriesNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderPutInventoriesNotFound(t, fakeServer)

	opts := resourceproviders.UpdateInventoriesOpts(ExpectedInventories)
	_, err := resourceproviders.UpdateInventories(context.TODO(), client.ServiceClient(fakeServer), NonExistentRPID, opts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestUpdateResourceProviderInventory(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderPutInventory(t, fakeServer)

	opts := resourceproviders.UpdateInventoryOpts(ExpectedInventory)
	actual, err := resourceproviders.UpdateInventory(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, PresentInventoryResourceClass, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedInventory, *actual)
}

func TestUpdateResourceProviderInventoryNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderPutInventoryNotFound(t, fakeServer)

	opts := resourceproviders.UpdateInventoryOpts(ExpectedInventory)
	_, err := resourceproviders.UpdateInventory(context.TODO(), client.ServiceClient(fakeServer), NonExistentRPID, PresentInventoryResourceClass, opts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestDeleteResourceProviderInventorySuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderDeleteInventorySuccess(t, fakeServer)

	err := resourceproviders.DeleteInventory(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, PresentInventoryResourceClass).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteResourceProviderInventoryInUse(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderDeleteInventoryInUse(t, fakeServer)

	err := resourceproviders.DeleteInventory(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, PresentInventoryResourceClass).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

func TestDeleteResourceProviderInventoriesSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderDeleteInventoriesSuccess(t, fakeServer)

	err := resourceproviders.DeleteInventories(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteResourceProviderInventoriesConflict(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderDeleteInventoriesConflict(t, fakeServer)

	err := resourceproviders.DeleteInventories(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

func TestGetResourceProvidersAllocations(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetAllocations(t, fakeServer)

	actual, err := resourceproviders.GetAllocations(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocations, *actual)
}

func TestGetResourceProvidersTraits(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetTraits(t, fakeServer)

	actual, err := resourceproviders.GetTraits(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedTraits, *actual)
}

func TestUpdateResourceProvidersTraits(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderPutTraits(t, fakeServer)

	opts := resourceproviders.UpdateTraitsOpts(ExpectedTraits)
	actual, err := resourceproviders.UpdateTraits(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedTraits, *actual)
}

func TestDeleteResourceProvidersTraits(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderDeleteTraits(t, fakeServer)

	err := resourceproviders.DeleteTraits(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetResourceProviderAggregatesSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetAggregatesSuccess(t, fakeServer)

	actual, err := resourceproviders.GetAggregates(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAggregates, *actual)
}

func TestGetResourceProviderAggregatesPreGenerationSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetAggregatesPreGenerationSuccess(t, fakeServer)

	actual, err := resourceproviders.GetAggregates(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAggregatesPreGeneration, *actual)
}

func TestGetResourceProviderAggregatesNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderGetAggregatesNotFound(t, fakeServer)

	_, err := resourceproviders.GetAggregates(context.TODO(), client.ServiceClient(fakeServer), AbsentResourceProviderID).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestUpdateResourceProviderAggregatesSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderUpdateAndGetAggregatesSuccess(t, fakeServer)

	updateOpts := resourceproviders.UpdateAggregatesOpts(ExpectedUpdatedAggregates)
	_, err := resourceproviders.UpdateAggregates(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	actual, err := resourceproviders.GetAggregates(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedUpdatedAggregates, *actual)
}

func TestUpdateResourceProviderAggregatesPreGenerationSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderUpdateAndGetAggregatesPreGenerationSuccess(t, fakeServer)

	updateOpts := resourceproviders.UpdateAggregatesOpts{
		Aggregates: ExpectedUpdatedAggregatesPreGeneration.Aggregates,
	}
	_, err := resourceproviders.UpdateAggregates(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	actual, err := resourceproviders.GetAggregates(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedUpdatedAggregatesPreGeneration, *actual)
}

func TestUpdateResourceProviderAggregatesConflict(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleResourceProviderUpdateAggregatesConflict(t, fakeServer)

	updateOpts := resourceproviders.UpdateAggregatesOpts(ExpectedUpdatedAggregates)
	_, err := resourceproviders.UpdateAggregates(context.TODO(), client.ServiceClient(fakeServer), ResourceProviderTestID, updateOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}
