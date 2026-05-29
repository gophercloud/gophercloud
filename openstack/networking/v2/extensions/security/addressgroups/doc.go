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

Example to Get an Address Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	addressGroup, err := addressgroups.Get(context.TODO(), networkClient, groupID).Extract()
	if err != nil {
		panic(err)
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
	name := "ADDR_GP_2"
	description := "new description"
	updateOpts := addressgroups.UpdateOpts{
		Name:        &name,
		Description: &description,
	}
	addressGroup, err := addressgroups.UpdateAddressGroup(context.TODO(), networkClient, groupID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to add addresses to an existing Address Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	createOpts := addressgroups.UpdateAddressesOpts{
		Addresses: []string{"10.2.30.4/32", "10.2.30.6/32"},
	}
	addressGroup, err := addressgroups.AddAddresses(context.TODO(), networkClient, groupID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to remove addresses from an existing Address Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	createOpts := addressgroups.UpdateAddressesOpts{
		Addresses: []string{"10.2.30.4/32", "10.2.30.6/32"},
	}
	addressGroup, err := addressgroups.RemoveAddresses(context.TODO(), networkClient, groupID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/

package addressgroups
