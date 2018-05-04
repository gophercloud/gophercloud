/*
Package policies provides information and interaction with the policies API
resource for the OpenStack Identity service.

Example to List Policies

	listOpts := policies.ListOpts{
		Type: "application/json",
	}

	allPages, err := policies.List(identityClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allPolicies, err := policies.ExtractPolicies(allPages)
	if err != nil {
		panic(err)
	}

	for _, policy := range allPolicies {
		fmt.Printf("%+v\n", policy)
	}

Example to Create a Policy

	createOpts := policies.CreateOpts{
		Type: "application/json",
		Blob: "{'foobar_user': 'role:compute-user'}",
		Extra: map[string]interface{}{
			"description": "policy for foobar_user",
		},
	}

	policy, err := policies.Create(identityClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package policies
