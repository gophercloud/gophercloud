/*
Package attachments provides access to OpenStack Block Storage Attachment
API's. Use of this package requires Cinder version 3.27 at a minimum.

For more information, see:
https://docs.openstack.org/api-ref/block-storage/v3/index.html#attachments

Example to List Attachments

	listOpts := &attachments.ListOpts{
		InstanceID: "uuid",
	}

	client.Microversion = "3.27"
	allPages, err := attachments.List(client, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allAttachments, err := attachments.ExtractAttachments(allPages)
	if err != nil {
		panic(err)
	}

	for _, attachment := range allAttachments {
		fmt.Println(attachment)
	}

Example to Create Attachment

	opts := &attachments.CreateOpts{
		InstanceiUUID: "uuid",
		VolumeUUID: "uuid"
	}

	client.Microversion = "3.27"
	attachment, err := attachments.Create(client, opts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(attachment)

Example to Get Attachment

	client.Microversion = "3.27"
	attachment, err := attachments.Get(client, "uuid").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(attachment)

Example to Update Attachment

	opts := &attachments.UpdateOpts{
		Connector: map[string]interface{}{
			"mode": "ro",
		}
	}

	client.Microversion = "3.27"
	attachment, err := attachments.Update(client, "uuid", opts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(attachment)

Example to Complete Attachment

	client.Microversion = "3.44"
	err := attachments.Complete(client, "uuid").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Delete Attachment

	client.Microversion = "3.27"
	err := attachments.Delete(client, "uuid").ExtractErr()
	if err != nil {
		panic(err)
	}

*/
package attachments
