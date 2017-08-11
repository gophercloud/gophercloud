/*
Package schedulerstats returns information about block storage pool capacity
and utilisation. Example:

	listOpts := schedulerstats.ListOpts{
		Detail: true,
	}

	allPages, err := schedulerstats.List(client, listOpts).AllPages()
	if err != nil {
		panic("Unable to query schedulerstats: %s", err)
	}

	allStats, err := schedulerstats.ExtractStoragePools(allPages)
	if err != nil {
		panic("Unable to extract pools: %s", err)
	}

	for _, stat := range allStats {
		fmt.Printf("%+v\n", stat)
	}
*/
package schedulerstats
