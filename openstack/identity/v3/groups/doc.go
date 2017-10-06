/*
Package groups manages and retrieves Groups in the OpenStack Identity Service.

Example to List Groups

	listOpts := groups.ListOpts{
		DomainID: "default",
	}

	allPages, err := groups.List(identityClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		panic(err)
	}

	for _, group := range allGroups {
		fmt.Printf("%+v\n", group)
	}
*/
package groups
