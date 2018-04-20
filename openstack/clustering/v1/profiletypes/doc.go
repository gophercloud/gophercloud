/*
Package profiletypes lists all profile types and shows details for a profile type from the OpenStack
Clustering Service.

Example to List ProfileType

	err = profiletypes.List(serviceClient).EachPage(func(page pagination.Page) (bool, error) {
		profileTypes, err := profiletypes.ExtractProfileTypes(page)
		if err != nil {
			return false, err
		}

		for _, profileType := range profileTypes {
			fmt.Println("%+v\n", profileType)
		}
		return true, nil
	})

Example to Get a ProfileType

	profileTypeName := "os.nova.server"
	profileType, err := profiletypes.Get(clusteringClient, profileTypeName).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", profileType)

*/
package profiletypes
