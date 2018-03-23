/*
Package clusters provides information and interaction with the clusters through
the OpenStack Clustering service.

Example to Create a cluster

	createOpts := clusters.CreateOpts{
		Name:            "test-cluster",
		DesiredCapacity: 1,
		ProfileUUID:     "b7b870ee-d3c5-4a93-b9d7-846c53b2c2da",
	}

	cluster, err := clusters.Create(serviceClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get Clusters

	clusterName := "cluster123"
	cluster, err := clusters.Get(serviceClient, clusterName).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cluster)

Example to List Clusters

	listOpts := clusters.ListOpts{
		Name: "testcluster",
	}

	allPages, err := clusters.ListDetail(serviceClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allClusters, err := clusters.ExtractClusters(allPages)
	if err != nil {
		panic(err)
	}

	for _, cluster := range allClusters {
		fmt.Printf("%+v\n", cluster)
	}

Example to Update a cluster

	updateOpts := clusters.UpdateOpts{
		Name:       "testcluster",
		ProfileID:  "b7b870ee-d3c5-4a93-b9d7-846c53b2c2da",
	}

	clusterID := "7d85f602-a948-4a30-afd4-e84f47471c15"
	cluster, err := clusters.Update(serviceClient, clusterName, opts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cluster)

Example to Delete a cluster

	clusterID := "dc6d336e3fc4c0a951b5698cd1236ee"
	err := clusters.Delete(serviceClient, clusterID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Resize a cluster

	number := 1
	maxSize := 5
	minSize := 1
	minStep := 1
	strict := true

	resizeOpts := clusters.ResizeOpts{
		AdjustmentType: clusters.ChangeInCapacityAdjustment,
		Number:         number,
		MaxSize:        &maxSize,
		MinSize:        &minSize,
		MinStep:        &minStep,
		Strict:         &strict,
	}

	actionID, err := clusters.Resize(client, clusterName, resizeOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to resize cluster: %v", err)
	}
	fmt.Println("Resize actionID", actionID)

Example to ScaleIn a cluster

	count := 2
	scaleInOpts := clusters.ScaleInOpts{
		Count: &count,
	}
	clusterID:  "b7b870e3-d3c5-4a93-b9d7-846c53b2c2da"

	action, err := clusters.ScaleIn(computeClient, clusterID, scaleInOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package clusters
