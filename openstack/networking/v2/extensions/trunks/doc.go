/*
Package trunks provides the ability to retrieve and manage trunks through the Neutron API.
Trunks allow you to multiplex multiple ports traffic on a single port. For example, you could
have a compute instance port be the parent port of a trunk and inside the VM run workloads
using other ports, without the need of plugging those ports.

Example of a new empty Trunk creation

	iTrue := true
	createOpts := trunks.CreateOpts{
		Name:         "gophertrunk",
		Description:  "Trunk created by gophercloud",
		AdminStateUp: &iTrue,
		PortID:       "a6f0560c-b7a8-401f-bf6e-d0a5c851ae10",
	}

	trunk, err := trunks.Create(networkClient, trunkOpts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", trunk)

Example of a new Trunk creation with 2 subports

	iTrue := true
	createOpts := trunks.CreateOpts{
		Name:         "gophertrunk",
		Description:  "Trunk created by gophercloud",
		AdminStateUp: &iTrue,
		PortID:       "a6f0560c-b7a8-401f-bf6e-d0a5c851ae10",
		Subports: []trunks.Subport{
			{
				SegmentationID:   1,
				SegmentationType: "vlan",
				PortID:           "bf4efcc0-b1c7-4674-81f0-31f58a33420a",
			},
			{
				SegmentationID:   10,
				SegmentationType: "vlan",
				PortID:           "2cf671b9-02b3-4121-9e85-e0af3548d112",
			},
		},
	}

	trunk, err := trunks.Create(client, trunkOpts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", trunk)

Example of deleting a Trunk

	trunkID := "c36e7f2e-0c53-4742-8696-aee77c9df159"
	err := trunks.Delete(networkClient, trunkID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example of listing Trunks

	listOpts := trunks.ListOpts{}
	allPages, err := trunks.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}
	allTrunks, err := trunks.ExtractTrunks(allPages)
	if err != nil {
		panic(err)
	}
	for _, trunk := range allTrunks {
		fmt.Printf("%+v\n", trunk)
	}


*/
package trunks
