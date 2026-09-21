/*
Package federation provides information and interaction with the OS-FEDERATION
API for the OpenStack Identity service.

Example to Create an Identity Provider

	createOpts := federation.CreateIdentityProviderOpts{
		Description: "Stores ACME identities",
		Enabled:     gophercloud.Enabled,
		RemoteIDs:   []string{"acme_id_1", "acme_id_2"},
	}
	identityProvider, err := federation.CreateIdentityProvider(
		context.TODO(), identityClient, "ACME", createOpts,
	).Extract()
	if err != nil {
		panic(err)
	}

Example to List Identity Providers

	listOpts := federation.ListIdentityProvidersOpts{
		Enabled: gophercloud.Enabled,
	}
	allPages, err := federation.ListIdentityProviders(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}
	allIdentityProviders, err := federation.ExtractIdentityProviders(allPages)
	if err != nil {
		panic(err)
	}

Example to List Mappings

	allPages, err := federation.ListMappings(identityClient).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}
	allMappings, err := federation.ExtractMappings(allPages)
	if err != nil {
		panic(err)
	}

Example to Create Mappings

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

	createdMapping, err := federation.CreateMapping(context.TODO(), identityClient, "ACME", createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a Mapping

	mapping, err := federation.GetMapping(context.TODO(), identityClient, "ACME").Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Mapping

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
	updatedMapping, err := federation.UpdateMapping(context.TODO(), identityClient, "ACME", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Mapping

	err := federation.DeleteMapping(context.TODO(), identityClient, "ACME").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Add a Protocol to an Identity Provider

	createOpts := federation.CreateProtocolOpts{
		MappingID:         "ACME",
		RemoteIDAttribute: "HTTP_OIDC_ISS",
	}
	protocol, err := federation.CreateProtocol(
		context.TODO(), identityClient, "ACME", "openid", createOpts,
	).Extract()
	if err != nil {
		panic(err)
	}

Example to Create a Service Provider

	createOpts := federation.CreateServiceProviderOpts{
		AuthURL:          "https://example.com/v3/OS-FEDERATION/identity_providers/acme/protocols/saml2/auth",
		Description:      "Remote Service Provider",
		Enabled:          gophercloud.Enabled,
		RelayStatePrefix: "ss:mem:",
		SPURL:            "https://example.com/Shibboleth.sso/SAML2/ECP",
	}
	serviceProvider, err := federation.CreateServiceProvider(
		context.TODO(), identityClient, "ACME", createOpts,
	).Extract()
	if err != nil {
		panic(err)
	}
*/
package federation
