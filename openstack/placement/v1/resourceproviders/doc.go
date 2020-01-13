/*
Package resourceproviders lists all resource providers from the OpenStack Placement service.

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

*/
package resourceproviders
