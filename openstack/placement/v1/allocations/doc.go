/*
Package allocations manages consumer allocations from the OpenStack Placement service.

Allocation API requests are available starting from microversion 1.0.

# Example to get allocations for a consumer

	allocs, err := allocations.Get(context.TODO(), placementClient, consumerUUID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", allocs)
*/
package allocations
