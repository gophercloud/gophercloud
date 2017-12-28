/*
Package subnetpools provides the ability to retrieve and manage subnetpools through the Neutron API.

Example of Listing Subnetpools.

	listOpts := subnets.ListOpts{
		IPVersion: 6,
	}

	allPages, err := subnetpools.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allSubnetpools, err := subnetpools.ExtractSubnetPools(allPages)
	if err != nil {
		panic(err)
	}

	for _, subnetpools := range allSubnetpools {
		fmt.Printf("%+v\n", subnetpools)
	}
*/
package subnetpools
