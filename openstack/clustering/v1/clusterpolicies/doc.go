/*
Package clusterpolicies enables Lists all cluster policies and shows information for a cluster policy from the OpenStack
Clustering Service.

Example to list cluster policies for a Senlin deployment

	clusterID := "7d85f602-a948-4a30-afd4-e84f47471c15"
		allPages, err := clusterpolicies.List(serviceClient, clusterID, clusterpolicies.ListOpts{}).AllPages()
	if err != nil {
		panic(err)
	}

	allClusterPolicies, err := clusterpolicies.ExtractClusterPolicies(allPages)
	if err != nil {
		panic(err)
	}

	for _, clusterPolicy := range allClusterPolicies {
		fmt.Printf("%+v\n", clusterPolicy)
	}

*/
package clusterpolicies
