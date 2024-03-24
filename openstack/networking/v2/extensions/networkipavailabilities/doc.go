/*
Package networkipavailabilities provides the ability to retrieve and manage
networkipavailabilities through the Neutron API.

Example of Listing NetworkIPAvailabilities

	allPages, err := networkipavailabilities.List(networkClient, networkipavailabilities.ListOpts{}).AllPages(context.TODO())
	if err != nil {
	  panic(err)
	}

	allAvailabilities, err := networkipavailabilities.ExtractNetworkIPAvailabilities(allPages)
	if err != nil {
	  panic(err)
	}

	for _, availability := range allAvailabilities {
	  fmt.Printf("%+v\n", availability)
	}

Example of Getting a single NetworkIPAvailability

	availability, err := networkipavailabilities.Get(context.TODO(), networkClient, "cf11ab78-2302-49fa-870f-851a08c7afb8").Extract()
	if err != nil {
	  panic(err)
	}

	fmt.Printf("%+v\n", availability)
*/
package networkipavailabilities
