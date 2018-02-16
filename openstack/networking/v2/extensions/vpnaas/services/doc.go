/*
Package services allows management and retrieval of VPN services in the
OpenStack Networking Service.


Example to Create a Service

	createOpts := services.CreateOpts{
		Name:        "vpnservice1",
		Description: "A service",
		RouterID:	 "2512e759-e8d7-4eea-a0af-4a85927a2e59",
		AdminStateUp: gophercloud.Enabled,
	}

	service, err := services.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/
package services
