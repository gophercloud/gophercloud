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

Example to Create a Group

	createOpts := groups.CreateOpts{
		Name:             "groupname",
		DomainID:         "default",
		Extra: map[string]interface{}{
			"email": "groupname@example.com",
		}
	}

	group, err := groups.Create(identityClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package groups
