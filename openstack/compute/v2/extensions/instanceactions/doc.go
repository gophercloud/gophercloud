package instanceactions

/*
Package instanceactions provides the ability to list or get a server instance-action.
Example:

	pages, err := instanceactions.List(client, "server-id").AllPages()
	if err != nil {
		panic("fail to get actions pages")
	}
	actions, err := instanceactions.ExtractInstanceActions(pages)
	if err != nil {
		panic("fail to list instance actions")
	}

	for _, action := range actions {
		action, err = instanceactions.Get(client, "server-id", action.RequestID).Extract()
		if err != nil {
			panic("fail to get instance action")
		}
		fmt.Println(action)
	}
*/
