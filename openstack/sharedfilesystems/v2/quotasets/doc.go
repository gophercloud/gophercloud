/*
Package quotasets provides information and interaction with the quotasets API for the OpenStack Shared Filesystems service.

Example to Get a Quota Set

	quotaset, err := quotasets.Get(sharedfilesystemsClient, "tenant-id").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Update a Quota Set

	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.IntToPointer(100),
	}

	quotaset, err := quotasets.Update(sharedfilesystemsClient, "tenant-id", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Get a Quota Set by Share Type

	quotaset, err := quotasets.GetByShareType(sharedfilesystemsClient, "tenant-id", "default").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Update a Quota Set by Share Type

	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.IntToPointer(100),
	}

	quotaset, err := quotasets.UpdateByShareType(sharedfilesystemsClient, "tenant-id", "default", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Get a Quota Set by User

	quotaset, err := quotasets.GetByUser(sharedfilesystemsClient, "tenant-id", "user-id").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Update a Quota Set by User

	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.IntToPointer(100),
	}

	quotaset, err := quotasets.UpdateByUser(sharedfilesystemsClient, "tenant-id", "user-id", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)
*/
package quotasets
