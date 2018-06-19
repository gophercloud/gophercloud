/*
Package receivers provides information and interaction with the receivers through
the OpenStack Compute service.

Lists all receivers and creates, shows information for, and deletes a receiver.

Example to Delete receiver

	receiverID := "6dc6d336e3fc4c0a951b5698cd1236ee"
	err := receivers.Delete(serviceClient, receiverID).ExtractErr()
	if err != nil {
		panic(err)
	}

*/
package receivers
