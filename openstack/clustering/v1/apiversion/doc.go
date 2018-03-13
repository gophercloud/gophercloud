/*
Package apiversion enables listing API version for senlin engine from the OpenStack
Clustering Service.

Example to list API version

  allPages, err := apiversion.ListDetail(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allVersions, err := actions.ExtractVersions(allPages)
	if err != nil {
		panic(err)
	}

	for _, version := range allVersions {
		fmt.Printf("%+v\n", version)
	}
*/
package apiversion
