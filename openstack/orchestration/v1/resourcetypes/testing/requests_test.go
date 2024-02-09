package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/resourcetypes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestBasicListResourceTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	result := resourcetypes.List(context.TODO(), fake.ServiceClient(), nil)
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, BasicListExpected, actual)
}

func TestFullListResourceTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	result := resourcetypes.List(context.TODO(), fake.ServiceClient(), resourcetypes.ListOpts{
		WithDescription: true,
	})
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, FullListExpected, actual)
}

func TestFilteredListResourceTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	result := resourcetypes.List(context.TODO(), fake.ServiceClient(), resourcetypes.ListOpts{
		NameRegex:       listFilterRegex,
		WithDescription: true,
	})
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, FilteredListExpected, actual)
}

func TestGetSchema(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSchemaSuccessfully(t)

	result := resourcetypes.GetSchema(context.TODO(), fake.ServiceClient(), "OS::Test::TestServer")
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, GetSchemaExpected, actual)
}

func TestGenerateTemplate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGenerateTemplateSuccessfully(t)

	result := resourcetypes.GenerateTemplate(context.TODO(), fake.ServiceClient(), "OS::Heat::None", nil)
	th.AssertNoErr(t, result.Err)

	actual, err := result.Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "2016-10-14", actual["heat_template_version"])
}
