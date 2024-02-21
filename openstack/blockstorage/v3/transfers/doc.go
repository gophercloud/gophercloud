/*
Package transfers provides an interaction with volume transfers in the
OpenStack Block Storage service. A volume transfer allows to transfer volumes
between projects withing the same OpenStack region.

Example to List all Volume Transfer requests being an OpenStack admin

	listOpts := &transfers.ListOpts{
		// this option is available only for OpenStack cloud admin
		AllTenants: true,
	}

	allPages, err := transfers.List(client, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allTransfers, err := transfers.ExtractTransfers(allPages)
	if err != nil {
		panic(err)
	}

	for _, transfer := range allTransfers {
		fmt.Println(transfer)
	}

Example to Create a Volume Transfer request

	createOpts := transfers.CreateOpts{
		VolumeID: "uuid",
		Name:	  "my-volume-transfer",
	}

	transfer, err := transfers.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(transfer)
	// secret auth key is returned only once as a create response
	fmt.Printf("AuthKey: %s\n", transfer.AuthKey)

Example to Accept a Volume Transfer request from the target project

	acceptOpts := transfers.AcceptOpts{
		// see the create response above
		AuthKey: "volume-transfer-secret-auth-key",
	}

	// see the transfer ID from the create response above
	transfer, err := transfers.Accept(context.TODO(), client, "transfer-uuid", acceptOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(transfer)

Example to Delete a Volume Transfer request from the source project

	err := transfers.Delete(context.TODO(), client, "transfer-uuid").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package transfers
