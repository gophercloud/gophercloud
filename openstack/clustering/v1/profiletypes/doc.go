/*
Package profile types lists all profile types and shows details for a profile type from the OpenStack
Clustering Service.

Example to list profile types for a Senlin deployment

  allPages, err := profiletypes.ListDetail(computeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allProfileTypes, err := profiletypes.ExtractProfileTypes(allPages)
	if err != nil {
		panic(err)
	}

	for _, profileType := range allProfileTypes {
		fmt.Printf("%+v\n", profileType)
	}
*/
package profiletypes
