/*
Package receivers provides information and interaction with the receivers through
the OpenStack Compute service.

Lists all receivers and creates, shows information for, and deletes a receiver.

Example to List Receivers

	listOpts := receivers.ListOpts{
		Limit: 2,
	}

	allPages, err := images.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allReceivers, err := images.ExtractImages(allPages)
	if err != nil {
		panic(err)
	}

	for _, receiver := range allReceivers {
		fmt.Printf("%+v\n", receiver)
	}

Example to Create a Receiver

	createOpts := receivers.CreateOpts{
    Action:     "CLUSTER_DEL_NODES"
    ClusterID:  "b7b870ee-d3c5-4a93-b9d7-846c53b2c2dc"
    Name:       "test_receiver"
    Type:       "webhook"
	}

	receiver, err := receivers.Create(computeClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package receivers
