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
*/
package resourceproviders
