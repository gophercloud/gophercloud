/*
Package clusters contains functionality for working with Magnum Cluster resources.

Example to Create a Cluster

	masterCount := 1
	nodeCount := 1
	createTimeout := 30
	opts := clusters.CreateOpts{
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
	}

	cluster, err := clusters.Create(serviceClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/
package clusters
