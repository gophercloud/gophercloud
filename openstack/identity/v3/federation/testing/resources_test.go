package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestIdentityProviders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	const createOutput = `
{
	"identity_provider": {
		"authorization_ttl": 30,
		"domain_id": "1789d1",
		"description": "Stores ACME identities",
		"enabled": true,
		"id": "ACME",
		"links": {
			"protocols": "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols",
			"self": "http://example.com/OS-FEDERATION/identity_providers/ACME"
		},
		"remote_ids": ["acme_id_1", "acme_id_2"]
	}
}`
	const updateOutput = `
{
	"identity_provider": {
		"authorization_ttl": 15,
		"domain_id": "1789d1",
		"description": "Updated ACME identities",
		"enabled": false,
		"id": "ACME",
		"links": {
			"protocols": "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols",
			"self": "http://example.com/OS-FEDERATION/identity_providers/ACME"
		},
		"remote_ids": []
	}
}`

	fakeServer.Mux.HandleFunc("/OS-FEDERATION/identity_providers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"identity_providers":[%s],"links":{"next":null,"previous":null}}`,
			`{"authorization_ttl":null,"domain_id":"1789d1","description":null,"enabled":true,"id":"ACME","links":{"self":"http://example.com/OS-FEDERATION/identity_providers/ACME"},"remote_ids":["acme_id_1","acme_id_2"]}`)
	})

	fakeServer.Mux.HandleFunc("/OS-FEDERATION/identity_providers/ACME", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		switch r.Method {
		case http.MethodPut:
			th.TestJSONRequest(t, r, `{
				"identity_provider": {
					"authorization_ttl": 30,
					"domain_id": "1789d1",
					"description": "Stores ACME identities",
					"enabled": true,
					"remote_ids": ["acme_id_1", "acme_id_2"]
				}
			}`)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, createOutput)
		case http.MethodGet:
			th.TestHeader(t, r, "Accept", "application/json")
			fmt.Fprint(w, createOutput)
		case http.MethodPatch:
			th.TestJSONRequest(t, r, `{
				"identity_provider": {
					"authorization_ttl": 15,
					"description": "Updated ACME identities",
					"enabled": false,
					"remote_ids": []
				}
			}`)
			fmt.Fprint(w, updateOutput)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})

	allPages, err := federation.ListIdentityProviders(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	identityProviders, err := federation.ExtractIdentityProviders(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(identityProviders))
	th.CheckEquals(t, (*int)(nil), identityProviders[0].AuthorizationTTL)
	th.CheckEquals(t, "1789d1", identityProviders[0].DomainID)
	th.CheckEquals(t, (*string)(nil), identityProviders[0].Description)
	th.CheckEquals(t, true, identityProviders[0].Enabled)
	th.CheckEquals(t, "ACME", identityProviders[0].ID)
	th.CheckDeepEquals(t, map[string]any{
		"self": "http://example.com/OS-FEDERATION/identity_providers/ACME",
	}, identityProviders[0].Links)
	th.CheckDeepEquals(t, []string{"acme_id_1", "acme_id_2"}, identityProviders[0].RemoteIDs)

	createOpts := federation.CreateIdentityProviderOpts{
		AuthorizationTTL: gophercloud.IntToPointer(30),
		DomainID:         "1789d1",
		Description:      "Stores ACME identities",
		Enabled:          gophercloud.Enabled,
		RemoteIDs:        []string{"acme_id_1", "acme_id_2"},
	}
	created, err := federation.CreateIdentityProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME", createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "ACME", created.ID)
	th.CheckEquals(t, "1789d1", created.DomainID)
	th.AssertTrue(t, created.Description != nil)
	th.CheckEquals(t, "Stores ACME identities", *created.Description)
	th.CheckEquals(t, true, created.Enabled)
	th.CheckEquals(t, 30, *created.AuthorizationTTL)
	th.CheckDeepEquals(t, map[string]any{
		"protocols": "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols",
		"self":      "http://example.com/OS-FEDERATION/identity_providers/ACME",
	}, created.Links)
	th.CheckDeepEquals(t, []string{"acme_id_1", "acme_id_2"}, created.RemoteIDs)

	actual, err := federation.GetIdentityProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, created, actual)

	description := "Updated ACME identities"
	remoteIDs := []string{}
	updateOpts := federation.UpdateIdentityProviderOpts{
		AuthorizationTTL: gophercloud.IntToPointer(15),
		Description:      &description,
		Enabled:          gophercloud.Disabled,
		RemoteIDs:        &remoteIDs,
	}
	updated, err := federation.UpdateIdentityProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 15, *updated.AuthorizationTTL)
	th.AssertTrue(t, updated.Description != nil)
	th.CheckEquals(t, "Updated ACME identities", *updated.Description)
	th.CheckEquals(t, false, updated.Enabled)
	th.CheckDeepEquals(t, []string{}, updated.RemoteIDs)

	err = federation.DeleteIdentityProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestIdentityProviderUpdateOptsNullableFields(t *testing.T) {
	actual, err := (federation.UpdateIdentityProviderOpts{}).ToIdentityProviderUpdateMap()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `{"identity_provider": {}}`, actual)

	description := ""
	remoteIDs := []string{}
	opts := federation.UpdateIdentityProviderOpts{
		AuthorizationTTL: gophercloud.IntToPointer(-1),
		Description:      &description,
		Enabled:          gophercloud.Disabled,
		RemoteIDs:        &remoteIDs,
	}
	actual, err = opts.ToIdentityProviderUpdateMap()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `{
		"identity_provider": {
			"authorization_ttl": null,
			"description": null,
			"enabled": false,
			"remote_ids": []
		}
	}`, actual)

	opts = federation.UpdateIdentityProviderOpts{
		AuthorizationTTL: gophercloud.IntToPointer(0),
	}
	actual, err = opts.ToIdentityProviderUpdateMap()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `{
		"identity_provider": {
			"authorization_ttl": 0
		}
	}`, actual)
}

func TestProtocols(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	const createOutput = `
{
	"protocol": {
		"id": "saml2",
		"links": {
			"identity_provider": "http://example.com/OS-FEDERATION/identity_providers/ACME",
			"self": "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols/saml2"
		},
		"mapping_id": "mapping-1",
		"remote_id_attribute": "Shib-Identity-Provider"
	}
}`
	const updateOutput = `
{
	"protocol": {
		"id": "saml2",
		"links": {
			"identity_provider": "http://example.com/OS-FEDERATION/identity_providers/ACME",
			"self": "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols/saml2"
		},
		"mapping_id": "mapping-2",
		"remote_id_attribute": "HTTP_OIDC_ISS"
	}
}`

	fakeServer.Mux.HandleFunc("/OS-FEDERATION/identity_providers/ACME/protocols", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"protocols": [{
				"id": "saml2",
				"links": {
					"identity_provider": "http://example.com/OS-FEDERATION/identity_providers/ACME",
					"self": "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols/saml2"
				},
				"mapping_id": "mapping-1"
			}],
			"links": {"next": null, "previous": null}
		}`)
	})

	fakeServer.Mux.HandleFunc("/OS-FEDERATION/identity_providers/ACME/protocols/saml2", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		switch r.Method {
		case http.MethodPut:
			th.TestJSONRequest(t, r, `{
				"protocol": {
					"mapping_id": "mapping-1",
					"remote_id_attribute": "Shib-Identity-Provider"
				}
			}`)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, createOutput)
		case http.MethodGet:
			th.TestHeader(t, r, "Accept", "application/json")
			fmt.Fprint(w, createOutput)
		case http.MethodPatch:
			th.TestJSONRequest(t, r, `{
				"protocol": {
					"mapping_id": "mapping-2",
					"remote_id_attribute": "HTTP_OIDC_ISS"
				}
			}`)
			fmt.Fprint(w, updateOutput)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})

	allPages, err := federation.ListProtocols(client.ServiceClient(fakeServer), "ACME").AllPages(context.TODO())
	th.AssertNoErr(t, err)
	protocols, err := federation.ExtractProtocols(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(protocols))
	th.CheckEquals(t, "saml2", protocols[0].ID)
	th.CheckDeepEquals(t, map[string]any{
		"identity_provider": "http://example.com/OS-FEDERATION/identity_providers/ACME",
		"self":              "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols/saml2",
	}, protocols[0].Links)
	th.CheckEquals(t, "mapping-1", protocols[0].MappingID)
	th.CheckEquals(t, (*string)(nil), protocols[0].RemoteIDAttribute)

	createOpts := federation.CreateProtocolOpts{
		MappingID:         "mapping-1",
		RemoteIDAttribute: "Shib-Identity-Provider",
	}
	created, err := federation.CreateProtocol(context.TODO(), client.ServiceClient(fakeServer), "ACME", "saml2", createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "saml2", created.ID)
	th.CheckEquals(t, "mapping-1", created.MappingID)
	th.CheckDeepEquals(t, map[string]any{
		"identity_provider": "http://example.com/OS-FEDERATION/identity_providers/ACME",
		"self":              "http://example.com/OS-FEDERATION/identity_providers/ACME/protocols/saml2",
	}, created.Links)
	th.AssertTrue(t, created.RemoteIDAttribute != nil)
	th.CheckEquals(t, "Shib-Identity-Provider", *created.RemoteIDAttribute)

	actual, err := federation.GetProtocol(context.TODO(), client.ServiceClient(fakeServer), "ACME", "saml2").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, created, actual)

	remoteIDAttribute := "HTTP_OIDC_ISS"
	updateOpts := federation.UpdateProtocolOpts{
		MappingID:         "mapping-2",
		RemoteIDAttribute: &remoteIDAttribute,
	}
	updated, err := federation.UpdateProtocol(context.TODO(), client.ServiceClient(fakeServer), "ACME", "saml2", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "mapping-2", updated.MappingID)
	th.AssertTrue(t, updated.RemoteIDAttribute != nil)
	th.CheckEquals(t, "HTTP_OIDC_ISS", *updated.RemoteIDAttribute)

	err = federation.DeleteProtocol(context.TODO(), client.ServiceClient(fakeServer), "ACME", "saml2").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestProtocolCreateOptsValidation(t *testing.T) {
	_, err := (federation.CreateProtocolOpts{}).ToProtocolCreateMap()
	th.AssertErr(t, err)
}

func TestServiceProviders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	const createOutput = `
{
	"service_provider": {
		"auth_url": "https://example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth",
		"description": "Remote Service Provider",
		"enabled": true,
		"id": "ACME",
		"links": {
			"self": "https://example.com/v3/OS-FEDERATION/service_providers/ACME"
		},
		"relay_state_prefix": "ss:mem:",
		"sp_url": "https://example.com/Shibboleth.sso/SAML2/ECP"
	}
}`
	const updateOutput = `
{
	"service_provider": {
		"auth_url": "https://new.example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth",
		"description": "Updated Service Provider",
		"enabled": false,
		"id": "ACME",
		"links": {
			"self": "https://example.com/v3/OS-FEDERATION/service_providers/ACME"
		},
		"relay_state_prefix": "ss:temp:",
		"sp_url": "https://new.example.com/Shibboleth.sso/SAML2/ECP"
	}
}`

	fakeServer.Mux.HandleFunc("/OS-FEDERATION/service_providers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"service_providers": [{
				"auth_url": "https://example.com/auth",
				"description": null,
				"enabled": true,
				"id": "ACME",
				"links": {
					"self": "https://example.com/v3/OS-FEDERATION/service_providers/ACME"
				},
				"relay_state_prefix": null,
				"sp_url": "https://example.com/sp"
			}],
			"links": {"next": null, "previous": null}
		}`)
	})

	fakeServer.Mux.HandleFunc("/OS-FEDERATION/service_providers/ACME", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		switch r.Method {
		case http.MethodPut:
			th.TestJSONRequest(t, r, `{
				"service_provider": {
					"auth_url": "https://example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth",
					"description": "Remote Service Provider",
					"enabled": true,
					"relay_state_prefix": "ss:mem:",
					"sp_url": "https://example.com/Shibboleth.sso/SAML2/ECP"
				}
			}`)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, createOutput)
		case http.MethodGet:
			th.TestHeader(t, r, "Accept", "application/json")
			fmt.Fprint(w, createOutput)
		case http.MethodPatch:
			th.TestJSONRequest(t, r, `{
				"service_provider": {
					"auth_url": "https://new.example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth",
					"description": "Updated Service Provider",
					"enabled": false,
					"relay_state_prefix": "ss:temp:",
					"sp_url": "https://new.example.com/Shibboleth.sso/SAML2/ECP"
				}
			}`)
			fmt.Fprint(w, updateOutput)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})

	allPages, err := federation.ListServiceProviders(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	serviceProviders, err := federation.ExtractServiceProviders(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(serviceProviders))
	th.CheckEquals(t, "https://example.com/auth", serviceProviders[0].AuthURL)
	th.CheckEquals(t, (*string)(nil), serviceProviders[0].Description)
	th.CheckEquals(t, true, serviceProviders[0].Enabled)
	th.CheckEquals(t, "ACME", serviceProviders[0].ID)
	th.CheckDeepEquals(t, map[string]any{
		"self": "https://example.com/v3/OS-FEDERATION/service_providers/ACME",
	}, serviceProviders[0].Links)
	th.CheckEquals(t, (*string)(nil), serviceProviders[0].RelayStatePrefix)
	th.CheckEquals(t, "https://example.com/sp", serviceProviders[0].SPURL)

	createOpts := federation.CreateServiceProviderOpts{
		AuthURL:          "https://example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth",
		Description:      "Remote Service Provider",
		Enabled:          gophercloud.Enabled,
		RelayStatePrefix: "ss:mem:",
		SPURL:            "https://example.com/Shibboleth.sso/SAML2/ECP",
	}
	created, err := federation.CreateServiceProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME", createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "ACME", created.ID)
	th.CheckEquals(t, createOpts.AuthURL, created.AuthURL)
	th.AssertTrue(t, created.Description != nil)
	th.CheckEquals(t, createOpts.Description, *created.Description)
	th.CheckEquals(t, true, created.Enabled)
	th.CheckDeepEquals(t, map[string]any{
		"self": "https://example.com/v3/OS-FEDERATION/service_providers/ACME",
	}, created.Links)
	th.AssertTrue(t, created.RelayStatePrefix != nil)
	th.CheckEquals(t, createOpts.RelayStatePrefix, *created.RelayStatePrefix)
	th.CheckEquals(t, createOpts.SPURL, created.SPURL)

	actual, err := federation.GetServiceProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, created, actual)

	authURL := "https://new.example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth"
	description := "Updated Service Provider"
	relayStatePrefix := "ss:temp:"
	spURL := "https://new.example.com/Shibboleth.sso/SAML2/ECP"
	updateOpts := federation.UpdateServiceProviderOpts{
		AuthURL:          &authURL,
		Description:      &description,
		Enabled:          gophercloud.Disabled,
		RelayStatePrefix: &relayStatePrefix,
		SPURL:            &spURL,
	}
	updated, err := federation.UpdateServiceProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, authURL, updated.AuthURL)
	th.AssertTrue(t, updated.Description != nil)
	th.CheckEquals(t, description, *updated.Description)
	th.CheckEquals(t, false, updated.Enabled)
	th.AssertTrue(t, updated.RelayStatePrefix != nil)
	th.CheckEquals(t, relayStatePrefix, *updated.RelayStatePrefix)
	th.CheckEquals(t, spURL, updated.SPURL)

	err = federation.DeleteServiceProvider(context.TODO(), client.ServiceClient(fakeServer), "ACME").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestServiceProviderUpdateOptsNullableFields(t *testing.T) {
	actual, err := (federation.UpdateServiceProviderOpts{}).ToServiceProviderUpdateMap()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `{"service_provider": {}}`, actual)

	description := ""
	relayStatePrefix := ""
	opts := federation.UpdateServiceProviderOpts{
		Description:      &description,
		RelayStatePrefix: &relayStatePrefix,
	}
	actual, err = opts.ToServiceProviderUpdateMap()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `{
		"service_provider": {
			"description": null,
			"relay_state_prefix": null
		}
	}`, actual)
}

func TestServiceProviderCreateOptsValidation(t *testing.T) {
	_, err := (federation.CreateServiceProviderOpts{}).ToServiceProviderCreateMap()
	th.AssertErr(t, err)
}
