//go:build acceptance
// +build acceptance

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/federation"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestListMappings(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := federation.ListMappings(client).AllPages()
	th.AssertNoErr(t, err)

	mappings, err := federation.ExtractMappings(allPages)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, mappings)
}

func TestMappingsCRUD(t *testing.T) {
	mappingName := tools.RandomString("TESTMAPPING-", 8)

	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

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

	actual, err := federation.CreateMapping(client, mappingName, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, createOpts.Rules[0], actual.Rules[0])
}
