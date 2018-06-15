// Package clustertemplates contains functionality for working with Magnum Cluster Templates
// resources.
/*
Package clustertemplates provides information and interaction with the cluster-templates through
the OpenStack Container Infra service.

Example to Create Cluster Template

	createOpts := clustertemplates.CreateOpts{
		Name:                "test-cluster-template",
		Labels:              map[string]string{},
		FixedSubnet:         nil,
		MasterFlavorID:      nil,
		NoProxy:             "10.0.0.0/8,172.0.0.0/8,192.0.0.0/8,localhost",
		HTTPSProxy:          "http://10.164.177.169:8080",
		TLSDisabled:         false,
		KeyPairID:           "kp",
		Public:              false,
		HTTPProxy:           "http://10.164.177.169:8080",
		DockerVolumeSize:    3,
		ServerType:          "vm",
		ExternalNetworkID:   "public",
		ImageID:             "fedora-atomic-latest",
		VolumeDriver:        "cinder",
		RegistryEnabled:     false,
		DockerStorageDriver: "devicemapper",
		NetworkDriver:       "flannel",
		FixedNetwork:        nil,
		COE:                 "kubernetes",
		FlavorID:            "m1.small",
		MasterLBEnabled:     true,
		DNSNameServer:       "8.8.8.8",
	}

	clustertemplate, err := clustertemplates.Create(serviceClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to List Clusters Templates

	listOpts := clustertemplates.ListOpts{
		Limit: 20,
	}

	allPages, err := clustertemplates.List(serviceClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allClusterTemplates, err := clusters.ExtractClusterTemplates(allPages)
	if err != nil {
		panic(err)
	}

	for _, clusterTemplate := range allClusterTemplates {
		fmt.Printf("%+v\n", clusterTemplate)
	}

Example to Update Cluster Template

	updateOpts := []clustertemplates.UpdateOptsBuilder{
		clustertemplates.UpdateOpts{
			Path:  "/master_lb_enabled",
			Value: true,
			Op:    "replace",
		},
		clustertemplates.UpdateOpts{
			Path:  "/registry_enabled",
			Value: true,
			Op:    "replace",
		},
	}

	clustertemplate, err := clustertemplates.Update(serviceClient, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete Cluster Template

	clusterTemplateID := "dc6d336e3fc4c0a951b5698cd1236ee"
	err := clustertemplates.Delete(serviceClient, clusterTemplateID).ExtractErr()
	if err != nil {
		panic(err)
	}

*/
package clustertemplates
