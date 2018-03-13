/*
Package clusters provides information and interaction with the clusters through
the OpenStack Compute service.

Lists all clusters and creates, shows information for, updates, deletes, and triggers an action on a cluster.

Example to List Clusters

	listOpts := clusters.ListOpts{
		Name: "testcluster",
	}

	allPages, err := clusters.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allClusters, err := images.ExtractClusters(allPages)
	if err != nil {
		panic(err)
	}

	for _, cluster := range allClusters {
		fmt.Printf("%+v\n", cluster)
	}

Example to Create a cluster

	createOpts := clusters.CreateOpts{
    Name:       "testcluster"
    DesiredCapacity: 1
    ProfileID:  "b7b870ee-d3c5-4a93-b9d7-846c53b2c2da"
	}

	cluster, err := clusters.Create(computeClient, createOpts).ExtractCluster()
	if err != nil {
		panic(err)
	}

Example to ScaleIn a cluster

	scaleOpts := clusters.ScaleOpts{
    Count: 2,
	}
  clusterID:  "b7b870e3-d3c5-4a93-b9d7-846c53b2c2da"

	action, err := clusters.ScaleIn(computeClient, clusterID, scaleOpts).ExtractAction()
	if err != nil {
		panic(err)
	}
*/
package clusters
