/*
Package manageablevolumes information and interaction with manageable volumes
for the OpenStack Block Storage service.

NOTE: Requires at least microversion 3.8

Example to manage an existing volume

	manageOpts := manageablevolumes.ManageExistingOpts{
		Host:             "host@lvm#LVM",
		Ref:              map[string]string{
			"source-name": "volume-73796b96-169f-4675-a5bc-73fc0f8f9a17",
		},
		Name:             "New Volume",
		AvailabilityZone: "nova",
		Description:      "Volume imported from existingLV",
		VolumeType:       "lvm",
		Bootable:         true,
		Metadata:         map[string]string{
			"key1": "value1",
			"key2": "value2"
		},
	}

	managedVolume, err := manageablevolumes.ManageExisting(context.TODO(), client, manageOpts).Extract()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Managed volume: %+v\n", managedVolume)
*/
package manageablevolumes
