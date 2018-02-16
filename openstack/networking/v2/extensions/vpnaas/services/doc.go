/*
Package services allows management and retrieval of VPN services in the
OpenStack Networking Service.


Example to Create a Service

	createOpts := services.CreateOpts{
		Name:        "vpnservice1",
		Description: "A service"
	}

	service, err := services.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/
package services
