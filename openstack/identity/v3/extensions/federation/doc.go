/*
Package federation provides information and interaction with OS-FEDERATION API for the
Openstack Identity service.

Example to List Mappings

	allPages, err := federation.ListMappings(identityClient).AllPages()
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

	createdMapping, err := federation.CreateMapping(identityClient, "ACME", createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a Mapping

	mapping, err := federation.GetMapping(identityClient, "ACME").Extract()
	if err != nil {
		panic(err)
	}
*/
package federation
