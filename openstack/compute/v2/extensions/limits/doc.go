/*
Package limits shows rate and limit information for a tenant/project.

Example to Retrieve Limits for a Tenant

	getOpts := limits.GetOpts{
		TenantID: "tenant-id",
	}

	limits, err := limits.Get(computeClient, getOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println("%+v\n", limits)
*/
package limits
