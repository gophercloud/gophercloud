package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListMappings(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListMappingsSuccessfully(t)

	count := 0
	err := federation.ListMappings(client.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := federation.ExtractMappings(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedMappingsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListMappingsAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListMappingsSuccessfully(t)

	allPages, err := federation.ListMappings(client.ServiceClient()).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := federation.ExtractMappings(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedMappingsSlice, actual)
}

func TestCreateMappings(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateMappingSuccessfully(t)

	createOpts := federation.CreateMappingOpts{
		Rules: []federation.MappingRule{
			{
				Local: []federation.RuleLocal{
					{
						User: &federation.RuleUser{
							Name: "{0}",
						},
					},
					{
						Group: &federation.Group{
							ID: "0cd5e9",
						},
					},
				},
				Remote: []federation.RuleRemote{
					{
						Type: "UserName",
					},
					{
						Type: "orgPersonType",
						NotAnyOf: []string{
							"Contractor",
							"Guest",
						},
					},
				},
			},
		},
	}

	actual, err := federation.CreateMapping(context.TODO(), client.ServiceClient(), "ACME", createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, MappingACME, *actual)
}

func TestGetMapping(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetMappingSuccessfully(t)

	actual, err := federation.GetMapping(context.TODO(), client.ServiceClient(), "ACME").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, MappingACME, *actual)
}

func TestUpdateMapping(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateMappingSuccessfully(t)

	updateOpts := federation.UpdateMappingOpts{
		Rules: []federation.MappingRule{
			{
				Local: []federation.RuleLocal{
					{
						User: &federation.RuleUser{
							Name: "{0}",
						},
					},
					{
						Group: &federation.Group{
							ID: "0cd5e9",
						},
					},
				},
				Remote: []federation.RuleRemote{
					{
						Type: "UserName",
					},
					{
						Type: "orgPersonType",
						AnyOneOf: []string{
							"Contractor",
							"SubContractor",
						},
					},
				},
			},
		},
	}

	actual, err := federation.UpdateMapping(context.TODO(), client.ServiceClient(), "ACME", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, MappingUpdated, *actual)
}

func TestDeleteMapping(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteMappingSuccessfully(t)

	res := federation.DeleteMapping(context.TODO(), client.ServiceClient(), "ACME")
	th.AssertNoErr(t, res.Err)
}
