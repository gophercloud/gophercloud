/*
Package registeredlimits provides information and interaction with registered limits for the
Openstack Identity service.

Example to List RegisteredLimits

	listOpts := registeredlimits.ListOpts{
		ResourceName: "image_size_total",
	}

	allPages, err := registeredlimits.List(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allLimits, err := limits.ExtractLimits(allPages)
	if err != nil {
		panic(err)
	}

Example to Create a RegisteredLimit

	batchCreateOpts := registeredlimits.BatchCreateOpts{
		registeredlimits.CreateOpts{
			ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
			RegionID:      "RegionOne",
			ResourceName:  "snapshot",
			DefaultLimit: 5,
		},
		registeredlimits.CreateOpts{
			ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
			RegionID:      "RegionOne",
			ResourceName:  "volume",
			DefaultLimit: 10,
			Description:   "Number of volumes for service 9408080f1970482aa0e38bc2d4ea34b7",
		},
	}

	createdRegisteredLimits, err := limits.Create(context.TODO(), identityClient, batchCreateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a RegisteredLimit

	    registeredLimitID := "966b3c7d36a24facaf20b7e458bf2192"
	    registered_limit, err := registeredlimits.Get(context.TODO(), client, registeredLimitID).Extract()
		if err != nil {
			panic(err)
		}

Example to Update a RegisteredLimit

	    Either ServiceID, ResourceName, or RegionID must be different than existing value otherwise it will raise 409.

		registeredLimitID := "966b3c7d36a24facaf20b7e458bf2192"

		resourceName := "images"
		description := "Number of images for service 9408080f1970482aa0e38bc2d4ea34b7"
		defaultLimit := 10
		updateOpts := registeredlimits.UpdateOpts{
			Description:  &description,
			DefaultLimit: &defaultLimit,
			ResourceName: resourceName,
			ServiceID:    "9408080f1970482aa0e38bc2d4ea34b7",
		}

		registered_limit, err := registeredlimits.Update(context.TODO(), client, registeredLimitID, updateOpts).Extract()
		if err != nil {
			panic(err)
		}

Example to Delete a RegisteredLimit

	registeredLimitID := "966b3c7d36a24facaf20b7e458bf2192"
	err := registeredlimits.Delete(context.TODO(), identityClient, registeredLimitID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package registeredlimits
