/*
Package queues provides information and interaction with the queues through
the OpenStack Messaging (Zaqar) service.

Lists all queues and creates, shows information for updates, deletes, and actions on a queue.

Example to List Queues

	listOpts := queues.ListOpts{
		Limit: 10,
	}

	clientID := "3381af92-2b9e-11e3-b191-71861300734d"

	pager := queues.List(client, clientID, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		queues, err := queues.ExtractQueues(page)
		if err != nil {
			panic(err)
		}

		for _, queue := range queues {
			fmt.Printf("%+v\n", queue)
		}

		return true, nil
	})

Example to Create a Queue

	createOpts := queues.CreateOpts{
		QueueName:                  "My_Queue",
		MaxMessagesPostSize:        262143,
		DefaultMessageTTL:          3700,
		DefaultMessageDelay:        25,
		DeadLetterQueueMessageTTL:  3500,
		MaxClaimCount:              10,
		Extra:                      map[string]interface{}{"description": "Test queue."},
	}

	err := queues.Create(client, createOpts).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package queues
