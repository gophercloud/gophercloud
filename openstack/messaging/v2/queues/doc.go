/*
Package queues provides information and interaction with the queues through
the OpenStack Messaging (Zaqar) service.

Lists all queues and creates, shows information for updates, deletes, and actions on a queue.

Example to Create a Queue

	createOpts := queues.CreateOpts{
		QueueName:					"My_Queue",
		MaxMessagesPostSize: 		262143,
		DefaultMessageTTL: 			3700,
		DefaultMessageDelay: 		25,
		DeadLetterQueueMessageTTL: 	3500,
		MaxClaimCount: 				10,
		Extra: 						map[string]interface{}{"description": "Test queue."},
	}

	err := queues.Create(client, clientID, createOpts).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package queues
