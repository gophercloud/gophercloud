/*
Package quotas provides the ability to retrieve DNS quotas through the Designate API.

Example to Get a Quota Set

	projectID = "23d5d3f79dfa4f73b72b8b0b0063ec55"
	quotasInfo, err := quotas.Get(context.TODO(), dnsClient, projectID).Extract()
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Printf("quotas: %#v\n", quotasInfo)

Example to Update a Quota Set

	projectID = "23d5d3f79dfa4f73b72b8b0b0063ec55"
	zones := 10
	quota := &quotas.UpdateOpts{
	   Zones: &zones,
	}
	quotasInfo, err := quotas.Update(context.TODO(), dnsClient, projectID, quota).Extract()
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Printf("quotas: %#v\n", quotasInfo)
*/
package quotas
