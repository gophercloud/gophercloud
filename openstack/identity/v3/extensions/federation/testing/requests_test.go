package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/federation"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListMappings(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListMappingsSuccessfully(t)

	count := 0
	err := federation.ListMappings(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
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

	allPages, err := federation.ListMappings(client.ServiceClient()).AllPages()
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

	actual, err := federation.CreateMapping(client.ServiceClient(), "ACME", createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, MappingACME, *actual)
}

func TestGetMapping(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetMappingSuccessfully(t)

	actual, err := federation.GetMapping(client.ServiceClient(), "ACME").Extract()
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

	actual, err := federation.UpdateMapping(client.ServiceClient(), "ACME", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, MappingUpdated, *actual)
}
