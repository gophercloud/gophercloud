/*
Package quotasets enables retrieving and managing Compute quotas.

Example to Get a Quota Set

	quotaset, err := quotasets.Get(context.TODO(), computeClient, "tenant-id").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Get a Detailed Quota Set

	quotaset, err := quotasets.GetDetail(context.TODO(), computeClient, "tenant-id").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Update a Quota Set

	updateOpts := quotasets.UpdateOpts{
		FixedIPs: gophercloud.IntToPointer(100),
		Cores:    gophercloud.IntToPointer(64),
	}

	quotaset, err := quotasets.Update(context.TODO(), computeClient, "tenant-id", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)
*/
package quotasets
