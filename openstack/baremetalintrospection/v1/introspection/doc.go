/*
Package introspection contains the functionality for Starting introspection,
Get introspection status, List all introspection statuses, Abort an
introspection, Get stored introspection data and reapply introspection on
stored data.

API reference https://developer.openstack.org/api-ref/baremetal-introspection/#node-introspection

Example to Start Introspection

	err := introspection.StartIntrospection(context.TODO(), client, NodeUUID, introspection.StartOpts{}).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Get an Introspection status

	_, err := introspection.GetIntrospectionStatus(context.TODO(), client, NodeUUID).Extract()
	if err != nil {
		panic(err)
	}

Example to List all introspection statuses

	introspection.ListIntrospections(client.ServiceClient(), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		introspectionsList, err := introspection.ExtractIntrospections(page)
		if err != nil {
			return false, err
		}

		for _, n := range introspectionsList {
			// Do something
		}

		return true, nil
	})

Example to Abort an Introspection

	err := introspection.AbortIntrospection(context.TODO(), client, NodeUUID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Get stored Introspection Data

	v, err := introspection.GetIntrospectionData(c, NodeUUID).Extract()
	if err != nil {
		panic(err)
	}

Example to apply Introspection Data

	err := introspection.ApplyIntrospectionData(c, NodeUUID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package introspection
