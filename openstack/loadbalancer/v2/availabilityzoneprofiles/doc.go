/*
Package availabilityzoneprofiles provides information and interaction
with AvailabilityZoneProfiles for the OpenStack Load-balancing service.

Example to List AvailabilityZoneProfiles

	listOpts := availabilityzoneprofiles.ListOpts{}

	allPages, err := availabilityzoneprofiles.List(octaviaClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allAvailabilityZoneProfiles, err := availabilityzoneprofiles.ExtractAvailabilityZoneProfiles(allPages)
	if err != nil {
		panic(err)
	}

	for _, availabilityZoneProfile := range allAvailabilityZoneProfiles {
		fmt.Printf("%+v\n", availabilityZoneProfile)
	}

Example to Create a AvailabilityZoneProfile

	createOpts := availabilityzoneprofiles.CreateOpts{
		Name:                 "availability-zone-profile",
		ProviderName:         "amphora",
		AvailabilityZoneData: "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}",
	}

	availabilityZoneProfile, err := availabilityzoneprofiles.Create(context.TODO(), octaviaClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a AvailabilityZoneProfile

	availabilityZoneProfileID := "0c359d38-6164-498f-8409-5b11d05b6226"

	updateOpts := availabilityzoneprofiles.UpdateOpts{
		Name: "availability-zone-profile-updated",
	}

	availabilityZoneProfile, err := availabilityzoneprofiles.Update(context.TODO(), octaviaClient, availabilityZoneProfileID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a AvailabilityZoneProfile

	availabilityZoneProfileID := "0c359d38-6164-498f-8409-5b11d05b6226"
	err := availabilityzoneprofiles.Delete(context.TODO(), octaviaClient, availabilityZoneProfileID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package availabilityzoneprofiles
