/*
Package policytypes lists all policy types and shows details for a policy type from the OpenStack
Clustering Service.

Example to list policy types for a Senlin deployment

  allPages, err := policytypes.ListDetail(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allPolicyTypes, err := actions.ExtractPolicyTypes(allPages)
	if err != nil {
		panic(err)
	}

	for _, policyType := range allPolicyTypes {
		fmt.Printf("%+v\n", policyType)
	}
*/
package policytypes
