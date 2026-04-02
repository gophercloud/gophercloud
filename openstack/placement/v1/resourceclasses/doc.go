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

Example to Create a resource class using POST

	placementClient.Microversion = "1.2"

	createOpts := resourceclasses.CreateOpts{
		Name: "CUSTOM_RESOURCE_CLASS",
	}

	err := resourceclasses.Create(context.TODO(), placementClient, createOpts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to ensure the existence of a resource class using PUT (idempotent creation)

	placementClient.Microversion = "1.7"

	err := resourceclasses.Update(context.TODO(), placementClient, "CUSTOM_RESOURCE_CLASS").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package resourceclasses
