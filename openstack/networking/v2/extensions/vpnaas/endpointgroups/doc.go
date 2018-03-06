/*
Package endpointgroups allows management of endpoint groups in the Openstack Network Service

Example to create an Endpoint Group

	createOpts := endpointgroups.CreateOpts{
		Name: groupName,
		Type: endpointgroups.TypeCIDR,
		Endpoints: []string{
			"10.2.0.0/24",
			"10.3.0.0/24",
		},
	}
	group, err := endpointgroups.Create(client, createOpts).Extract()
	if err != nil {
		return group, err
	}

Example to retrieve an Endpoint Group

	group, err := endpointgroups.Get(client, "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a").Extract()
	if err != nil {
		panic(err)
	}
*/
package endpointgroups
