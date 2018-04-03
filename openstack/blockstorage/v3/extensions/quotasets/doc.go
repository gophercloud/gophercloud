/*
Package quotasets enables retrieving and managing Compute quotas.

Example to Get a Quota Set

	quotaset, err := quotasets.Get(computeClient, "project-id").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Get a Detailed Quota Set

	quotaset, err := quotasets.GetDetail(computeClient, "project-id").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)

Example to Update a Quota Set

	updateOpts := quotasets.UpdateOpts{
		FixedIPs: gophercloud.IntToPointer(100),
		Cores:    gophercloud.IntToPointer(64),
	}

	quotaset, err := quotasets.Update(computeClient, "project-id", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)
*/
package quotasets
