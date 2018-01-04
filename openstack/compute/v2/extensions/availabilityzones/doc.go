/*
Package availabilityzones provides the ability to get lists and detailed
availability zone information.

Example of Get Availability Zone Information

	allPages, err := az.List(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	availabilityZoneInfo, err := az.ExtractAvailabilityZones(allPages)
	if err != nil {
		panic(err)
	}

	for _, zoneInfo := range availabilityZoneInfo {
		fmt.Printf("Zone name: %s\nAvailable: %v\n", zoneInfo.ZoneName,
					zoneInfo.ZoneState.Available)
	}

Example of Get Detailed Availability Zone Information

	allPages, err := az.ListDetail(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	availabilityZoneInfo, err := az.ExtractAvailabilityZones(allPages)
	if err != nil {
		panic(err)
	}

	for _, zoneInfo := range availabilityZoneInfo {
		fmt.Printf("Zone name: %s\nAvailable: %v\n", zoneInfo.ZoneName,
					zoneInfo.ZoneState.Available)
		for hostname, services := range zoneInfo.Hosts {
			fmt.Println(hostname)
			// to be continued ...
			// cut for brevity :)
		}
	}
*/
package availabilityzones
