//go:build acceptance || identity || federation

package v3

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListMappings(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := federation.ListMappings(client).AllPages(context.TODO())
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

	createdMapping, err := federation.CreateMapping(context.TODO(), client, mappingName, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(createOpts.Rules), len(createdMapping.Rules))
	th.CheckDeepEquals(t, createOpts.Rules[0], createdMapping.Rules[0])

	mapping, err := federation.GetMapping(context.TODO(), client, mappingName).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(createOpts.Rules), len(mapping.Rules))
	th.CheckDeepEquals(t, createOpts.Rules[0], mapping.Rules[0])

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

	updatedMapping, err := federation.UpdateMapping(context.TODO(), client, mappingName, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(updateOpts.Rules), len(updatedMapping.Rules))
	th.CheckDeepEquals(t, updateOpts.Rules[0], updatedMapping.Rules[0])

	err = federation.DeleteMapping(context.TODO(), client, mappingName).ExtractErr()
	th.AssertNoErr(t, err)

	resp := federation.GetMapping(context.TODO(), client, mappingName)
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(resp.Err, http.StatusNotFound))
}
