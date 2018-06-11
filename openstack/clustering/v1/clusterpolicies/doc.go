/*
Package clusterpolicies Gets and Lists all cluster policies and shows detailed information for a cluster policy from the OpenStack
Clustering Service.

Example to get cluster-policies

	clusterID := "7d85f602-a948-4a30-afd4-e84f47471c15"
	profileID := "714fe676-a08f-4196-b7af-61d52eeded15"
	clusterPolicy, err := clusterpolicies.Get(serviceCLient, clusterID, profileID).Extract()
	fmt.Println("ClusterPolicy=", clusterPolicy)

*/
package clusterpolicies
