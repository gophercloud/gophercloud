/*
Package receivers provides information and interaction with the receivers through
the OpenStack Clustering service.

Example to Create a Receiver

	createOpts := receivers.CreateOpts{
		Action:     "CLUSTER_DEL_NODES",
		ClusterID:  "b7b870ee-d3c5-4a93-b9d7-846c53b2c2dc",
		Name:       "test_receiver",
		Type:       "webhook",
	}

	receiver, err := receivers.Create(serviceClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Println("receiver", receiver)

Example to Get a Receiver

	receiver, err := receivers.Get(serviceClient, "receiver-name").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Print("receiver", receiver)

Example to Delete receiver

	receiverID := "6dc6d336e3fc4c0a951b5698cd1236ee"
	err := receivers.Delete(serviceClient, receiverID).ExtractErr()
	if err != nil {
		panic(err)
	}

	fmt.Print("receiver", receiver)

Example to Update Receiver

	receiver, err := receivers.Update(serviceClient, receiverName, receivers.UpdateOpts{Name: newReceiverName}).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Print("receiver", receiver)

Example to List Receivers

	receivers.List(serviceClient, receivers.ListOpts{Limit: 2}).EachPage(func(page pagination.Page) (bool, error) {
		allReceivers, err := receivers.ExtractReceivers(page)
		if err != nil {
			panic(err)
		}

		for _, receiver := range allReceivers {
			fmt.Printf("%+v\n", receiver)
		}
		return true, nil
	})
*/
package receivers
