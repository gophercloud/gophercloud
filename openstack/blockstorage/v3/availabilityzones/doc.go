/*
Package availabilityzones provides information and interaction with the
availability zone API in the Openstack Block Storage service.

Example to list availability zones

	allPages, err := availabilityzones.List(client).AllPages()
	if err != nil {
		panic(err)
	}

	allAavailabilityZones, err := availabilityzones.ExtractAvailabilityZone(allPages)
	if err != nil {
		panic(err)
	}

	for _, availabilityZone := range allAvailabilityZones{
		fmt.Printf("List: %+v\n", availabilityZone)
	}
*/
package availabilityzones
