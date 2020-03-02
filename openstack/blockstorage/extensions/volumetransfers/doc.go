/*
Package volumetransfers provides an interaction with volume transfers in the
OpenStack Block Storage service. A volume transfer allows to transfer volumes
between projects withing the same OpenStack region.

Example to List all Volume Transfer requests being an OpenStack admin

	listOpts := &volumetransfers.ListOpts{
		// this option is available only for OpenStack cloud admin
		AllTenants: true,
	}

	allPages, err := volumetransfers.List(client, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allTransfers, err := volumetransfers.ExtractTransfers(allPages)
	if err != nil {
		panic(err)
	}

	for _, transfer := range allTransfers {
		fmt.Println(transfer)
	}

Example to Create a Volume Transfer request

	createOpts := volumetransfers.CreateOpts{
                VolumeID: "uuid",
		Name: "my-volume-transfer",
        }

	transfer, err := volumetransfers.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(transfer)
	// secret auth key is returned only once as a create response
	fmt.Printf("AuthKey: %s\n", transfer.AuthKey)

Example to Accept a Volume Transfer request from the target project

	acceptOpts := volumetransfers.AcceptOpts{
		// see the create response above
		AuthKey: "volume-transfer-secret-auth-key",
	}

	// see the transfer ID from the create response above
	transfer, err := volumetransfers.Accept(client, "transfer-uuid", acceptOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(transfer)

Example to Delete a Volume Transfer request from the source project

	err := volumetransfers.Delete(client, "transfer-uuid").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package volumetransfers
