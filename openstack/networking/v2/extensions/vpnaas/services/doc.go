/*
Package services allows management and retrieval of VPN services in the
OpenStack Networking Service.

Example to List Services

	listOpts := services.ListOpts{
		TenantID: "966b3c7d36a24facaf20b7e458bf2192",
	}

	allPages, err := services.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allPolicies, err := services.ExtractServices(allPages)
	if err != nil {
		panic(err)
	}

	for _, service := range allServices {
		fmt.Printf("%+v\n", service)
	}
*/
package services

