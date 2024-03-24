/*
Package flavorprofiles provides information and interaction
with FlavorProfiles for the OpenStack Load-balancing service.

Example to List FlavorProfiles

	listOpts := flavorprofiles.ListOpts{}

	allPages, err := flavorprofiles.List(octaviaClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allFlavorProfiles, err := flavorprofiles.ExtractFlavorProfiles(allPages)
	if err != nil {
		panic(err)
	}

	for _, flavorProfile := range allFlavorProfiles {
		fmt.Printf("%+v\n", flavorProfile)
	}

Example to Create a FlavorProfile

	createOpts := flavorprofiles.CreateOpts{
		Name:         "amphora-single",
		ProviderName: "amphora",
		FlavorData:   "{\"loadbalancer_topology\": \"SINGLE\"}",
	}

	flavorProfile, err := flavorprofiles.Create(context.TODO(), octaviaClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a FlavorProfile

	flavorProfileID := "dd6a26af-8085-4047-a62b-3080f4c76521"

	updateOpts := flavorprofiles.UpdateOpts{
		Name: "amphora-single-updated",
	}

	flavorProfile, err := flavorprofiles.Update(context.TODO(), octaviaClient, flavorProfileID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a FlavorProfile

	flavorProfileID := "dd6a26af-8085-4047-a62b-3080f4c76521"
	err := flavorprofiles.Delete(context.TODO(), octaviaClient, flavorProfileID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package flavorprofiles
