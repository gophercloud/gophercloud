/*
Package addressgroups provides information and interaction with Address Groups
for the OpenStack Networking services.

Example to List Address Groups

	listOpts := addressgroups.ListOpts{
	}

	allPages, err := addressgroups.List(networkClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allAddressGroups, err := addressgroups.ExtractGroups(allPages)
	if err != nil {
		panic(err)
	}

	for _, addressGroup := range allAddressGroups {
		fmt.Printf("%+v\n", addressGroup)
	}

Example to Create an Address Group

	createOpts := addressgroups.CreateOpts{
		Name:        "addressGroupName",
		Addresses:   []string{"10.2.30.4/32", "10.2.30.6/32"},
		Description: "Created address group",
	}

	addressGroup, err := addressgroups.Create(context.TODO(), networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Address Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	err := addressgroups.Delete(context.TODO(), computeClient, groupID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to update an existing Address Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	createOpts := addressGroups.CreateOpts{
		Addresses: []string{"10.2.30.4/32", "10.2.30.6/32"},
		Description: "new description",
		Name: "ADDR_GP_2"
	}
	addressGroup, err := addressgroups.UpdateAddressGroup(context.TODO(), networkClient, groupID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to add addresses to an existing Address Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	createOpts := addressGroups.CreateOpts{
		Addresses: []string{"10.2.30.4/32", "10.2.30.6/32"},
	}
	addressGroup, err := addressgroups.AddAddresses(context.TODO(), networkClient, groupID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to remove addresses from an existing Address Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	createOpts := addressGroups.CreateOpts{
		Addresses: []string{"10.2.30.4/32", "10.2.30.6/32"},
	}
	addressGroup, err := addressgroups.RemoveAddresses(context.TODO(), networkClient, groupID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/

package addressgroups
