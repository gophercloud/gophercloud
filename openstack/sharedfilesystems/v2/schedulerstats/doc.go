/*
Package schedulerstats returns information about shared file systems capacity
and utilisation. Example:

	listOpts := schedulerstats.ListOpts{
	}

	allPages, err := schedulerstats.List(client, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allStats, err := schedulerstats.ExtractPools(allPages)
	if err != nil {
		panic(err)
	}

	for _, stat := range allStats {
		fmt.Printf("%+v\n", stat)
	}
*/
package schedulerstats
