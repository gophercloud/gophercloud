/*
Package endpointgroups enables management of OpenStack Identity Endpoint Groups
and Endpoint associations.

Example to Get an Endpoint Group

    err := endpointgroups.Get(identityClient, endpointGroupID).ExtractErr()
    if err != nil {
    	panic(err)
    }

Example to List all Endpoint Groups by name

	listOpts := endpointgropus.ListOpts{
		Name: "mygroup",
	}

    allPages, err := endpointgroups.List(identityClient, listOpts).AllPages()
    if err != nil {
    	panic(err)
    }

	allGroups, err := endpointgroups.ExtractEndpointGroups(allPages)
    if err != nil {
    	panic(err)
    }

	for _, endpointgroup := range allGroups {
		fmt.Printf("%+v\n", endpointgroup)
	}
*/
package endpointgroups
