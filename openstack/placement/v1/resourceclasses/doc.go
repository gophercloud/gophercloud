/*
Package resourceclasses manages resource classes from the OpenStack Placement service.

Resource Class API requests are available starting from microversion 1.2.

Example to list resource classes

	placementClient.Microversion = "1.2"

	allPages, err := resourceclasses.List(placementClient).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allResourceClasses, err := resourceclasses.ExtractResourceClasses(allPages)
	if err != nil {
		panic(err)
	}

	for _, rc := range allResourceClasses {
		fmt.Printf("%+v\n", rc)
	}

Example to Get a resource class

	placementClient.Microversion = "1.2"

	resourceClassName := "VCPU"
	resourceClass, err := resourceclasses.Get(context.TODO(), placementClient, resourceClassName).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resourceClass)
*/
package resourceclasses
