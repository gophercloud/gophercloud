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

Example to Delete a Service

	serviceID := "38aee955-6283-4279-b091-8b9c828000ec"
	err := policies.Delete(networkClient, serviceID).ExtractErr()
	if err != nil {
		panic(err)
	}

*/
package services
