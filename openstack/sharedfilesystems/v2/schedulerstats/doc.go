/*
Package schedulerstats returns information about shared file systems capacity
and utilisation.

Pool, Host, and Backend filters are regular expressions matched from the start
of the corresponding name. Use ^ and $ to require an exact match. Capability
filtering is available through ShareType, which matches the share type's extra
specs against pool capabilities and requires microversion 2.23 or later. The
project ID belongs in the service endpoint, not in the filter options.

# Listing pools

	listOpts := schedulerstats.ListOpts{
		Pool:    "^pool1$",
		Backend: "^backend1$",
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

# Listing pools with details

ListDetail accepts the same ListOpts. To include pool capabilities and filter
by a share type's extra specs, reuse the options above:

	client.Microversion = "2.23"
	listOpts.ShareType = "share-type-name-or-uuid"
	allPages, err = schedulerstats.ListDetail(client, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allStats, err = schedulerstats.ExtractPools(allPages)
	if err != nil {
		panic(err)
	}

	for _, stat := range allStats {
		fmt.Printf("%+v\n", stat)
	}
*/
package schedulerstats
