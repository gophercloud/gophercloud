/*
Package receivers provides information and interaction with the receivers through
the OpenStack Compute service.

Lists all receivers and creates, shows information for, and deletes a receiver.

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
