package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceclasses"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListResourceClasses(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListResourceClasses(t, fakeServer)

	allPages, err := resourceclasses.List(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := resourceclasses.ExtractResourceClasses(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedResourceClassesList, actual)
}

func TestGetResourceClassSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetResourceClassSuccess(t, fakeServer)

	actual, err := resourceclasses.Get(context.TODO(), client.ServiceClient(fakeServer), PresentResourceClass).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedResourceClass, actual)
}

func TestGetResourceClassNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetResourceClassNotFound(t, fakeServer)

	_, err := resourceclasses.Get(context.TODO(), client.ServiceClient(fakeServer), AbsentResourceClass).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestCreateResourceClassSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateResourceClassSuccess(t, fakeServer)

	createOpts := resourceclasses.CreateOpts{
		Name: NewResourceClass,
	}

	err := resourceclasses.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestCreateResourceClassConflict(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateResourceClassConflict(t, fakeServer)

	createOpts := resourceclasses.CreateOpts{
		Name: PresentResourceClass,
	}

	err := resourceclasses.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

func TestUpdateResourceClassCreateSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateResourceClassSuccess(t, fakeServer)

	err := resourceclasses.Update(context.TODO(), client.ServiceClient(fakeServer), NewResourceClass).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdateResourceClassExists(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateResourceClassExists(t, fakeServer)

	err := resourceclasses.Update(context.TODO(), client.ServiceClient(fakeServer), PresentResourceClass).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdateResourceClassNonCustom(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateResourceClassNonCustom(t, fakeServer)

	err := resourceclasses.Update(context.TODO(), client.ServiceClient(fakeServer), "VCPU").ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}
