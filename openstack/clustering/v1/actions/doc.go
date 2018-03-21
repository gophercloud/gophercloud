/*
Package actions provides listing and retrieving of senlin actions for the OpenStack Clustering Service.

Example to list actions

	allPages, err := actions.List(serviceClient, actions.ListOpts{Limit: 5}).AllPages()
	if err != nil {
		panic(err)
	}

	allActions, err := actions.ExtractActions(allPages)
	if err != nil {
		panic(err)
	}

	for _, action := range allActions {
		fmt.Printf("%+v\n", action)
	}
*/
package actions
