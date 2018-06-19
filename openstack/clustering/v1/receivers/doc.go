/*
Package receivers provides information and interaction with the receivers through
the OpenStack Compute service.

Lists all receivers and creates, shows information for, and deletes a receiver.

Example to Update Receiver

	receiver, err := receivers.Update(serviceClient, receiverName, receivers.UpdateOpts{Name: newReceiverName}).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Print("receiver", receiver)

*/
package receivers
