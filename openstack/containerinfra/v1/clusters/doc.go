/*
Package clusters contains functionality for working with Magnum Cluster resources.

Example to Create a Cluster

	masterCount := 1
	nodeCount := 1
	createTimeout := 30
	masterLBEnabled := true
	createOpts := clusters.CreateOpts{
		ClusterTemplateID: "0562d357-8641-4759-8fed-8173f02c9633",
		CreateTimeout:     &createTimeout,
		DiscoveryURL:      "",
		FlavorID:          "m1.small",
		KeyPair:           "my_keypair",
		Labels:            map[string]string{},
		MasterCount:       &masterCount,
		MasterFlavorID:    "m1.small",
		Name:              "k8s",
		NodeCount:         &nodeCount,
		MasterLBEnabled:   &masterLBEnabled,
	}

	cluster, err := clusters.Create(context.TODO(), serviceClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a Cluster

	clusterName := "cluster123"
	cluster, err := clusters.Get(context.TODO(), serviceClient, clusterName).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cluster)

Example to List Clusters

	listOpts := clusters.ListOpts{
		Limit: 20,
	}

	allPages, err := clusters.List(serviceClient, listOpts).AllPages(context.TODO())
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

Example to List Clusters with detailed information

	allPagesDetail, err := clusters.ListDetail(serviceClient, clusters.ListOpts{}).AllPages(context.TODO())
	if err != nil {
	    panic(err)
	}

	allClustersDetail, err := clusters.ExtractClusters(allPagesDetail)
	if err != nil {
	    panic(err)
	}

	for _, clusterDetail := range allClustersDetail {
	    fmt.Printf("%+v\n", clusterDetail)
	}

Example to Update a Cluster

	updateOpts := []clusters.UpdateOptsBuilder{
		clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  "/master_lb_enabled",
			Value: "True",
		},
		clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  "/registry_enabled",
			Value: "True",
		},
	}
	clusterUUID, err := clusters.Update(context.TODO(), serviceClient, clusterUUID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", clusterUUID)

Example to Upgrade a Cluster

	upgradeOpts := clusters.UpgradeOpts{
		ClusterTemplate: "0562d357-8641-4759-8fed-8173f02c9633",
	}
	clusterUUID, err := clusters.Upgrade(context.TODO(), serviceClient, clusterUUID, upgradeOpts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", clusterUUID)

Example to Delete a Cluster

	clusterUUID := "dc6d336e3fc4c0a951b5698cd1236ee"
	err := clusters.Delete(context.TODO(), serviceClient, clusterUUID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package clusters
