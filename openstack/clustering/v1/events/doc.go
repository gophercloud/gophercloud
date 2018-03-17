/*
Package events provides information and interaction with the events through
the OpenStack Compute service.

Lists all events and shows information for an event.

Example to List Events

	listOpts := events.ListOpts{
		Limit: 2,
	}

	allPages, err := events.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allEvents, err := events.ExtractEvents(allPages)
	if err != nil {
		panic(err)
	}

	for _, event := range allEvents {
		fmt.Printf("%+v\n", event)
	}
*/
package events
