/*
Package buildinfo enables listing build information for senlin engine from the OpenStack
Clustering Service.

Example to list build information for a Senlin deployment

  allPages, err := buildinfo.ListDetail(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allBuildInfos, err := actions.ExtractBuildInfos(allPages)
	if err != nil {
		panic(err)
	}

	for _, buildinfo := range allBuildInfos {
		fmt.Printf("%+v\n", buildinfo)
	}
*/
package buildinfo
