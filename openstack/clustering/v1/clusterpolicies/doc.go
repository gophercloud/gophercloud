/*
Package clusterpolicies enables Lists all cluster policies and shows information for a cluster policy from the OpenStack
Clustering Service.

Example to list cluster policies for a Senlin deployment

  allPages, err := clusterpolicies.ListDetail(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allClusterPolicies, err := actions.ExtractClusterPolicies(allPages)
	if err != nil {
		panic(err)
	}

	for _, clusterPolicy := range allClusterPolicies {
		fmt.Printf("%+v\n", clusterPolicy)
	}
*/
package clusterpolicies
