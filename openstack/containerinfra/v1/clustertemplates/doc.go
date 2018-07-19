// Package clustertemplates contains functionality for working with Magnum Cluster Templates
// resources.
/*
Package clustertemplates provides information and interaction with the cluster-templates through
the OpenStack Container Infra service.

Example to Create Cluster Template

	boolFalse := false
	boolTrue := true
	createOpts := clustertemplates.CreateOpts{
		Name:                "test-cluster-template",
		Labels:              map[string]string{},
		FixedSubnet:         "",
		MasterFlavorID:      "",
		NoProxy:             "10.0.0.0/8,172.0.0.0/8,192.0.0.0/8,localhost",
		HTTPSProxy:          "http://10.164.177.169:8080",
		TLSDisabled:         &boolFalse,
		KeyPairID:           "kp",
		Public:              &boolFalse,
		HTTPProxy:           "http://10.164.177.169:8080",
		DockerVolumeSize:    3,
		ServerType:          "vm",
		ExternalNetworkID:   "public",
		ImageID:             "fedora-atomic-latest",
		VolumeDriver:        "cinder",
		RegistryEnabled:     &boolFalse,
		DockerStorageDriver: "devicemapper",
		NetworkDriver:       "flannel",
		FixedNetwork:        "",
		COE:                 "kubernetes",
		FlavorID:            "m1.small",
		MasterLBEnabled:     &boolTrue,
		DNSNameServer:       "8.8.8.8",
	}

	clustertemplate, err := clustertemplates.Create(serviceClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/
package clustertemplates
