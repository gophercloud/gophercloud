/*
Package services returns information about the compute services in the OpenStack
cloud.

Example of Retrieving list of all services

	opts := services.ListOpts{
		Binary: "nova-scheduler",
	}

	allPages, err := services.List(computeClient, opts).AllPages()
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

Example of delete a service

	updated, err := services.Delete(client, serviceID).Extract()
	if err != nil {
		panic(err)
	}
*/

package services
