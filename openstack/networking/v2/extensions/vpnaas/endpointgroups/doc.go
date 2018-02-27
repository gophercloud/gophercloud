/*
Package endpointgroups allows management of endpoint groups in the Openstack Network Service

Example to create an Endpoint Group

	createOpts := endpointgroups.CreateOpts{
		Name: groupName,
		Type: endpointgroups.TypeCidr,
		Endpoints: []string{
			"10.2.0.0/24",
			"10.3.0.0/24",
		},
	}
	group, err := endpointgroups.Create(client, createOpts).Extract()
	if err != nil {
		return group, err
	}
*/
package endpointgroups
