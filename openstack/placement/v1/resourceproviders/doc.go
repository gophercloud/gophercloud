/*
Package resourceproviders creates and lists all resource providers from the OpenStack Placement service.

Example to list resource providers

	allPages, err := resourceproviders.List(placementClient, resourceproviders.ListOpts{}).AllPages()
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
	}

	rp, err := resourceproviders.Create(placementClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers usages

	rp, err := resourceproviders.GetUsages(placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers inventories

	rp, err := resourceproviders.GetInventories(placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers traits

	rp, err := resourceproviders.GetTraits(placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

Example to get resource providers allocations

	rp, err := resourceproviders.GetAllocations(placementClient, resourceProviderID).Extract()
	if err != nil {
		panic(err)
	}

*/
package resourceproviders
