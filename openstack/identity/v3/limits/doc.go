/*
Package limits provides information and interaction with limits for the
Openstack Identity service.

Example to List Limits

	listOpts := limits.ListOpts{
		ProjectID: "3d596369fd2043bf8aca3c8decb0189e",
	}

	allPages, err := limits.List(identityClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allLimits, err := limits.ExtractLimits(allPages)
	if err != nil {
		panic(err)
	}
*/
package limits
