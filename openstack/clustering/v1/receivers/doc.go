/*
Package receivers provides information and interaction with the receivers through
the OpenStack Compute service.

Lists all receivers and creates, shows information for, and deletes a receiver.

Example to Get a Receiver

	receiver, err := receivers.Get(serviceClient, "receiver-name").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Print("receiver", receiver)

*/
package receivers
