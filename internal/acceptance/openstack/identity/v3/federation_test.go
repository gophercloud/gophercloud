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
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := federation.ListMappings(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	mappings, err := federation.ExtractMappings(allPages)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, mappings)
}

func TestMappingsCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	mappingName := tools.RandomString("TESTMAPPING-", 8)

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
	th.AssertTrue(t, gophercloud.ResponseCodeIs(resp.Err, http.StatusNotFound))
}

func TestIdentityProvidersCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	identityProviderID := tools.RandomString("TESTIDP-", 8)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := federation.CreateIdentityProviderOpts{
		Description: "Gophercloud acceptance test identity provider",
		Enabled:     gophercloud.Enabled,
		RemoteIDs:   []string{tools.RandomString("remote-", 8)},
	}
	created, err := federation.CreateIdentityProvider(context.TODO(), client, identityProviderID, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer federation.DeleteIdentityProvider(context.TODO(), client, identityProviderID)

	th.CheckEquals(t, identityProviderID, created.ID)
	th.AssertTrue(t, created.Description != nil)
	th.CheckEquals(t, createOpts.Description, *created.Description)
	th.CheckEquals(t, true, created.Enabled)
	th.CheckDeepEquals(t, createOpts.RemoteIDs, created.RemoteIDs)

	actual, err := federation.GetIdentityProvider(context.TODO(), client, identityProviderID).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, identityProviderID, actual.ID)

	listOpts := federation.ListIdentityProvidersOpts{
		ID:      identityProviderID,
		Enabled: gophercloud.Enabled,
	}
	allPages, err := federation.ListIdentityProviders(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	identityProviders, err := federation.ExtractIdentityProviders(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(identityProviders))
	th.CheckEquals(t, identityProviderID, identityProviders[0].ID)
	tools.PrintResource(t, identityProviders)

	description := "Updated Gophercloud acceptance test identity provider"
	remoteIDs := []string{tools.RandomString("updated-remote-", 8)}
	updateOpts := federation.UpdateIdentityProviderOpts{
		Description: &description,
		Enabled:     gophercloud.Disabled,
		RemoteIDs:   &remoteIDs,
	}
	updated, err := federation.UpdateIdentityProvider(context.TODO(), client, identityProviderID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertTrue(t, updated.Description != nil)
	th.CheckEquals(t, description, *updated.Description)
	th.CheckEquals(t, false, updated.Enabled)
	th.CheckDeepEquals(t, remoteIDs, updated.RemoteIDs)

	err = federation.DeleteIdentityProvider(context.TODO(), client, identityProviderID).ExtractErr()
	th.AssertNoErr(t, err)

	resp := federation.GetIdentityProvider(context.TODO(), client, identityProviderID)
	th.AssertTrue(t, gophercloud.ResponseCodeIs(resp.Err, http.StatusNotFound))
}

func TestProtocolsCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	identityProviderID := tools.RandomString("TESTIDP-", 8)
	mappingID := tools.RandomString("TESTMAPPING-", 8)
	updatedMappingID := tools.RandomString("TESTMAPPING-", 8)
	protocolID := tools.RandomString("testprotocol-", 8)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	mappingOpts := federation.CreateMappingOpts{
		Rules: []federation.MappingRule{
			{
				Local: []federation.RuleLocal{
					{User: &federation.RuleUser{Name: "{0}"}},
				},
				Remote: []federation.RuleRemote{
					{Type: "UserName"},
				},
			},
		},
	}
	_, err = federation.CreateMapping(context.TODO(), client, mappingID, mappingOpts).Extract()
	th.AssertNoErr(t, err)
	defer federation.DeleteMapping(context.TODO(), client, mappingID)

	_, err = federation.CreateMapping(context.TODO(), client, updatedMappingID, mappingOpts).Extract()
	th.AssertNoErr(t, err)
	defer federation.DeleteMapping(context.TODO(), client, updatedMappingID)

	_, err = federation.CreateIdentityProvider(context.TODO(), client, identityProviderID, federation.CreateIdentityProviderOpts{
		Enabled: gophercloud.Enabled,
	}).Extract()
	th.AssertNoErr(t, err)
	defer federation.DeleteIdentityProvider(context.TODO(), client, identityProviderID)

	createOpts := federation.CreateProtocolOpts{
		MappingID:         mappingID,
		RemoteIDAttribute: "HTTP_OIDC_ISS",
	}
	created, err := federation.CreateProtocol(context.TODO(), client, identityProviderID, protocolID, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer federation.DeleteProtocol(context.TODO(), client, identityProviderID, protocolID)

	th.CheckEquals(t, protocolID, created.ID)
	th.CheckEquals(t, mappingID, created.MappingID)

	actual, err := federation.GetProtocol(context.TODO(), client, identityProviderID, protocolID).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, protocolID, actual.ID)

	allPages, err := federation.ListProtocols(client, identityProviderID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	protocols, err := federation.ExtractProtocols(allPages)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, protocols)

	remoteIDAttribute := "HTTP_UPDATED_OIDC_ISS"
	updateOpts := federation.UpdateProtocolOpts{
		MappingID:         updatedMappingID,
		RemoteIDAttribute: &remoteIDAttribute,
	}
	updated, err := federation.UpdateProtocol(context.TODO(), client, identityProviderID, protocolID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, updatedMappingID, updated.MappingID)

	err = federation.DeleteProtocol(context.TODO(), client, identityProviderID, protocolID).ExtractErr()
	th.AssertNoErr(t, err)

	resp := federation.GetProtocol(context.TODO(), client, identityProviderID, protocolID)
	th.AssertTrue(t, gophercloud.ResponseCodeIs(resp.Err, http.StatusNotFound))
}

func TestServiceProvidersCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	serviceProviderID := tools.RandomString("TESTSP-", 8)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := federation.CreateServiceProviderOpts{
		AuthURL:          "https://example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth",
		Description:      "Gophercloud acceptance test service provider",
		Enabled:          gophercloud.Enabled,
		RelayStatePrefix: "ss:mem:",
		SPURL:            "https://example.com/Shibboleth.sso/SAML2/ECP",
	}
	created, err := federation.CreateServiceProvider(context.TODO(), client, serviceProviderID, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer federation.DeleteServiceProvider(context.TODO(), client, serviceProviderID)

	th.CheckEquals(t, serviceProviderID, created.ID)
	th.CheckEquals(t, createOpts.AuthURL, created.AuthURL)
	th.AssertTrue(t, created.Description != nil)
	th.CheckEquals(t, createOpts.Description, *created.Description)
	th.CheckEquals(t, true, created.Enabled)
	th.AssertTrue(t, created.RelayStatePrefix != nil)
	th.CheckEquals(t, createOpts.RelayStatePrefix, *created.RelayStatePrefix)
	th.CheckEquals(t, createOpts.SPURL, created.SPURL)

	actual, err := federation.GetServiceProvider(context.TODO(), client, serviceProviderID).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, serviceProviderID, actual.ID)

	listOpts := federation.ListServiceProvidersOpts{
		ID:      serviceProviderID,
		Enabled: gophercloud.Enabled,
	}
	allPages, err := federation.ListServiceProviders(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	serviceProviders, err := federation.ExtractServiceProviders(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(serviceProviders))
	th.CheckEquals(t, serviceProviderID, serviceProviders[0].ID)
	tools.PrintResource(t, serviceProviders)

	authURL := "https://new.example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth"
	description := "Updated Gophercloud acceptance test service provider"
	relayStatePrefix := "ss:temp:"
	spURL := "https://new.example.com/Shibboleth.sso/SAML2/ECP"
	updateOpts := federation.UpdateServiceProviderOpts{
		AuthURL:          &authURL,
		Description:      &description,
		Enabled:          gophercloud.Disabled,
		RelayStatePrefix: &relayStatePrefix,
		SPURL:            &spURL,
	}
	updated, err := federation.UpdateServiceProvider(context.TODO(), client, serviceProviderID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, authURL, updated.AuthURL)
	th.AssertTrue(t, updated.Description != nil)
	th.CheckEquals(t, description, *updated.Description)
	th.CheckEquals(t, false, updated.Enabled)
	th.AssertTrue(t, updated.RelayStatePrefix != nil)
	th.CheckEquals(t, relayStatePrefix, *updated.RelayStatePrefix)
	th.CheckEquals(t, spURL, updated.SPURL)

	err = federation.DeleteServiceProvider(context.TODO(), client, serviceProviderID).ExtractErr()
	th.AssertNoErr(t, err)

	resp := federation.GetServiceProvider(context.TODO(), client, serviceProviderID)
	th.AssertTrue(t, gophercloud.ResponseCodeIs(resp.Err, http.StatusNotFound))
}
