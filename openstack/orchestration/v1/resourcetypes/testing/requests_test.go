package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/resourcetypes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestBasicListResourceTypes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	result := resourcetypes.List(context.TODO(), client.ServiceClient(fakeServer), nil)
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, BasicListExpected, actual)
}

func TestFullListResourceTypes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	result := resourcetypes.List(context.TODO(), client.ServiceClient(fakeServer), resourcetypes.ListOpts{
		WithDescription: true,
	})
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, FullListExpected, actual)
}

func TestFilteredListResourceTypes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	result := resourcetypes.List(context.TODO(), client.ServiceClient(fakeServer), resourcetypes.ListOpts{
		NameRegex:       listFilterRegex,
		WithDescription: true,
	})
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, FilteredListExpected, actual)
}

func TestGetSchema(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSchemaSuccessfully(t, fakeServer)

	result := resourcetypes.GetSchema(context.TODO(), client.ServiceClient(fakeServer), "OS::Test::TestServer")
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, GetSchemaExpected, actual)
}

func TestGenerateTemplate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGenerateTemplateSuccessfully(t, fakeServer)

	result := resourcetypes.GenerateTemplate(context.TODO(), client.ServiceClient(fakeServer), "OS::Heat::None", nil)
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "2016-10-14", actual["heat_template_version"])
}
