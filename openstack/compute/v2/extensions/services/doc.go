/*
Package services returns information about the compute services in the OpenStack
cloud.

Example of Retrieving list of all services

	allPages, err := services.List(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allServices, err := services.ExtractServices(allPages)
	if err != nil {
		panic(err)
	}

	for _, service := range allServices {
		fmt.Printf("%+v\n", service)
	}

Example of updating a service

	opts := services.UpdateOpts{
		Status: services.ServiceDisabled,
	}

	updated, err := services.Update(client, serviceID, opts).Extract()
	if err != nil {
		panic(err)
	}
*/

package services
