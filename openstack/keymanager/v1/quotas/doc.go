/*
Package quotas provides the ability to retrieve and manage Barbican quotas

Example to Get project quotas

	projectID = "23d5d3f79dfa4f73b72b8b0b0063ec55"
	quotasInfo, err := quotas.Get(keyManagerClient, projectID).Extract()
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Printf("quotas: %#v\n", quotasInfo)

Example to Update project quotas

	projectID = "23d5d3f79dfa4f73b72b8b0b0063ec55"

	updateOpts := quotas.UpdateOpts{
	    Secrets:    gophercloud.IntToPointer(10),
	    Orders:     gophercloud.IntToPointer(20),
	    Containers: gophercloud.IntToPointer(10),
	    Consumers:  gophercloud.IntToPointer(-1),
	    Cas:        gophercloud.IntToPointer(5),
	}
	quotasInfo, err := quotas.Update(keyManagerClient, projectID)
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Printf("quotas: %#v\n", quotasInfo)
*/
package quotas
