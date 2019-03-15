package introspection

/*
Package introspection contains the functionality for Starting introspection,
Get introspection status, List all introspection statuses, Abort an
introspection, Get stored introspection data and reapply introspection on
stored data.

API reference https://developer.openstack.org/api-ref/baremetal-introspection/#node-introspection

    // Example to Start Introspection
    introspection.StartIntrospection(client, NodeUUID, introspection.StartOpts{}).ExtractErr()

    // Example to Get an Introspection status
    introspection.GetIntrospectionStatus(client, NodeUUID).Extract()
    if err != nil {
        panic(err)
    }

    // Example to List all introspection statuses
    introspection.ListIntrospections(client.ServiceClient(), introspection.ListIntrospectionsOpts{}).EachPage(func(page pagination.Page) (bool, error) {
        introspectionsList, err := introspection.ExtractIntrospections(page)
            if err != nil {
                return false, err
            }
            for _, n := range introspectionsList {
                // Do something
            }
        return true, nil
    })

    // Example to Abort an Introspection
    introspection.AbortIntrospection(client, NodeUUID).ExtractErr()

    // Example to Get stored Introspection Data
    introspection.GetIntrospectionData(c, NodeUUID).Extract()

    // Example to apply Introspection Data
    introspection.ApplyIntrospectionData(c, NodeUUID).ExtractErr()
*/
