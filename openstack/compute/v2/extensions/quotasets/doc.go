/*
Package quotasets enables retrieving and managing Compute quotas.

Example to Get a Quota Set

	quotaset, err := quotasets.Get(computeClient, "tenant-id").Extract()
	if err != nil {
		panic("Unable to retrieve quotaset: %s", err)
	}

	fmt.Println("%+v\n", quotaset)

Example to Update a Quota Set

	updateOpts := quotasets.UpdateOpts{
		FixedIPs: gophercloud.IntToPointer(100),
		Cores:    gophercloud.IntToPointer(64),
	}

	quotaset, err := quotasets.Update(computeClient, "tenant-id", updateOpts).Extract()
	if err != nil {
		panic("Unable to update quotaset: %s", err)
	}

	fmt.Println("%+v\n", quotaset)
*/
package quotasets
