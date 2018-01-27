/*
Package aggregates creates aggregates and returns information about the host aggregates in the
OpenStack cloud.

Example of Create an aggregate

	opts := aggregates.CreateOpts{
		Name:             "name",
		AvailabilityZone: "london",
	}

	aggregate, err := aggregates.Create(computeClient, opts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", aggregate)

Example of Retrieving list of all aggregates

	allPages, err := aggregates.List(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allAggregates, err := aggregates.ExtractAggregates(allPages)
	if err != nil {
		panic(err)
	}

	for _, aggregate := range allAggregates {
		fmt.Printf("%+v\n", aggregate)
	}
*/
package aggregates
