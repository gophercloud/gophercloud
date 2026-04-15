/*
Package usages retrieves total resource usage from the OpenStack Placement service.

Usage API requests are available starting from microversion 1.9.

# Example to get total usages grouped by consumer type (microversion 1.38+)

	placementClient.Microversion = "1.38"

	totalUsages, err := usages.Get(context.TODO(), placementClient, usages.GetOpts{
		ProjectID: projectID,
	}).Extract()
	if err != nil {
		panic(err)
	}

	for consumerType, usage := range totalUsages.Usages {
		fmt.Printf("%s: VCPU=%d, consumer_count=%d\n",
			consumerType, usage["VCPU"], usage["consumer_count"])
	}

# Example to get total usages without consumer type grouping (microversion 1.9–1.37)

	placementClient.Microversion = "1.9"

	totalUsages, err := usages.Get(context.TODO(), placementClient, usages.GetOpts{
		ProjectID: projectID,
	}).ExtractPre138()
	if err != nil {
		panic(err)
	}

	fmt.Printf("VCPU usage: %d\n", totalUsages.Usages["VCPU"])

# Example to get total usages for a specific project and user

	placementClient.Microversion = "1.38"

	totalUsages, err := usages.Get(context.TODO(), placementClient, usages.GetOpts{
		ProjectID: projectID,
		UserID:    userID,
	}).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", totalUsages)
*/
package usages
