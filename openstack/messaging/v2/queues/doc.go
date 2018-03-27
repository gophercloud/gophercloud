/*
Package queues provides information and interaction with the queues through
<<<<<<< HEAD
the OpenStack Messaging (Zaqar) service.
=======
the OpenStack Messaging(Zaqar) service.
>>>>>>> Add create function for queues

Lists all queues and creates, shows information for updates, deletes, and actions on a queue.

Example to Create a Queue

	createOpts := queues.CreateOpts{
<<<<<<< HEAD
		QueueName:                  "My_Queue",
		MaxMessagesPostSize:        262143,
		DefaultMessageTTL:          3700,
		DefaultMessageDelay:        25,
		DeadLetterQueueMessageTTL:  3500,
		MaxClaimCount:              10,
		Extra:                      map[string]interface{}{"description": "Test queue."},
	}

	err := queues.Create(client, createOpts).ExtractErr()
=======
		MaxMessagesPostSize: 262143,
		DefaultMessageTTL: 3700,
		DefaultMessageDelay: 25,
		DeadLetterQueueMessageTTL: 3500,
		MaxClaimCount: 10,
		Description: "Test queue.",
	}

	err := queues.Create(client, "My Queue", clientID, createOpts).ExtractErr()
>>>>>>> Add create function for queues
	if err != nil {
		panic(err)
	}
*/
package queues
