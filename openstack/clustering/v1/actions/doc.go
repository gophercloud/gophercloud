/*
Package actions enables listing actions for senlin engine from the OpenStack
Clustering Service.

Example to list actions

  listOpts := actions.ListOpts{
    Limit: 2,
  }

  allPages, err := actions.ListDetail(computeClient, listOpts).AllPages()
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
