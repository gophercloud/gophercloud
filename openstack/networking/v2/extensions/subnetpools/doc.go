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

Example to Get a Subnetpool

	subnetPoolID = "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55"
	subnetPool, err := subnetpools.Get(networkClient, subnetPoolID).Extract()
	if err != nil {
		panic(err)
	}

Example to Create a new Subnetpool.

	subnetPoolName := "private_pool"
	subnetPoolPrefixes := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	subnetPoolOpts := subnetpools.CreateOpts{
		Name: subnetPoolName,
		Prefixes: subnetPoolPrefixes,
	}
	subnetPool, err := subnetpools.Create(networkClient, subnetPoolOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package subnetpools
