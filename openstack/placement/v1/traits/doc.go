/*
Package traits lists traits from the OpenStack Placement service.

Traits API requests are available starting from microversion 1.6.

Example to list traits

	placementClient.Microversion = "1.6"

	allPages, err := traits.List(placementClient, traits.ListOpts{}).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allTraits, err := traits.ExtractTraits(allPages)
	if err != nil {
		panic(err)
	}

	for _, t := range allTraits {
		fmt.Println(t)
	}

Example to check if a trait exists

	placementClient.Microversion = "1.6"

	traitName := "CUSTOM_HW_FPGA_CLASS1"
	err := traits.Get(context.TODO(), placementClient, traitName).ExtractErr()
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			// 404 Not Found - The trait does not exist
			fmt.Println("Trait does not exist.")
		} else {
			// Another error occurred
			panic(err)
		}
	} else {
		fmt.Println("Trait exists!")
	}

Example to create a trait

	placementClient.Microversion = "1.6"

	traitName := "CUSTOM_HW_FPGA_CLASS1"
	createOpts := traits.CreateOpts{}
	err := traits.Create(context.TODO(), placementClient, traitName, createOpts).ExtractErr()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Trait created successfully!")
	}
*/
package traits
