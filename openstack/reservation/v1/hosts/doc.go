/*
Package hosts manages compute hosts in the OpenStack Reservation service.

A host must be enrolled into Blazar's freepool before any reservation can draw
on it. Listing the freepool is administrative.

Example to list hosts

	listOpts := hosts.ListOpts{}

	allPages, err := hosts.List(reservationClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allHosts, err := hosts.ExtractHosts(allPages)
	if err != nil {
		panic(err)
	}

	for _, h := range allHosts {
		fmt.Printf("%+v\n", h)
	}

Example to create a host

	createOpts := hosts.CreateOpts{
		Name: "compute-1.example.com",
		ExtraCapabilities: map[string]any{
			"gpu": "a100",
		},
	}

	host, err := hosts.Create(context.TODO(), reservationClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to get a host

	host, err := hosts.Get(context.TODO(), reservationClient, "18").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", host)

Example to update the extra capabilities of a host

	updateOpts := hosts.UpdateOpts{
		ExtraCapabilities: map[string]any{
			"gpu":  "h100",
			"rack": nil,
		},
	}

	host, err := hosts.Update(context.TODO(), reservationClient, "18", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to delete a host

	err := hosts.Delete(context.TODO(), reservationClient, "18").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to list the allocations of a lease

	listOpts := hosts.ListAllocationsOpts{
		LeaseID: "98c3544d-0afe-4251-8556-700c847a127f",
	}

	allPages, err := hosts.ListAllocations(reservationClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allAllocations, err := hosts.ExtractAllocations(allPages)
	if err != nil {
		panic(err)
	}

	for _, allocation := range allAllocations {
		fmt.Printf("%+v\n", allocation)
	}

Example to get the allocations of a host

	allocation, err := hosts.GetAllocation(context.TODO(), reservationClient, "18", nil).Extract()
	if err != nil {
		panic(err)
	}

Example to list the resource properties with their values

	listOpts := hosts.ListResourcePropertiesOpts{
		Detail: true,
	}

	allPages, err := hosts.ListResourceProperties(reservationClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allProperties, err := hosts.ExtractResourceProperties(allPages)
	if err != nil {
		panic(err)
	}

Example to hide a resource property from users

	updateOpts := hosts.UpdateResourcePropertyOpts{
		Private: true,
	}

	property, err := hosts.UpdateResourceProperty(context.TODO(), reservationClient, "gpu", updateOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package hosts
