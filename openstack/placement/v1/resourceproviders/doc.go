/*
Package resourceproviders creates and lists all resource providers from the OpenStack Placement service.

Example to list resource providers

	allPages, err := resourceproviders.List(placementClient, resourceproviders.ListOpts{}).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allResourceProviders, err := resourceproviders.ExtractResourceProviders(allPages)
	if err != nil {
		panic(err)
	}

	for _, r := range allResourceProviders {
		fmt.Printf("%+v\n", r)
	}

Example to create resource providers

	createOpts := resourceproviders.CreateOpts{
		Name: "new-rp",
		UUID: "b99b3ab4-3aa6-4fba-b827-69b88b9c544a",
		ParentProvider: "c7f50b40-6f32-4d7a-9f32-9384057be83b"
	}

	rp, err := resourceproviders.Create(context.TODO(), placementClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a resource provider

	resourceProviderID := "b99b3ab4-3aa6-4fba-b827-69b88b9c544a"
	err := resourceproviders.Delete(context.TODO(), placementClient, resourceProviderID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Get a resource provider

	resourceProviderID := "b99b3ab4-3aa6-4fba-b827-69b88b9c544a"
	resourceProvider, err := resourceproviders.Get(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a resource provider

	resourceProviderID := "b99b3ab4-3aa6-4fba-b827-69b88b9c544a"

	updateOpts := resourceproviders.UpdateOpts{
		Name: "new-rp",
		ParentProvider: "c7f50b40-6f32-4d7a-9f32-9384057be83b"
	}

	placementClient.Microversion = "1.37"
	resourceProvider, err := resourceproviders.Update(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers usages

	rp, err := resourceproviders.GetUsages(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers inventories

	rp, err := resourceproviders.GetInventories(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get one resource provider inventory

	rpInventory, err := resourceproviders.GetInventory(context.TODO(), placementClient, resourceProviderID, "VCPU").Extract()
	if err != nil {
		panic(err)
	}

Example to update (replace) all resource provider inventories

	inventories, err := resourceproviders.GetInventories(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

	updateInventoriesOpts := resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.Inventory{
			"VCPU": {
				Total:           4,
				Reserved:        0,
				MinUnit:         1,
				MaxUnit:         4,
				StepSize:        1,
				AllocationRatio: 16.0,
			},
		},
	}

	rp, err = resourceproviders.UpdateInventories(context.TODO(), placementClient, resourceProviderID, updateInventoriesOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to update one existing resource provider inventory

	inventories, err := resourceproviders.GetInventories(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

	// UpdateInventory updates an existing resource class inventory.
	updateInventoryOpts := resourceproviders.UpdateInventoryOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventory: resourceproviders.Inventory{
			Total:           4,
			Reserved:        0,
			MinUnit:         1,
			MaxUnit:         4,
			StepSize:        1,
			AllocationRatio: 16.0,
		},
	}

	rpInventory, err := resourceproviders.UpdateInventory(context.TODO(), placementClient, resourceProviderID, "VCPU", updateInventoryOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to delete one existing resource provider inventory
Since this request does not accept the resource provider generation, it is not safe to use when multiple threads are managing inventories for a single provider. In such situations use UpdateInventories with the empty inventory.

	err = resourceproviders.DeleteInventory(context.TODO(), placementClient, resourceProviderID, "VCPU").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to delete all resource provider inventories
Since this request does not accept the resource provider generation, it is not safe to use when multiple threads are managing inventories for a single provider. In such situations use UpdateInventories with an empty inventory map.

	err = resourceproviders.DeleteInventories(context.TODO(), placementClient, resourceProviderID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to get resource providers traits

	rp, err := resourceproviders.GetTraits(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers allocations

	rp, err := resourceproviders.GetAllocations(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers aggregates

	placementClient.Microversion = "1.1"

	rp, err := resourceproviders.GetAggregates(context.TODO(), placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

# Example to update resource providers aggregates

For microversion 1.18 and earlier the ResourceProviderGeneration is optional and would be ignored if provided,
as it was not supported pre-1.19. For greater safety, it is recommended to use the newer microversion.

	placementClient.Microversion = "1.19"

	updateOpts := resourceproviders.UpdateAggregatesOpts{
		ResourceProviderGeneration: rp.ResourceProviderGeneration,
		Aggregates: []string{
			"6d84f6f6-7736-40ff-84d2-7db47f18ea25",
			"f11f14bc-6f17-4f0a-b7c2-44b3e685ccf4",
		},
	}

	rp, err = resourceproviders.UpdateAggregates(context.TODO(), placementClient, resourceProviderID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package resourceproviders
