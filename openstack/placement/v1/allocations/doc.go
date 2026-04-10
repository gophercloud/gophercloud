/*
Package allocations manages consumer allocations from the OpenStack Placement service.

Allocation API requests are available starting from microversion 1.0.

# Example to get allocations for a consumer

	allocs, err := allocations.Get(context.TODO(), placementClient, consumerUUID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", allocs)

# Example to set allocations for a new consumer (microversion 1.28+)

When creating allocations for a consumer that does not yet exist, set
ConsumerGeneration to nil. It will be serialized as JSON null, telling the
server that no prior allocation is expected. This is different from omitting
the field: nil is required, not optional.

	placementClient.Microversion = "1.28"

	err := allocations.Update(context.TODO(), placementClient, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			providerUUID: {
				Resources: map[string]int{"VCPU": 2, "MEMORY_MB": 2048},
			},
		},
		ProjectID:          projectID,
		UserID:             userID,
		ConsumerGeneration: nil,
	}).ExtractErr()
	if err != nil {
		panic(err)
	}

# Example to update allocations for an existing consumer (microversion 1.28+)

When updating allocations for an existing consumer, retrieve the current
generation first and pass it in.

	placementClient.Microversion = "1.28"

	existing, err := allocations.Get(context.TODO(), placementClient, consumerUUID).Extract()
	if err != nil {
		panic(err)
	}

	err = allocations.Update(context.TODO(), placementClient, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			providerUUID: {
				Resources: map[string]int{"VCPU": 4, "MEMORY_MB": 4096},
			},
		},
		ProjectID:          *existing.ProjectID,
		UserID:             *existing.UserID,
		ConsumerGeneration: existing.ConsumerGeneration,
	}).ExtractErr()
	if err != nil {
		panic(err)
	}

# Example to delete all allocations for a consumer

Note: using Update with an empty Allocations map is generally safer because
it is protected by the consumer generation check. Use Delete only when you
do not need that protection.

	err = allocations.Delete(context.TODO(), placementClient, consumerUUID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package allocations
