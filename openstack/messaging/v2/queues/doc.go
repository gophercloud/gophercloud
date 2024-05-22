/*
Package queues provides information and interaction with the queues through
the OpenStack Messaging (Zaqar) service.

Lists all queues and creates, shows information for updates, deletes, and actions on a queue.

Example to List Queues

	listOpts := queues.ListOpts{
		Limit: 10,
	}

	pager := queues.List(cclient, listOpts)

	err = pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
		Extra:                      map[string]any{"description": "Test queue."},
	}

	err := queues.Create(context.TODO(), client, createOpts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Update a Queue

	updateOpts := queues.BatchUpdateOpts{
		queues.UpdateOpts{
			Op:    "replace",
			Path:  "/metadata/_max_claim_count",
			Value: 15,
		},
		queues.UpdateOpts{
			Op: "replace",
			Path: "/metadata/description",
			Value: "Updated description test queue.",
		},
	}

	updateResult, err := queues.Update(context.TODO(), client, queueName, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a Queue

	queue, err := queues.Get(context.TODO(), client, queueName).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Queue

	err := queues.Delete(context.TODO(), client, queueName).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Get Message Stats of a Queue

	queueStats, err := queues.GetStats(context.TODO(), client, queueName).Extract()
	if err != nil {
		panic(err)
	}

Example to Share a queue

	shareOpts := queues.ShareOpts{
		Paths:   []queues.SharePath{queues.ShareMessages},
		Methods: []queues.ShareMethod{queues.MethodGet},
	}

	queueShare, err := queues.Share(context.TODO(), client, queueName, shareOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Purge a queue

	purgeOpts := queues.PurgeOpts{
		ResourceTypes: []queues.PurgeResource{
			queues.ResourceMessages,
		},
	}

	err := queues.Purge(context.TODO(), client, queueName, purgeOpts).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package queues
