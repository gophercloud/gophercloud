/*
Package limits provides information and interaction with limits for the
Openstack Identity service.

Example to Get EnforcementModel

	model, err := limits.GetEnforcementModel(context.TODO(), identityClient).Extract()
	if err != nil {
		panic(err)
	}

Example to List Limits

	listOpts := limits.ListOpts{
		ProjectID: "3d596369fd2043bf8aca3c8decb0189e",
	}

	allPages, err := limits.List(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allLimits, err := limits.ExtractLimits(allPages)
	if err != nil {
		panic(err)
	}

Example to Create Limits

	batchCreateOpts := limits.BatchCreateOpts{
		limits.CreateOpts{
			ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
			ProjectID:     "3a705b9f56bb439381b43c4fe59dccce",
			RegionID:      "RegionOne",
			ResourceName:  "snapshot",
			ResourceLimit: 5,
		},
		limits.CreateOpts{
			ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
			DomainID:      "edbafc92be354ffa977c58aa79c7bdb2",
			ResourceName:  "volume",
			ResourceLimit: 10,
			Description:   "Number of volumes for project 3a705b9f56bb439381b43c4fe59dccce",
		},
	}

	createdLimits, err := limits.Create(context.TODO(), identityClient, batchCreateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a Limit

	limit, err := limits.Get(context.TODO(), identityClient, "25a04c7a065c430590881c646cdcdd58").Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Limit

	limitID := "0fe36e73809d46aeae6705c39077b1b3"

	description := "Number of snapshots for project 3a705b9f56bb439381b43c4fe59dccce"
	resourceLimit := 5
	updateOpts := limits.UpdateOpts{
	  Description:   &description,
	  ResourceLimit: &resourceLimit,
	}

	limit, err := limits.Update(context.TODO(), identityClient, limitID, updateOpts).Extract()
	if err != nil {
	  panic(err)
	}

Example to Delete a Limit

	limitID := "0fe36e73809d46aeae6705c39077b1b3"
	err := limits.Delete(context.TODO(), identityClient, limitID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package limits
